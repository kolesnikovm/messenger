package messenger

import (
	"context"
	"fmt"

	"github.com/kolesnikovm/messenger/proto"
	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/metadata"
)

func (h *Handler) GetMessageHistory(ctx context.Context, req *proto.HistoryRequest) (*proto.HistoryResponse, error) {
	const op = "Handler.GetMessageHistory"

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("%s: no metadata in request", op)
	}

	userID, err := getHeader(md, "x-user-id")
	if err != nil {
		return nil, err
	}

	messageID := ulid.ULID{}
	if err := messageID.UnmarshalBinary(req.MessageID); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if req.ChatID == "" {
		return nil, fmt.Errorf("%s: no chat id in request", op)
	}

	messages, err := h.Usecase.GetHistory(ctx, req.ChatID, messageID, userID)
	if err != nil {
		return nil, err
	}

	protoMessages := make([]*proto.Message, 0, len(messages))
	for _, message := range messages {
		protoMessages = append(protoMessages, &proto.Message{
			MessageID:   message.MessageID.Bytes(),
			SenderID:    message.SenderID,
			RecipientID: message.RecipientID,
			Text:        message.Text,
		})
	}

	return &proto.HistoryResponse{
		Messages: protoMessages,
	}, nil
}
