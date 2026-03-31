// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package sidecarexec

import "syscall"

// CmdExec_ForTesting creates a CmdExec with the given reaper and allowed
// command names. This is only available in test builds.
func CmdExec_ForTesting(reaper *Reaper, allowedCmds []string) CmdExec {
	allowed := map[string]struct{}{}
	for _, cmd := range allowedCmds {
		allowed[cmd] = struct{}{}
	}
	return CmdExec{reaper: reaper, allowedCmdNames: allowed}
}

// Dispatch_ForTesting exposes the reaper's dispatch method for unit tests.
func Dispatch_ForTesting(r *Reaper, pid int, status syscall.WaitStatus) bool {
	return r.dispatch(pid, status)
}
