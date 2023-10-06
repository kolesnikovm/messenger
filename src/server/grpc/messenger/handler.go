package messenger

import (
	"github.com/kolesnikovm/messenger/entity"
	pb "github.com/kolesnikovm/messenger/proto"
	usecase "github.com/kolesnikovm/messenger/usecase"
)

type Handler struct {
	Usecase usecase.Message
}

func (s *Handler) transformMessageRPC(msg *pb.Message) entity.Message {
	res := entity.Message{
		Text: msg.Text,
	}

	return res
}
