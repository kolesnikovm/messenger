package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up0002, Down0002)
}

func Up0002(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `create table if not exists last_read_message (
		id                    serial          primary key,
		chat_id               varchar(255)    not null,
		user_id               bigint          not null,
		message_order_id      bigint          not null
	);`)
	if err != nil {
		return err
	}

	return nil
}

func Down0002(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `drop table if exists last_read_message;`)
	if err != nil {
		return err
	}

	return nil
}
