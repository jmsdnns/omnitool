package main

// The purpose of this file is to provide an interface for reading ini-like
// files; .mini files. The format is simply a group name with hostnames on each
// line below, until a new group starts or the file is finished.
//
// Example .mini file:
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

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

// File is a map of group name => group
type File map[string]HostGroup

// HostGroup is a list of hosts
type HostGroup []string

// Get takes a group name and returns a HostGroup
func (f File) Get(groupName string) HostGroup {
	group := f[groupName]
	if group == nil {
		group = make(HostGroup, 0)
		f[groupName] = group
	}

	return group
}

// Set associates a HostGroup with a group name
func (f File) Set(groupName string, group HostGroup) {
	f[groupName] = group
}

// Load parses an io.Reader and populates itself
func (f File) Load(in io.Reader) (err error) {
	bufin, ok := in.(*bufio.Reader)
	if !ok {
		bufin = bufio.NewReader(in)
	}

	return parseFile(bufin, f)
}

// LoadFile opens a file by name and passes io.Reader to File.Load
func (f File) LoadFile(file string) (err error) {
	in, err := os.Open(file)
	if err != nil {
		return
	}
	defer in.Close()

	return f.Load(in)
}

//
// Config Parsing
//

// Load parses an io.Reader into a File instance
func Load(in io.Reader) (File, error) {
	file := make(File)
	err := file.Load(in)
	return file, err
}

// LoadFile reads a file into a File instance
func LoadFile(filename string) (File, error) {
	file := make(File)
	err := file.LoadFile(filename)
	return file, err
}

func parseFile(in *bufio.Reader, file File) (err error) {
	// [group name]
	groupNamePattern := regexp.MustCompile(`^\[(.*)\]$`)

	// Tracks last group found
	currentGroup := ""

	for done := false; !done; {
		var line string
		if line, err = in.ReadString('\n'); err != nil {
			if err == io.EOF {
				done = true
				continue
			} else {
				return err
			}
		}

		line = strings.TrimSpace(line)

		// Skip blank lines
		if len(line) == 0 {
			continue
		}

		// Skip comments
		if line[0] == ';' || line[0] == '#' {
			continue
		}

		// Group names are the only pattern in the file. Anything that
		// isn't a group name should be treated like a host name
		if groups := groupNamePattern.FindStringSubmatch(line); groups != nil {
			name := strings.TrimSpace(groups[1])
			currentGroup = name
			file.Get(name)
		} else {
			host := line
			group := file.Get(currentGroup)
			group = append(group, host)
			file.Set(currentGroup, group)
		}
	}

	return nil
}
