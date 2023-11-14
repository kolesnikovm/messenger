package migrations

import (
	"fmt"

	"github.com/kolesnikovm/messenger/configs"
	store "github.com/kolesnikovm/messenger/store/migrations"
	"github.com/kolesnikovm/messenger/store/postgres"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "migrations [goose-command]",
	Short: "Perform database migrations",
	Args:  cobra.MinimumNArgs(1),
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
		defer migrations.Close()

		arguments := []string{}
		if len(args) > 1 {
			arguments = append(arguments, args[1:]...)
		}

		if err := migrations.Run(args[0], arguments...); err != nil {
			log.Fatal().Err(err).Msgf("failed to perform migration command: %s", args[0])
		}
	},
}

func new(config configs.ServerConfig) (*store.Migrations, error) {
	const op = "migrations.new"

	db, err := postgres.New(config.Postgres)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	migrations, err := store.New(db)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return migrations, nil
}
