package migrations

import (
	"fmt"

	"github.com/pressly/goose/v3"
)

func (m *Migrations) Up() error {
	const op = "Migrations.Up"

	if err := goose.Up(m.DB, "sql"); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
