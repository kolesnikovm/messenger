package messenger

import (
	"context"
	"fmt"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/kolesnikovm/messenger/proto"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetMessageHistory(ctx context.Context, req *proto.HistoryRequest) (*proto.HistoryResponse, error) {
	const op = "Handler.GetMessageHistory"

	userID := ctx.Value("userID").(uint64)

	messageID, err := ulid.Parse(req.MessageID)
	if err != nil {
		log.Error().Err(err).Str("op", op).Send()

		statusError := composeInvalidArgumentError("HistoryRequest.messageID", fmt.Sprintf("failed to get message id if from: %s", req.MessageID))

		return nil, statusError
	}

	user1, user2, err := entity.ParseChatID(req.ChatID)
	if err != nil {
		log.Error().Err(err).Send()

		statusError := composeInvalidArgumentError("HistoryRequest.chatID", fmt.Sprintf("failed to parse chat id from: %s", req.ChatID))

		return nil, statusError
	}

	if user1 != userID && user2 != userID {
		log.Error().Msgf("%s: permission denied for user %d on chat %s", op, userID, req.ChatID)

		return nil, status.Error(codes.NotFound, "Chat not found")
	}

	messages, err := h.Usecase.GetHistory(ctx, req.ChatID, messageID, userID)
	if err != nil {
		return nil, err
	}

	protoMessages := make([]*proto.Message, 0, len(messages))
	for _, message := range messages {
		protoMessages = append(protoMessages, &proto.Message{
			MessageID:   message.MessageID.String(),
			SenderID:    message.SenderID,
			RecipientID: message.RecipientID,
			Text:        message.Text,
		})
	}

	return &proto.HistoryResponse{
		Messages: protoMessages,
	}, nil
}

func composeInvalidArgumentError(agrument string, details string) error {
	const op = "composeInvalidArgumentError"

	st := status.New(codes.InvalidArgument, codes.InvalidArgument.String())
	fv := &errdetails.BadRequest_FieldViolation{
		Field:       agrument,
		Description: details,
	}
	br := &errdetails.BadRequest{}
	br.FieldViolations = append(br.FieldViolations, fv)

	st, err := st.WithDetails(br)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return st.Err()
}
