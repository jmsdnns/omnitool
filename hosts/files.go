package hosts

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

//
// File Parsing
//

// LoadHostsFile reads a file into a File instance
func LoadHostsFile(filename string) (Config, error) {
	config := make(Config)
	err := config.LoadFile(filename)
	return config, err
}

func parseFile(in *bufio.Reader, c *Config) (err error) {
	// [group name]
	groupNamePattern := regexp.MustCompile(`^\[(.*)\]$`)

	// Tracks last group name found
	currentGroup := ""

	// Walk across each line input
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

		// Group names are the only pattern we care about in the file. Anything that
		// isn't a group name should be treated like a host address
		if groups := groupNamePattern.FindStringSubmatch(line); groups != nil {
			name := strings.TrimSpace(groups[1])
			currentGroup = name
			c.Get(name)
		} else {
			host := line
			group := c.Get(currentGroup)
			group = append(group, host)
			c.Set(currentGroup, group)
		}
	}

	return nil
}
