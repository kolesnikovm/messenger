package server

import (
	"fmt"
	"os"

	"github.com/kolesnikovm/messenger/configs"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "server",
	Short: "Start messenger in server mode",
	Run: func(cmd *cobra.Command, args []string) {
		configFile, err := cmd.InheritedFlags().GetString("config")
		if err != nil {
			fmt.Printf("%s\n", err)
		}

		// check config is ok
		config, err := configs.NewServerConfig(configFile)
		if err != nil {
			fmt.Printf("failed to instantiate config: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Messenger server listening on %s\n", config.ListenAddress)
	},
}
