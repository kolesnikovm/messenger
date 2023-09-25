package client

import (
	"fmt"

	"github.com/kolesnikovm/messanger/configs"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "client",
	Short: "Start messanger in client mode",
	Run: func(cmd *cobra.Command, args []string) {
		conf, _ := configs.New()
		fmt.Printf("Messanger client connected to %s\n", conf.Client.ServerAddress)
	},
}
