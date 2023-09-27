package client

import (
	"github.com/kolesnikovm/messenger/configs"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "client",
	Short: "Start messenger in client mode",
	Run: func(cmd *cobra.Command, args []string) {
		configFile, err := cmd.InheritedFlags().GetString("config")
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}

		// check config is ok
		config, err := configs.NewClientConfig(configFile)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to instantiate config")
		}

		log.Info().Msgf("Messenger client connected to %+v", config.ServerAddress)
	},
}
