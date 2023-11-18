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

const maxMessageCount = 1_000

func (h *Handler) GetMessageHistory(ctx context.Context, req *proto.HistoryRequest) (*proto.HistoryResponse, error) {
	const op = "Handler.GetMessageHistory"

	userID := ctx.Value(StringContextKey("userID")).(uint64)

	messageID, err := ulid.Parse(req.MessageID)
	if err != nil {
		log.Error().Err(err).Str("op", op).Send()

		statusError := composeInvalidArgumentError("HistoryRequest.messageID", fmt.Sprintf("failed to get message id if from: %s", req.MessageID))

		return nil, statusError
	}

	if req.MessageCount == 0 || req.MessageCount > maxMessageCount {
		statusError := composeInvalidArgumentError("HistoryRequest.messageCount", fmt.Sprintf("messageCount must be in (0, %d]", maxMessageCount))

		return nil, statusError
	}

	if err := checkPermission(userID, req.ChatID); err != nil {
		return nil, err
	}

	messages, err := h.Usecase.GetHistory(ctx, req.ChatID, messageID, userID, req.MessageCount, req.GetDirection().String())
	if err != nil {
		return nil, err
	}

	protoMessages := make([]*proto.Message, 0, len(messages))
	for _, message := range messages {
		protoMessages = append(protoMessages, convertEntityToPb(message))
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

func checkPermission(userID uint64, chatID string) error {
	const op = "Handler.checkPermission"

	switch entity.GetChatType(chatID) {
	case entity.Group:
		groupID, err := entity.GetGroupID(chatID)
		if err != nil {
			log.Error().Err(err).Send()

			statusError := composeInvalidArgumentError("HistoryRequest.chatID", fmt.Sprintf("failed to parse group id from: %s", chatID))

			return statusError
		}

		if !isGroupMember(userID, groupID) {
			log.Error().Msgf("%s: permission denied for user %d on group %d", op, userID, groupID)

			return status.Error(codes.NotFound, "Chat not found")
		}
	case entity.Channel:
		channelID, err := entity.GetChannelID(chatID)
		if err != nil {
			log.Error().Err(err).Send()

			statusError := composeInvalidArgumentError("HistoryRequest.chatID", fmt.Sprintf("failed to parse channel id from: %s", chatID))

			return statusError
		}

		if !isChannelMember(userID, channelID) {
			log.Error().Msgf("%s: permission denied for user %d on channel %d", op, userID, channelID)

			return status.Error(codes.NotFound, "Chat not found")
		}
	case entity.P2P:
		user1, user2, err := entity.GetUserIDs(chatID)
		if err != nil {
			log.Error().Err(err).Send()

			statusError := composeInvalidArgumentError("HistoryRequest.chatID", fmt.Sprintf("failed to parse user ids from: %s", chatID))

			return statusError
		}

		if user1 != userID && user2 != userID {
			log.Error().Msgf("%s: permission denied for user %d on chat %s", op, userID, chatID)

			return status.Error(codes.NotFound, "Chat not found")
		}
	default:
		log.Error().Msgf("%s: failed to get chat type from: %s", op, chatID)

		statusError := composeInvalidArgumentError("HistoryRequest.chatID", fmt.Sprintf("failed to get chat type from: %s", chatID))

		return statusError
	}

	return nil
}

// TODO implemen for groups
func isGroupMember(userID, groupID uint64) bool {
	return false
}

// TODO implemen for groups
func isChannelMember(userID, groupID uint64) bool {
	return false
}
