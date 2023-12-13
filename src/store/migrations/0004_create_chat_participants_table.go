package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up0004, Down0004)
}

func Up0004(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `create table if not exists chat_participants (
		user_id                bigint          not null,
		chat_id                varchar(255)    not null,
		last_read_message      uuid            not null,
		primary key(user_id, chat_id),
		foreign key(chat_id) references chats(id)
	);`)
	if err != nil {
		return err
	}

	return nil
}

func Down0004(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `drop table if exists chat_participants;`)
	if err != nil {
		return err
	}

	return nil
}
