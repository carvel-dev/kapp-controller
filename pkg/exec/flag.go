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
		return Flag{}, fmt.Errorf("Expected flag to be in '--name=val' format")
	}
	name := match[1]
	val := match[3]
	return Flag{Name: name, Value: val}, nil
}
