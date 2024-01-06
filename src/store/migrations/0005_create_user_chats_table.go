package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up0005, Down0005)
}

func Up0005(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `create table if not exists user_chats (
		user_id         bigint          not null,
		chat_id         varchar(255)    not null,
		primary key(user_id, chat_id)
	);`)
	if err != nil {
		return err
	}

	return nil
}

func Down0005(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `drop table if exists user_chats;`)
	if err != nil {
		return err
	}

	return nil
}
