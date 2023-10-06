package messenger

import (
	"github.com/kolesnikovm/messenger/entity"
	pb "github.com/kolesnikovm/messenger/proto"
	usecase "github.com/kolesnikovm/messenger/usecase"
)

type Handler struct {
	Usecase usecase.Message
}

func NewHandler(usecase usecase.Message) *Handler {
	return &Handler{
		Usecase: usecase,
	}
}

func (s *Handler) transformMessageRPC(msg *pb.Message) entity.Message {
	res := entity.Message{
		Text: msg.Text,
	}

	return res
}
