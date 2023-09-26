package client

import (
	"fmt"

	"github.com/kolesnikovm/messenger/configs"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "client",
	Short: "Start messenger in client mode",
	Run: func(cmd *cobra.Command, args []string) {
		conf, _ := configs.New()
		fmt.Printf("Messenger client connected to %s\n", conf.Client.ServerAddress)
	},
}
