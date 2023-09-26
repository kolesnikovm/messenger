package cmd

import (
	"fmt"
	"os"

	"github.com/kolesnikovm/messenger/cmd/client"
	"github.com/kolesnikovm/messenger/cmd/server"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:  "messenger",
	Long: `Messenger application allows you to communicate with other clients via Messenger server.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("failed to execute root cmd: %s", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(
		server.Cmd,
		client.Cmd,
	)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./messenger.yaml)")
}
