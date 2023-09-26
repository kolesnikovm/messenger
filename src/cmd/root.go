package cmd

import (
	"fmt"
	"os"

	"github.com/kolesnikovm/messenger/cmd/client"
	"github.com/kolesnikovm/messenger/cmd/server"
	"github.com/kolesnikovm/messenger/configs"
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
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(
		server.Cmd,
		client.Cmd,
	)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./messenger.yaml)")
}

func initConfig() {
	if err := configs.Load(cfgFile); err != nil {
		fmt.Printf("failed to load config file: %s\n", err)
	}

	config, err := configs.New()
	if err != nil {
		fmt.Printf("failed to instantiate config file: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("using config: %+v\n", config)
}
