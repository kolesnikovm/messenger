package server

import (
	"fmt"

	"github.com/kolesnikovm/messanger/configs"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "server",
	Short: "Start messanger in server mode",
	Run: func(cmd *cobra.Command, args []string) {
		conf, _ := configs.New()
		fmt.Printf("Messanger server listening on %s\n", conf.Server.ListenAddress)
	},
}
