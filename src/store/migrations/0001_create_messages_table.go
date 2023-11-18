package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up0001, Down0001)
}

func Up0001(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `create table if not exists messages (
		id              uuid            primary key not null,
		sender_id       bigint          not null,
		chat_id         varchar(255)    not null,
		text            text            not null
	);`)
	if err != nil {
		return err
	}

	return nil
}

func Down0001(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `drop table if exists messages;`)
	if err != nil {
		return err
	}

	return nil
}
