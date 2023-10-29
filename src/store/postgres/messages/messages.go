package messages

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/entity"
	"github.com/kolesnikovm/messenger/store/postgres"
	"github.com/rs/zerolog/log"
)

type Messages struct {
	DB       *postgres.DB
	Messages chan *entity.Message
	Batches  chan *pgx.Batch
	Config   configs.Postgres
}

const insert = "insert into messenger.messages (id, sender_id, recipient_id, text) values ($1, $2, $3, $4) on conflict (id) do update set text = $4"

func New(ctx context.Context, db *postgres.DB, conf configs.Postgres) *Messages {
	messageChan := make(chan *entity.Message, conf.BatchSize)
	batchChan := make(chan *pgx.Batch, 4)

	messages := &Messages{
		DB:       db,
		Messages: messageChan,
		Batches:  batchChan,
		Config:   conf,
	}

	go messages.aggregate(ctx)

	for i := 0; i < 4; i++ {
		go messages.startSender(ctx)
	}

	return messages
}

func (m *Messages) aggregate(ctx context.Context) {
	batch := &pgx.Batch{}

	for {
		select {
		case msg := <-m.Messages:
			batch.Queue(insert, msg.MessageID, msg.SenderID, msg.RecipientID, msg.Text)

			if batch.Len() == m.Config.BatchSize {
				m.Batches <- batch
				batch = &pgx.Batch{}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (m *Messages) startSender(ctx context.Context) {
	for {
		select {
		case batch := <-m.Batches:
			err := m.sendBatch(batch)
			if err != nil {
				log.Error().Err(err).Msg("")
			}
		case <-ctx.Done():
			return
		}
	}
}

func (m *Messages) sendBatch(batch *pgx.Batch) error {
	const op = "Messages.sendBatch"

	results := m.DB.SendBatch(context.Background(), batch)

	for i := 0; i < batch.Len(); i++ {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return results.Close()
}
