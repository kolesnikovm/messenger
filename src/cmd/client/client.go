package client

import (
	"fmt"
	"os"

	"github.com/kolesnikovm/messenger/configs"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "client",
	Short: "Start messenger in client mode",
	Run: func(cmd *cobra.Command, args []string) {
		configFile, err := cmd.InheritedFlags().GetString("config")
		if err != nil {
			fmt.Printf("%s\n", err)
		}

		// check config is ok
		config, err := configs.NewClientConfig(configFile)
		if err != nil {
			fmt.Printf("failed to instantiate config: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Messenger client connected to %s\n", config.ServerAddress)
	},
}
