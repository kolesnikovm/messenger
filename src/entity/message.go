package entity

import "github.com/oklog/ulid/v2"

type Message struct {
	MessageID   ulid.ULID
	SenderID    uint64
	RecipientID uint64
	Text        string
}
