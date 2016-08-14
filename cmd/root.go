package cmd

import (
	"fmt"
	"os"

	"github.com/jmsdnns/omnitool/hosts"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "omnitool",
	Short: "A tool for managing machines via parallel SSH pools",
	Long: `Omnitool's goal is to let you think in terms of one machine while working with
N machines.`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringP("username", "u", "", "username for ssh")
	RootCmd.PersistentFlags().StringP("keyfile", "k", "", "path to ssh key")
	RootCmd.PersistentFlags().StringP("hostsfile", "", "hosts.list", "path to hosts file")
	RootCmd.PersistentFlags().StringP("group", "g", "", "host group for task")
}

// ParseHostArgs reads a hosts file, finds the requested group, and returns
// the corresponding list of hosts
func ParseHostArgs(hostsfile string, group string) ([]string, error) {
	// Load machine list file
	hostsConfig, err := hosts.LoadHostsFile(hostsfile)
	if err != nil {
		return nil, err
	}

	hostGroup := hostsConfig.Get(group)
	return hostGroup, nil
}
