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

type StringContextKey string

type Handler struct {
	Usecase usecase.Message
}

func NewHandler(usecase usecase.Message) *Handler {
	return &Handler{
		Usecase: usecase,
	}
}

func convertPbToEntity(msg *pb.Message) (*entity.Message, error) {
	const op = "convertPbToEntity"

	messageID, err := ulid.Parse(msg.MessageID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := entity.NewMessage(
		messageID,
		msg.SenderID,
		msg.RecipientID,
		msg.Text,
	)

	return res, nil
}

func convertEntityToPb(msg *entity.Message) *pb.Message {
	return &pb.Message{
		MessageID:   msg.MessageID.String(),
		SenderID:    msg.SenderID,
		RecipientID: msg.RecipientID,
		Text:        msg.Text,
	}
}

func GetHeader(md metadata.MD, header string) (uint64, error) {
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
