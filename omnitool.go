package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/codegangsta/cli"
)

//
// Commands
//

func cmdRun(c *cli.Context) {
	if len(c.Args()) != 1 {
		cli.ShowCommandHelp(c, "command")
		return
	}

	cmd := strings.Join(c.Args(), " ")

	u := c.GlobalString("username")
	k := c.GlobalString("keyfile")
	ml := c.GlobalString("list")
	mg := c.GlobalString("group")

	// Load machine list file
	mf, err := LoadFile(ml)
	if err != nil {
		fmt.Println("ERROR: Couldn't find %s", ml)
		return
	}

	// Lost machine addresses by group name
	machine_list := mf.Get(mg)

	// Receive responses on this channel
	results := make(chan SSHResponse)

	// Timeout after N seconds
	timeout := time.After(60 * time.Second)

	// Spawn the gorountines
	MapCmd(machine_list, u, k, cmd, results)

	// See how they did
	for i := 0; i < len(machine_list); i++ {
		select {
		case r := <-results:
			fmt.Printf("Hostname: %s\n", r.Hostname)
			fmt.Printf("Result:\n%s\n", r.Result)
		case <-timeout:
			fmt.Println("Timed out!")
		}
	}

	fmt.Println("CMD: ", c.Args())
}

func cmdScp(c *cli.Context) {
	if len(c.Args()) != 2 {
		cli.ShowCommandHelp(c, "scp")
		return
	}

	localPath := c.Args()[0]
	remotePath := c.Args()[1]

	u := c.GlobalString("username")
	k := c.GlobalString("keyfile")
	ml := c.GlobalString("list")
	mg := c.GlobalString("group")

	// Load machine list file
	mf, err := LoadFile(ml)
	if err != nil {
		fmt.Println("ERROR: Couldn't find %s", ml)
		return
	}

	// Lost machine addresses by group name
	machine_list := mf.Get(mg)

	// Receive responses on this channel
	results := make(chan SSHResponse)

	// Timeout after N seconds
	timeout := time.After(60 * time.Second)

	// Spawn the gorountines
	MapScp(machine_list, u, k, localPath, remotePath, results)

	// See how they did
	for i := 0; i < len(machine_list); i++ {
		select {
		case r := <-results:
			fmt.Printf("Hostname: %s\n", r.Hostname)
			fmt.Printf("Result:\n%s\n", r.Result)
		case <-timeout:
			fmt.Println("Timed out!")
		}
	}

	fmt.Println("SCP: ", c.Args())
}

func main() {
	// App setup

	app := cli.NewApp()
	app.Name = "omnitool"
	app.Usage = "Simple SSH pools, backed by machine lists"

	app.Version = "0.1"

	// Global Flags

	globalFlags := []cli.Flag{
		cli.StringFlag{
			Name:   "list, l",
			Value:  "machines.list",
			Usage:  "Path to machine list file",
			EnvVar: "OMNI_MACHINE_LIST",
		},

		cli.StringFlag{
			Name:   "username, u",
			Value:  "vagrant",
			Usage:  "Username for machine group",
			EnvVar: "OMNI_USERNAME",
		},

		cli.StringFlag{
			Name:   "keyfile, k",
			Value:  os.Getenv("HOME") + "/.vagrant.d/insecure_private_key",
			Usage:  "Path to auth key",
			EnvVar: "OMNI_KEYFILE",
		},

		cli.StringFlag{
			Name:   "group, g",
			Value:  "vagrants",
			Usage:  "Machine group to perform task on",
			EnvVar: "OMNI_MACHINE_GROUP",
		},
	}

	app.Flags = globalFlags

	// Subcommands

	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "Runs command on machine group",
			Action: cmdRun,
			Flags:  globalFlags,
		},
		{
			Name:   "scp",
			Usage:  "Copies file to machine group",
			Action: cmdScp,
			Flags:  globalFlags,
		},
	}

	// Do it

	app.Run(os.Args)
}
