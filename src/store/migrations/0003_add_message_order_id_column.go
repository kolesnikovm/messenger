package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up0003, Down0003)
}

func Up0003(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `alter table messages add column order_id bigint;`)
	if err != nil {
		return err
	}

	return nil
}

func Down0003(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `alter table messages drop column order_id;`)
	if err != nil {
		return err
	}

	return nil
}
