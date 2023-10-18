package messenger

import (
	"fmt"

	"github.com/kolesnikovm/messenger/entity"
	pb "github.com/kolesnikovm/messenger/proto"
	usecase "github.com/kolesnikovm/messenger/usecase"
	"github.com/oklog/ulid/v2"
)

type Handler struct {
	Usecase usecase.Message
}

func NewHandler(usecase usecase.Message) *Handler {
	return &Handler{
		Usecase: usecase,
	}
}

func (s *Handler) transformMessageRPC(msg *pb.Message) (entity.Message, error) {
	const op = "Handler.transformMessageRPC"

	messageID := ulid.ULID{}
	if err := messageID.UnmarshalBinary(msg.MessageID); err != nil {
		return entity.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	res := entity.Message{
		MessageID:   messageID,
		SenderID:    msg.SenderID,
		RecipientID: msg.RecipientID,
		Text:        msg.Text,
	}

	return res, nil
}
