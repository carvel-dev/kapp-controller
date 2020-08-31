// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"fmt"
	"regexp"
)

var (
	// Be very strict about what kind of format is allowed for flags
	fullFlagRegex = regexp.MustCompile(`\A(\-\-[a-z]+(\-[a-z]+)*)=(.+)\z`)
)

type Flag struct {
	Name  string // e.g. --name
	Value string
}

func NewFlagFromString(str string) (Flag, error) {
	match := fullFlagRegex.FindStringSubmatch(str)
	if len(match) != 4 {
		return Flag{}, fmt.Errorf("Expected flag '%s' to be in '--name=val' format", str)
	}
	name := match[1]
	val := match[3]
	return Flag{Name: name, Value: val}, nil
}

type FlagSet struct {
	included map[string]struct{}
}

func NewFlagSet(optss ...[]string) FlagSet {
	result := map[string]struct{}{}
	for _, opts := range optss {
		for _, opt := range opts {
			result[opt] = struct{}{}
		}
	}
	return FlagSet{result}
}

func (s FlagSet) Includes(name string) bool {
	_, found := s.included[name]
	return found
}
