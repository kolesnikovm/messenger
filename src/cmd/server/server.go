package server

import (
	"fmt"

	"github.com/kolesnikovm/messenger/configs"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "server",
	Short: "Start messenger in server mode",
	Run: func(cmd *cobra.Command, args []string) {
		conf, _ := configs.New()
		fmt.Printf("Messenger server listening on %s\n", conf.Server.ListenAddress)
	},
}
