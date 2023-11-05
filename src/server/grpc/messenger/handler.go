package messenger

import (
	"fmt"
	"strconv"

	"github.com/kolesnikovm/messenger/entity"
	pb "github.com/kolesnikovm/messenger/proto"
	usecase "github.com/kolesnikovm/messenger/usecase"
	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/metadata"
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

func getHeader(md metadata.MD, header string) (uint64, error) {
	if len(md.Get(header)) > 0 {
		id, err := strconv.ParseUint(md.Get(header)[0], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse header %s: %v", header, md.Get(header))
		}
		return id, nil
	} else {
		return 0, fmt.Errorf("no %s header in request", header)
	}
}
