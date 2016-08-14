package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmsdnns/omnitool/sessions"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a command on host group",
	Long: `Creates an SSH pool to a host group and executes a command. The output is then
collected from each host and displayed here.`,
	Run: cmdRun,
}

func init() {
	RootCmd.AddCommand(runCmd)
}

func cmdRun(cmd *cobra.Command, args []string) {
	username := cmd.Flags().Lookup("username").Value.String()
	keyfile := cmd.Flags().Lookup("keyfile").Value.String()
	hostsfile := cmd.Flags().Lookup("hostsfile").Value.String()
	group := cmd.Flags().Lookup("group").Value.String()

	hostList, err := ParseHostArgs(hostsfile, group)
	if err != nil {
		log.Fatal(err)
		return
	}

	userCmd := strings.Join(args, " ")

	results := make(chan sessions.SSHResponse)
	timeout := time.After(60 * time.Second)
	sessions.MapCmd(hostList, username, keyfile, userCmd, results)

	fmt.Printf("CMD: %s\n\n", userCmd)
	for i := 0; i < len(hostList); i++ {
		select {
		case r := <-results:
			fmt.Printf("Host: %s\n", r.Host)
			fmt.Printf("Result:\n%s\n", r.Result)
		case <-timeout:
			fmt.Println("Timed out!")
		}
	}

}
