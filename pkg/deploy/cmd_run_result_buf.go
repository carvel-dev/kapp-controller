// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package deploy

import (
	"sync"

	"github.com/vmware-tanzu/carvel-kapp-controller/pkg/exec"
)

type WriterFunc func([]byte) (int, error)

func (f WriterFunc) Write(data []byte) (int, error) {
	return f(data)
}

type CmdRunResultBuffer struct {
	result     *exec.CmdRunResult
	resultLock *sync.RWMutex
}

func NewCmdRunResultBuffer() *CmdRunResultBuffer {
	return &CmdRunResultBuffer{result: &exec.CmdRunResult{}, resultLock: &sync.RWMutex{}}
}

func (w *CmdRunResultBuffer) WriteStdout(data []byte) (int, error) {
	w.resultLock.Lock()
	w.result.Stdout += string(data)
	w.resultLock.Unlock()
	return len(data), nil
}

func (w *CmdRunResultBuffer) WriteStderr(data []byte) (int, error) {
	w.resultLock.Lock()
	w.result.Stderr += string(data)
	w.resultLock.Unlock()
	return len(data), nil
}

func (w *CmdRunResultBuffer) Copy() exec.CmdRunResult {
	w.resultLock.RLock()
	result := *w.result
	w.resultLock.RUnlock()
	return result
}
