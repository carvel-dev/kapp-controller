package fetch

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	"github.com/k14s/kapp-controller/pkg/memdir"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Git struct {
	opts       v1alpha1.AppFetchGit
	nsName     string
	coreClient kubernetes.Interface
}

func NewGit(opts v1alpha1.AppFetchGit, nsName string, coreClient kubernetes.Interface) *Git {
	return &Git{opts, nsName, coreClient}
}

func (t *Git) Retrieve(dstPath string) error {
	if len(t.opts.URL) == 0 {
		return fmt.Errorf("Expected non-empty URL")
	}
	if len(t.opts.Ref) == 0 {
		return fmt.Errorf("Expected non-empty ref (could be branch, tag, commit)")
	}

	authOpts, err := t.getAuthOpts()
	if err != nil {
		return err
	}

	sshAuthDir := memdir.NewTmpDir("fetch-git")

	err = sshAuthDir.Create()
	if err != nil {
		return err
	}

	defer sshAuthDir.Remove()

	sshCmd := []string{"-o", "ServerAliveInterval=30", "-o", "ForwardAgent=no", "-F", "/dev/null"}

	if authOpts.PrivateKey != nil {
		path := filepath.Join(sshAuthDir.Path(), "private-key")

		err = ioutil.WriteFile(path, []byte(*authOpts.PrivateKey), 0600)
		if err != nil {
			return fmt.Errorf("Writing private key: %s", err)
		}

		sshCmd = append(sshCmd, "-i", path)
	}

	if authOpts.KnownHosts != nil {
		path := filepath.Join(sshAuthDir.Path(), "known-hosts")

		err = ioutil.WriteFile(path, []byte(*authOpts.KnownHosts), 0600)
		if err != nil {
			return fmt.Errorf("Writing known hosts: %s", err)
		}

		sshCmd = append(sshCmd, "-o", "StrictHostKeyChecking=yes", "-o", "UserKnownHostsFile="+path)
	} else {
		sshCmd = append(sshCmd, "-o", "StrictHostKeyChecking=no")
	}

	env := append(os.Environ(), "GIT_SSH_COMMAND="+strings.Join(sshCmd, " "))

	if t.opts.LFSSkipSmudge {
		env = append(env, "GIT_LFS_SKIP_SMUDGE=1")
	}

	argss := [][]string{
		{"init"},
		{"remote", "add", "origin", t.opts.URL},
		{"fetch", "origin"}, // TODO shallow clones?
		{"checkout", t.opts.Ref, "--recurse-submodules", "."},
	}

	for _, args := range argss {
		var stdoutBs, stderrBs bytes.Buffer

		cmd := exec.Command("git", args...)
		cmd.Env = env
		cmd.Dir = dstPath
		cmd.Stdout = &stdoutBs
		cmd.Stderr = &stderrBs

		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("Git %s: %s (stderr: %s)", args, err, stderrBs.String())
		}
	}

	return nil
}

type gitAuthOpts struct {
	PrivateKey *string
	KnownHosts *string
}

func (t *Git) getAuthOpts() (gitAuthOpts, error) {
	var opts gitAuthOpts

	if t.opts.SecretRef != nil {
		secret, err := t.coreClient.CoreV1().Secrets(t.nsName).Get(t.opts.SecretRef.Name, metav1.GetOptions{})
		if err != nil {
			return opts, err
		}

		for name, val := range secret.Data {
			switch name {
			case corev1.SSHAuthPrivateKey:
				key := string(val)
				opts.PrivateKey = &key
			case "ssh-knownhosts":
				hosts := string(val)
				opts.KnownHosts = &hosts
			default:
				return opts, fmt.Errorf("Unknown secret field '%s' in secret '%s'", name, secret.Name)
			}
		}
	}

	return opts, nil
}
