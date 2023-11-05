package migrations

import (
	"github.com/kolesnikovm/messenger/configs"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var cmdUp = &cobra.Command{
	Use:   "up",
	Short: "Perform up migration to the latest version",
	Run: func(cmd *cobra.Command, args []string) {
		configFile, err := cmd.InheritedFlags().GetString("config")
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		config, err := configs.NewServerConfig(configFile)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to instantiate config")
		}

		migrations, err := new(config)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create migrations")
		}

		if err := migrations.Up(); err != nil {
			log.Fatal().Err(err).Msg("failed to perform up migration")
		}
	},
}
