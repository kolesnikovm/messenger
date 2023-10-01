package grpc

import (
	"io"

	"github.com/kolesnikovm/messenger/internal/entity"
	usecase "github.com/kolesnikovm/messenger/internal/usecase/message"
	pb "github.com/kolesnikovm/messenger/proto/message"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type server struct {
	usecase usecase.Message

	pb.UnimplementedMessengerServer
}

func NewMessageServerGrpc(gs *grpc.Server, messageUsacase usecase.Message) {

	messageServer := &server{
		usecase: messageUsacase,
	}

	pb.RegisterMessengerServer(gs, messageServer)
}

func (s *server) transformMessageRPC(msg *pb.Message) entity.Message {
	res := entity.Message{
		Text: msg.Text,
	}

	return res
}

func (s *server) Send(stream pb.Messenger_SendMessageServer) error {
	errorCount := int32(0)
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Status{
				ErrorCount: errorCount,
			})
		}
		if err != nil {
			log.Error().Err(err).Msg("")
			return err
		}

		m := s.transformMessageRPC(message)

		err = s.usecase.Send(m)
		if err != nil {
			log.Error().Err(err).Msg("failed to send message")
			errorCount++
		}
	}
}
