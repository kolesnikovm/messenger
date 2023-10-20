package messenger

import (
	"fmt"
	"strconv"

	"github.com/kolesnikovm/messenger/proto"
	"google.golang.org/grpc/metadata"
)

func (h *Handler) GetMessage(msgRequest *proto.MessaggeRequest, stream proto.Messenger_GetMessageServer) error {
	const op = "Handler.GetMessage"

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return fmt.Errorf("%s: no metadata in request", op)
	}

	userID, err := strconv.Atoi(md.Get("x-user-id")[0])
	if err != nil {
		return fmt.Errorf("%s: failed to get user id from metadata: %v", op, md.Get("x-user-id"))
	}

	deviceID, err := strconv.Atoi(md.Get("x-device-id")[0])
	if err != nil {
		return fmt.Errorf("%s: failed to get device id from metadata: %v", op, md.Get("x-device-id"))
	}

	messageCh := h.Usecase.Get(stream.Context(), uint64(userID), deviceID)

	for message := range messageCh {
		protoMsg := &proto.Message{
			MessageID:   message.MessageID.Bytes(),
			SenderID:    message.SenderID,
			RecipientID: message.RecipientID,
			Text:        message.Text,
		}

		if err := stream.Send(protoMsg); err != nil {
			return err
		}
	}

	return nil
}
