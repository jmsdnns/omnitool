package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/jmsdnns/omnitool/sessions"
	"github.com/spf13/cobra"
)

// sftpCmd represents the sftp command
var sftpCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copies file to host group",
	Long: `Creates an SFTP pool to a host group and copies file to each host. Success of
failure is then collected from each host and reported here`,
	Run: cmdCopy,
}

func init() {
	RootCmd.AddCommand(sftpCmd)
}

func cmdCopy(cmd *cobra.Command, args []string) {
	username := cmd.Flags().Lookup("username").Value.String()
	keyfile := cmd.Flags().Lookup("keyfile").Value.String()
	hostsfile := cmd.Flags().Lookup("hostsfile").Value.String()
	group := cmd.Flags().Lookup("group").Value.String()

	hostList, err := ParseHostArgs(hostsfile, group)
	if err != nil {
		log.Fatal(err)
		return
	}

	if len(args) != 2 {
		log.Fatal("localPath and remotePath arguments not found")
		return
	}
	localPath := args[0]
	remotePath := args[1]

	results := make(chan sessions.SFTPResponse)
	timeout := time.After(60 * time.Second)
	sessions.MapCopy(hostList, username, keyfile, localPath, remotePath, results)

	for i := 0; i < len(hostList); i++ {
		select {
		case r := <-results:
			fmt.Printf("Host: %s\n", r.Host)
			fmt.Printf("Result: %s\n\n", r.Result)
		case <-timeout:
			fmt.Println("Timed out!")
		}
	}
}
