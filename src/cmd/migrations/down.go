package migrations

import (
	"github.com/kolesnikovm/messenger/configs"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var cmdDown = &cobra.Command{
	Use:   "down",
	Short: "Perform down migration to the first version",
	Run: func(cmd *cobra.Command, args []string) {
		configFile, err := cmd.InheritedFlags().GetString("config")
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}

		config, err := configs.NewServerConfig(configFile)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to instantiate config")
		}

		migrations, err := new(config)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create migrations")
		}

		if err := migrations.Down(); err != nil {
			log.Fatal().Err(err).Msg("failed to perform down migration")
		}
	},
}
