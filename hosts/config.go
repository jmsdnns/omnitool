package hosts

// This an interface for reading ini-like host groups from file. The format is
// simply a group name with hostnames on each line below, until a new group
// starts or EOF.
//
// Example:
//
//     [apiservers]
//     10.0.0.1
//     10.0.0.2
//     10.0.0.3
//     10.0.0.4
//
//     [dbservers]
//     10.0.0.5
//     10.0.0.6
//
// This file is stored as a map[group name] => list of hosts

import (
	"bufio"
	"os"
)

// Config is a map of group name => group
type Config map[string]Group

// Group is a list of hosts
type Group []string

// Get takes a group name and returns a HostGroup
func (c Config) Get(groupName string) Group {
	group := c[groupName]
	if group == nil {
		group = make(Group, 0)
		c[groupName] = group
	}

	return group
}

// Set associates a HostGroup with a group name
func (c Config) Set(groupName string, group Group) {
	c[groupName] = group
}

// LoadFile opens a file by name and passes io.Reader to File.Load
func (c Config) LoadFile(file string) (err error) {
	in, err := os.Open(file)
	if err != nil {
		return
	}
	defer in.Close()

	bufin := bufio.NewReader(in)

	err = parseFile(bufin, &c)
	return err
}
