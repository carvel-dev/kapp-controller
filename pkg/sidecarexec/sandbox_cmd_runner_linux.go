// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package sidecarexec

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	goexec "os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

var (
	limitDir = os.Getenv("KAPPCTRL_SIDECAREXEC_LIMIT_DIR")
)

// Run executes exec.Cmd under new process/mount namespace.
func (r SandboxCmdRunner) Run(cmd *goexec.Cmd, opts exec.RunOpts) error {
	newRootDir, err := ioutil.TempDir("", "sidecarexecwrap-cmd")
	if err != nil {
		return fmt.Errorf("Creating new root dir: %s", err)
	}

	wrapArgs := sandboxWrapArgs{
		NewRootDir: newRootDir,
		Posix:      r.opts.RequiresPosix != nil && r.opts.RequiresPosix[filepath.Base(cmd.Path)],
		Network:    r.opts.RequiresNetwork != nil && r.opts.RequiresNetwork[filepath.Base(cmd.Path)],

		CmdDir:  cmd.Dir,
		CmdPath: cmd.Path,
		CmdArgs: cmd.Args[1:], // drop first arg (binary name)

		VisiblePaths: opts.VisiblePaths,
	}

	// Roughly equivalent to: unshare -Urm -R / ...
	cmd.Path = "/usr/bin/kapp-controller"
	cmd.Args = []string{"kapp-controller", "--sidecarexecwrap", wrapArgs.AsString()}
	cmd.Dir = "/"
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// TODO syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID?
		Cloneflags:  syscall.CLONE_NEWUSER | syscall.CLONE_NEWNS,
		UidMappings: []syscall.SysProcIDMap{{ContainerID: 0, HostID: os.Geteuid(), Size: 1}},
		GidMappings: []syscall.SysProcIDMap{{ContainerID: 0, HostID: os.Getegid(), Size: 1}},
	}

	// Create new non-configured network namespace if networking is not required
	if !wrapArgs.Network {
		cmd.SysProcAttr.Cloneflags = cmd.SysProcAttr.Cloneflags | syscall.CLONE_NEWNET
	}

	err = r.local.Run(cmd, opts)
	if err != nil {
		_ = os.RemoveAll(newRootDir) // remove error is lower priority
		return err
	}

	// Propagate up any errors deleting root directory
	// (want to notice any deletion problems asap)
	err = os.RemoveAll(newRootDir)
	if err != nil {
		return fmt.Errorf("Deleting new root dir: %s", err)
	}

	return nil
}

// RunWithCancel executes exec.Cmd.
// Kills execution immediately if value is read from cancelCh.
func (SandboxCmdRunner) RunWithCancel(cmd *goexec.Cmd, cancelCh chan struct{}, opts exec.RunOpts) error {
	return fmt.Errorf("RunWithCancel not implemented")
}

type sandboxWrapArgs struct {
	NewRootDir string
	Posix      bool
	Network    bool

	CmdDir  string
	CmdPath string
	CmdArgs []string

	VisiblePaths []string
}

func newSandboxWrapArgs(str string) sandboxWrapArgs {
	var args sandboxWrapArgs
	err := json.Unmarshal([]byte(str), &args)
	if err != nil {
		panic(fmt.Sprintf("Internal inconsistency: invalid sandboxWrapArgs: %s", err))
	}
	return args
}

func (a sandboxWrapArgs) AsString() string {
	bs, err := json.Marshal(a)
	if err != nil {
		panic(fmt.Sprintf("Internal inconsistency: invalid sandboxWrapArgs: %s", err))
	}
	return string(bs)
}

// SandboxWrap represents continuation of sandbox setup (mounting, ...)
// which sandbox has initiated within a new user/mount namespace.
type SandboxWrap struct{}

type mountPoint struct {
	Path      string
	File      bool
	Writeable bool
	NonMount  bool
	Content   string
}

// ExecuteCmd mounts various directories needed to execute sandboxed command.
// As a last step it exec-s into given command, replacing the current process.
func (w SandboxWrap) ExecuteCmd(rawArgs []string) error {
	runtime.LockOSThread() // as a safety precaution for syscalls below

	args := newSandboxWrapArgs(rawArgs[0])

	// Keep default environment as minimal as possible
	mounts := []mountPoint{
		{Path: "usr/bin"},
		{Path: "tmp", NonMount: true},
	}

	for _, path := range args.VisiblePaths {
		if !strings.HasPrefix(path, limitDir) {
			return fmt.Errorf("Expected visible path to start with '%s', but was '%s'", limitDir, path)
		}
		mounts = append(mounts, mountPoint{Path: path[1:], Writeable: true})
	}

	if args.Posix {
		// Example: calling out to git
		mounts = append(mounts, []mountPoint{
			{Path: "etc", NonMount: true},
			{Path: "etc/passwd", File: true, NonMount: true, Content: "root:x:0:0:/home/kapp-controller:/usr/sbin/nologin\n"},
			{Path: "bin"}, // e.g. bin/sh was referenced by git executables
			{Path: "lib"},
			{Path: "lib64"}, // some dynamic libraries
			{Path: "var/lib"},
			{Path: "usr/lib"},
			{Path: "usr/libexec"}, // internal git binaries
			{Path: "usr/share"},   // internal git data
			{Path: "usr/sbin"},
			// TODO get rid of 'dev/null'? (used by vendir today)
			{Path: "dev", NonMount: true},
			{Path: "dev/null", NonMount: true, File: true},
		}...)
	}

	if args.Network {
		mounts = append(mounts, []mountPoint{
			// Cannot mount whole 'etc' directory directly (mount returns exit status 32)
			// Mounting individual child files/directories makes it work.
			{Path: "etc/pki"},
			{Path: "etc/resolv.conf", File: true},
			{Path: "etc/hostname", File: true},
			{Path: "etc/hosts", File: true},
		}...)
	}

	err := w.configureMounts(mounts, args.NewRootDir)
	if err != nil {
		return err
	}

	err = w.pivotRoot(args.NewRootDir)
	if err != nil {
		return err
	}

	// Not all command specify execution directory
	if len(args.CmdDir) > 0 {
		err = os.Chdir(args.CmdDir)
		if err != nil {
			return fmt.Errorf("Chdir '%s': %s", args.CmdDir, err)
		}
	}

	return syscall.Exec(args.CmdPath, append([]string{args.CmdPath}, args.CmdArgs...), os.Environ())
}

func (SandboxWrap) configureMounts(mounts []mountPoint, newRootDir string) error {
	for _, mount := range mounts {
		if mount.File {
			file, err := os.Create(filepath.Join(newRootDir, mount.Path))
			if err != nil {
				return fmt.Errorf("Creating mount file '%s': %s", mount.Path, err)
			}
			if len(mount.Content) > 0 {
				_, err = file.Write([]byte(mount.Content))
				if err != nil {
					return fmt.Errorf("Writing into file '%s': %s", mount.Path, err)
				}
			}
			err = file.Close()
			if err != nil {
				return fmt.Errorf("Closing mount file '%s': %s", mount.Path, err)
			}
		} else {
			err := os.MkdirAll(filepath.Join(newRootDir, mount.Path), 0700)
			if err != nil {
				return fmt.Errorf("Creating mount dir '%s': %s", mount.Path, err)
			}
		}
	}

	for _, mount := range mounts {
		if mount.NonMount {
			continue
		}

		srcPath := string(filepath.Separator) + mount.Path
		dstPath := filepath.Join(newRootDir, mount.Path)

		// Equivalent to: mount --bind /usr/sbin usr/sbin
		err := syscall.Mount(srcPath, dstPath, "bind", syscall.MS_BIND, "")
		if err != nil {
			return fmt.Errorf("Mounting path '%s': %s", mount.Path, err)
		}

		// TODO readonly for files?
		if !mount.Writeable && !mount.File {
			// Equivalent to: mount -o remount,ro,bind etc
			err = syscall.Mount("", dstPath, "", syscall.MS_REMOUNT|syscall.MS_BIND|syscall.MS_RDONLY, "")
			if err != nil {
				return fmt.Errorf("Remounting path '%s' as readonly: %s", mount.Path, err)
			}
		}
	}

	return nil
}

func (SandboxWrap) pivotRoot(newRootDir string) error {
	oldRootDir := filepath.Join(newRootDir, "/.pivot_root")

	// Satisfy requirement of pivot_root (https://man7.org/linux/man-pages/man2/pivot_root.2.html)
	err := syscall.Mount(newRootDir, newRootDir, "", syscall.MS_BIND|syscall.MS_REC, "")
	if err != nil {
		return fmt.Errorf("Remounting new root dir: %s", err)
	}

	err = os.MkdirAll(oldRootDir, 0700)
	if err != nil {
		return fmt.Errorf("Creating old root dir: %s", err)
	}

	err = syscall.PivotRoot(newRootDir, oldRootDir)
	if err != nil {
		return fmt.Errorf("Pivot root: %s", err)
	}

	err = os.Chdir("/")
	if err != nil {
		return fmt.Errorf("Chdir /: %s", err)
	}

	oldRootDir = "/.pivot_root"

	err = syscall.Unmount(oldRootDir, syscall.MNT_DETACH)
	if err != nil {
		return fmt.Errorf("Unmounting old root: %s", err)
	}

	err = os.RemoveAll(oldRootDir)
	if err != nil {
		return fmt.Errorf("Removing old root dir: %s", err)
	}

	return nil
}
