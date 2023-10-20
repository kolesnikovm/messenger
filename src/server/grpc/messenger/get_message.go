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

	userID, err := getHeader(md, "x-user-id")
	if err != nil {
		return err
	}

	deviceID, err := getHeader(md, "x-device-id")
	if err != nil {
		return err
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

func getHeader(md metadata.MD, header string) (int, error) {
	if len(md.Get(header)) > 0 {
		id, err := strconv.Atoi(md.Get(header)[0])
		if err != nil {
			return 0, fmt.Errorf("failed to parse header %s: %v", header, md.Get("x-user-id"))
		}
		return id, nil
	} else {
		return 0, fmt.Errorf("no %s header in request", header)
	}
}
