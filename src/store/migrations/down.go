package migrations

import (
	"fmt"

	"github.com/pressly/goose/v3"
)

func (m *Migrations) Down() error {
	const op = "Migrations.Down"

	if err := goose.Down(m.DB, "sql"); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
