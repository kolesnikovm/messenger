package migrations

import (
	"fmt"

	"github.com/kolesnikovm/messenger/configs"
	store "github.com/kolesnikovm/messenger/store/migrations"
	"github.com/kolesnikovm/messenger/store/postgres"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "migrations",
	Short: "Perform database migrations",
}

func init() {
	Cmd.AddCommand(
		cmdUp,
		cmdDown,
	)
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
