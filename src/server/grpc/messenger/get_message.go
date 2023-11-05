package messenger

import (
	"fmt"

	"github.com/kolesnikovm/messenger/proto"
	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/metadata"
)

func (h *Handler) GetMessage(msgRequest *proto.MessaggeRequest, stream proto.Messenger_GetMessageServer) error {
	const op = "Handler.GetMessage"

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return fmt.Errorf("%s: no metadata in request", op)
	}

	userID, err := getHeader(md, "x-user-id")
	if err != nil {
		return err
	}

	sessionID := ulid.Make()

	messageCh, cleanup := h.Usecase.Get(stream.Context(), userID, sessionID)
	defer cleanup()

	for {
		select {
		case message := <-messageCh:
			protoMsg := &proto.Message{
				MessageID:   message.MessageID.Bytes(),
				SenderID:    message.SenderID,
				RecipientID: message.RecipientID,
				Text:        message.Text,
			}

			if err := stream.Send(protoMsg); err != nil {
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}
