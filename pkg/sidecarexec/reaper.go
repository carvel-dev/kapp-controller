// Copyright 2024 The Carvel Authors.
// SPDX-License-Identifier: Apache-2.0

package sidecarexec

import (
	"sync"
	"syscall"
	"time"

	"github.com/go-logr/logr"
)

// Reaper is a centralized zombie process reaper that eliminates the race
// condition between a blanket Wait4(-1) and Go's os/exec cmd.Wait().
// It is the sole goroutine calling Wait4, and dispatches exit statuses
// back to callers that registered interest in specific PIDs.
type Reaper struct {
	mu      sync.Mutex
	tracked map[int]chan syscall.WaitStatus
	log     logr.Logger
}

// NewReaper creates a new Reaper. Start the reaping loop by calling Run
// in a separate goroutine.
func NewReaper(log logr.Logger) *Reaper {
	return &Reaper{
		tracked: make(map[int]chan syscall.WaitStatus),
		log:     log.WithName("reaper"),
	}
}

// Lock acquires the reaper's mutex. This must be held while starting a
// process and calling TrackLocked, so that the reaping loop cannot
// consume the process exit status before tracking is registered.
func (r *Reaper) Lock() { r.mu.Lock() }

// Unlock releases the reaper's mutex.
func (r *Reaper) Unlock() { r.mu.Unlock() }

// TrackLocked registers a PID so that when Wait4 reaps it, the wait status
// is delivered to the returned channel. The caller MUST hold r.Lock() and
// must have started the process while holding the lock, ensuring no race
// between process start and tracking registration.
func (r *Reaper) TrackLocked(pid int) <-chan syscall.WaitStatus {
	ch := make(chan syscall.WaitStatus, 1)
	r.tracked[pid] = ch
	return ch
}

// dispatch delivers a wait status to a tracked PID's channel.
// Returns true if the PID was being tracked (an RPC-managed process).
func (r *Reaper) dispatch(pid int, status syscall.WaitStatus) bool {
	r.mu.Lock()
	ch, ok := r.tracked[pid]
	if ok {
		delete(r.tracked, pid)
	}
	r.mu.Unlock()

	if ok {
		ch <- status
	}
	return ok
}

// Run is the reaping loop and must be the only goroutine calling Wait4
// in this process. It runs forever.
func (r *Reaper) Run() {
	r.log.Info("starting zombie reaper")

	for {
		var status syscall.WaitStatus

		// The mutex is acquired here so that Wait4 cannot race with
		// a concurrent cmd.Start() + TrackLocked() sequence.
		r.mu.Lock()
		pid, err := syscall.Wait4(-1, &status, syscall.WNOHANG, nil)

		if pid > 0 {
			ch, ok := r.tracked[pid]
			if ok {
				delete(r.tracked, pid)
			}
			r.mu.Unlock()

			if ok {
				ch <- status
				r.log.V(1).Info("reaped tracked process", "pid", pid)
			} else {
				r.log.V(1).Info("reaped orphaned zombie", "pid", pid)
			}
			continue
		}
		r.mu.Unlock()

		if pid == 0 || (pid == -1 && err == syscall.ECHILD) {
			time.Sleep(1 * time.Second)
			continue
		}

		r.log.V(1).Info("wait4 unexpected error", "err", err)
		time.Sleep(1 * time.Second)
	}
}
