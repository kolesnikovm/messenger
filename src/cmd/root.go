package cmd

import (
	"github.com/kolesnikovm/messenger/cmd/client"
	"github.com/kolesnikovm/messenger/cmd/migrations"
	"github.com/kolesnikovm/messenger/cmd/server"
	"github.com/rs/zerolog/log"
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
		log.Fatal().Err(err).Msg("failed to execute root cmd")
	}
}

func init() {
	rootCmd.AddCommand(
		server.Cmd,
		client.Cmd,
		migrations.Cmd,
	)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./messenger.yaml)")
}
