// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package sidecarexec_test

import (
	"sync"
	"syscall"
	"testing"
	"time"

	"carvel.dev/kapp-controller/pkg/sidecarexec"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Reaper_TrackLockedDeliversStatus(t *testing.T) {
	reaper := sidecarexec.NewReaper(logr.Discard())

	reaper.Lock()
	ch := reaper.TrackLocked(9999)
	reaper.Unlock()

	go func() {
		time.Sleep(50 * time.Millisecond)
		sidecarexec.Dispatch_ForTesting(reaper, 9999, syscall.WaitStatus(0))
	}()

	select {
	case status := <-ch:
		assert.Equal(t, syscall.WaitStatus(0), status)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for status delivery")
	}
}

func Test_Reaper_UntrackedPidIsDiscarded(t *testing.T) {
	reaper := sidecarexec.NewReaper(logr.Discard())

	reaper.Lock()
	ch := reaper.TrackLocked(1111)
	reaper.Unlock()

	dispatched := sidecarexec.Dispatch_ForTesting(reaper, 2222, syscall.WaitStatus(0))
	assert.False(t, dispatched, "untracked PID should not be dispatched")

	select {
	case <-ch:
		t.Fatal("channel for PID 1111 should not have received anything")
	default:
	}
}

func Test_Reaper_TrackIsRemovedAfterDispatch(t *testing.T) {
	reaper := sidecarexec.NewReaper(logr.Discard())

	reaper.Lock()
	ch := reaper.TrackLocked(3333)
	reaper.Unlock()

	ok := sidecarexec.Dispatch_ForTesting(reaper, 3333, syscall.WaitStatus(0))
	require.True(t, ok)
	<-ch

	ok2 := sidecarexec.Dispatch_ForTesting(reaper, 3333, syscall.WaitStatus(0))
	assert.False(t, ok2, "second dispatch to same PID should return false")
}

func Test_Reaper_MultipleTrackedPidsDispatchIndependently(t *testing.T) {
	reaper := sidecarexec.NewReaper(logr.Discard())

	reaper.Lock()
	ch1 := reaper.TrackLocked(100)
	ch2 := reaper.TrackLocked(200)
	ch3 := reaper.TrackLocked(300)
	reaper.Unlock()

	exitOk := syscall.WaitStatus(0)
	// Simulate an exit status with code 1 (on Linux, WaitStatus stores
	// the exit code in bits 8-15, so code<<8 gives exit status N).
	exitFail := syscall.WaitStatus(1 << 8)

	sidecarexec.Dispatch_ForTesting(reaper, 200, exitFail)
	sidecarexec.Dispatch_ForTesting(reaper, 100, exitOk)
	sidecarexec.Dispatch_ForTesting(reaper, 300, exitOk)

	assert.Equal(t, exitOk, <-ch1)
	assert.Equal(t, exitFail, <-ch2)
	assert.Equal(t, exitOk, <-ch3)
}

func Test_Reaper_ConcurrentDispatchIsSafe(t *testing.T) {
	reaper := sidecarexec.NewReaper(logr.Discard())

	const n = 50
	channels := make([]<-chan syscall.WaitStatus, n)

	reaper.Lock()
	for i := 0; i < n; i++ {
		channels[i] = reaper.TrackLocked(i + 1)
	}
	reaper.Unlock()

	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(pid int) {
			defer wg.Done()
			sidecarexec.Dispatch_ForTesting(reaper, pid, syscall.WaitStatus(0))
		}(i + 1)
	}
	wg.Wait()

	for i := 0; i < n; i++ {
		select {
		case status := <-channels[i]:
			assert.Equal(t, syscall.WaitStatus(0), status, "PID %d", i+1)
		case <-time.After(2 * time.Second):
			t.Fatalf("timed out waiting for PID %d", i+1)
		}
	}
}

func Test_CmdExec_DisallowedCommandRejected(t *testing.T) {
	reaper := sidecarexec.NewReaper(logr.Discard())

	cmdExec := sidecarexec.CmdExec_ForTesting(reaper, []string{"echo"})

	var output sidecarexec.CmdOutput
	err := cmdExec.Run(sidecarexec.CmdInput{
		Command: "rm",
		Args:    []string{"-rf", "/"},
	}, &output)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not allowed")
}

func Test_Reaper_NonZeroExitStatusDelivered(t *testing.T) {
	reaper := sidecarexec.NewReaper(logr.Discard())

	reaper.Lock()
	ch := reaper.TrackLocked(4444)
	reaper.Unlock()

	exitCode42 := syscall.WaitStatus(42 << 8)
	sidecarexec.Dispatch_ForTesting(reaper, 4444, exitCode42)

	status := <-ch
	assert.True(t, status.Exited())
	assert.Equal(t, 42, status.ExitStatus())
}

func Test_Reaper_LockPreventsDispatchDuringTrack(t *testing.T) {
	reaper := sidecarexec.NewReaper(logr.Discard())

	reaper.Lock()

	dispatched := make(chan bool, 1)
	go func() {
		// dispatch acquires the mutex internally, so it will block
		// until we release the lock
		ok := sidecarexec.Dispatch_ForTesting(reaper, 5555, syscall.WaitStatus(0))
		dispatched <- ok
	}()

	// Give the goroutine time to reach the mutex
	time.Sleep(100 * time.Millisecond)

	select {
	case <-dispatched:
		t.Fatal("dispatch should be blocked while lock is held")
	default:
	}

	ch := reaper.TrackLocked(5555)
	reaper.Unlock()

	ok := <-dispatched
	require.True(t, ok, "dispatch should succeed after unlock with tracking in place")

	status := <-ch
	assert.Equal(t, syscall.WaitStatus(0), status)
}
