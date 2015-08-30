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
//
// The file's syntax is close to ini but not quite. I think of it as an ini file
// for machines or .mini files.
//
// Example use:
//
//	f, _ := LoadFile("test.mini")
//
//	for groupName, group := range f {
//		fmt.Println("Group: ", groupName)
//		for _, hostname := range group {
//			fmt.Println("  - ", hostname)
//		}
//	}
//
//	g := f.Get("dbservers")
//	fmt.Println("Group: ", g)
//
// This code was adapted from: https://github.com/vaughan0/go-ini

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

//
// Machine Groups
//

// group name => group
type File map[string]HostGroup

// group == list of hosts
type HostGroup []string

func (f File) Get(groupName string) HostGroup {
	group := f[groupName]
	if group == nil {
		group = make(HostGroup, 0)
		f[groupName] = group
	}

	return group
}

func (f File) Set(groupName string, group HostGroup) {
	f[groupName] = group
}

func (f File) Load(in io.Reader) (err error) {
	bufin, ok := in.(*bufio.Reader)
	if !ok {
		bufin = bufio.NewReader(in)
	}

	return parseFile(bufin, f)
}

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

func Load(in io.Reader) (File, error) {
	file := make(File)
	err := file.Load(in)
	return file, err
}

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
