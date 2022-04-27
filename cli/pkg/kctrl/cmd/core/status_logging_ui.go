// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/cppforlife/color"
	"github.com/cppforlife/go-cli-ui/ui"
	"github.com/k14s/difflib"
	"k8s.io/apimachinery/pkg/util/duration"
)

type StatusLoggingUI struct {
	ui ui.UI
}

func NewStatusLoggingUI(ui ui.UI) StatusLoggingUI {
	return StatusLoggingUI{ui}
}

func (s StatusLoggingUI) PrintMessage(message string) {
	s.ui.BeginLinef("%s: %s\n", time.Now().Format("3:04:05PM"), message)
}

func (s StatusLoggingUI) PrintMessagef(message string, args ...interface{}) {
	message = fmt.Sprintf(message, args...)
	s.ui.BeginLinef("%s: %s\n", time.Now().Format("3:04:05PM"), message)
}

func (s StatusLoggingUI) PrintLogLine(message string, messageBlock string, errorBlock bool, timestamp time.Time) {
	messageAge := ""
	if time.Since(timestamp) > 1*time.Second {
		messageAge = fmt.Sprintf("(%s ago)", duration.ShortHumanDuration(time.Since(timestamp)))
	}
	s.ui.BeginLinef("%s: %s %s\n", timestamp.Local().Format("3:04:05PM"), message, messageAge)
	if len(messageBlock) > 0 {
		s.ui.PrintBlock([]byte(s.indentMessageBlock(messageBlock, errorBlock)))
	}
}

func (s StatusLoggingUI) indentMessageBlock(messageBlock string, errored bool) string {
	lines := strings.Split(messageBlock, "\n")
	for ind := range lines {
		if errored {
			lines[ind] = color.RedString(lines[ind])
		}
		lines[ind] = fmt.Sprintf("\t    | %s", lines[ind])
	}

	indentedBlock := strings.Join(lines, "\n")
	if strings.LastIndex(indentedBlock, "\n") != (len(indentedBlock) - 1) {
		indentedBlock += "\n"
	}
	return indentedBlock
}

func (s StatusLoggingUI) PrintMessageBlockDiff(oldBlock string, newBlock string, timestamp time.Time) {
	oldLines := strings.Split(oldBlock, "\n")
	newLines := strings.Split(newBlock, "\n")
	diff := difflib.Diff(oldLines, newLines)

	var lines []string
	for _, diffLine := range diff {
		switch diffLine.Delta {
		case difflib.RightOnly:
			lines = append(lines, diffLine.Payload)
		}
	}
	if len(lines) > 0 {
		for _, line := range lines {
			s.ui.BeginLinef("\t    | %s\n", line)
			time.Sleep(10 * time.Millisecond)
		}
	}
}
