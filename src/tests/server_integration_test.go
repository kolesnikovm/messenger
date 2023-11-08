package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/entity"
	"github.com/kolesnikovm/messenger/proto"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestSendMessage(t *testing.T) {
	config, err := configs.NewServerConfig("")
	require.NoError(t, err)

	suite, cleanup, err := InitializeSuite(t, config)
	require.NoError(t, err)
	defer cleanup()

	messageID := ulid.Make()
	entityMessage := entity.Message{
		MessageID:   messageID,
		SenderID:    1,
		RecipientID: 2,
		Text:        "test",
	}
	suite.messageSender.EXPECT().Send(mock.AnythingOfType("*context.valueCtx"), entityMessage).Return(nil)

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-user-id", "1")
	stream, err := suite.messengerServiceClient.SendMessage(ctx)
	require.NoErrorf(t, err, "Failed to create stream")

	message := &proto.Message{
		MessageID:   messageID.Bytes(),
		SenderID:    1,
		RecipientID: 2,
		Text:        "test",
	}
	err = stream.Send(message)
	require.NoErrorf(t, err, "Error in %v.Send(%v)", stream, message)

	reply, err := stream.CloseAndRecv()
	require.NoErrorf(t, err, "Error in %v.CloseAndRecv()", stream)

	require.Equal(t, 0, int(reply.GetErrorCount()))
}

func TestSendMessageInternalError(t *testing.T) {
	config, err := configs.NewServerConfig("")
	require.NoError(t, err)

	suite, cleanup, err := InitializeSuite(t, config)
	require.NoError(t, err)
	defer cleanup()

	messageID := ulid.Make()
	notifierError := errors.New("notifier error")
	suite.messageSender.EXPECT().Send(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("entity.Message")).Return(notifierError)

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-user-id", "1")
	stream, err := suite.messengerServiceClient.SendMessage(ctx)
	require.NoErrorf(t, err, "Failed to create stream")

	message := &proto.Message{
		MessageID:   messageID.Bytes(),
		SenderID:    1,
		RecipientID: 2,
		Text:        "test",
	}
	err = stream.Send(message)
	require.NoErrorf(t, err, "Error in %v.Send(%v)", stream, message)

	_, err = stream.CloseAndRecv()
	require.EqualError(t, err, "rpc error: code = Internal desc = Internal")
}

func TestSendMessageAuthError(t *testing.T) {
	config, err := configs.NewServerConfig("")
	require.NoError(t, err)

	suite, cleanup, err := InitializeSuite(t, config)
	require.NoError(t, err)
	defer cleanup()

	ctx := context.Background()
	stream, err := suite.messengerServiceClient.SendMessage(ctx)
	require.NoErrorf(t, err, "Failed to create stream")

	message := &proto.Message{
		MessageID:   ulid.Make().Bytes(),
		SenderID:    1,
		RecipientID: 2,
		Text:        "test",
	}
	err = stream.Send(message)
	require.NoErrorf(t, err, "Error in %v.Send(%v)", stream, message)

	_, err = stream.CloseAndRecv()
	require.EqualError(t, err, "rpc error: code = Unauthenticated desc = Unauthenticated")
}

func TestGetMessage(t *testing.T) {
	config, err := configs.NewServerConfig("")
	require.NoError(t, err)

	suite, cleanup, err := InitializeSuite(t, config)
	require.NoError(t, err)
	defer cleanup()

	messageID := ulid.Make()
	entityMessage := &entity.Message{
		MessageID:   messageID,
		SenderID:    1,
		RecipientID: 2,
		Text:        "test",
	}

	messageCh := make(chan *entity.Message, 1)
	messageCh <- entityMessage
	suite.messageSender.EXPECT().Get(mock.AnythingOfType("*context.valueCtx"), uint64(1), mock.AnythingOfType("ulid.ULID")).Return(messageCh, func() {})

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-user-id", "1")
	stream, err := suite.messengerServiceClient.GetMessage(ctx, &proto.MessaggeRequest{})
	require.NoErrorf(t, err, "Failed to create stream")

	message, err := stream.Recv()
	require.NoErrorf(t, err, "Error in %v.Send(%v)", stream, message)

	require.Equal(t, message.MessageID, messageID.Bytes())
}

func TestGetMessageHistory(t *testing.T) {
	config, err := configs.NewServerConfig("")
	require.NoError(t, err)

	suite, cleanup, err := InitializeSuite(t, config)
	require.NoError(t, err)
	defer cleanup()

	entityMessages := []*entity.Message{{
		MessageID:   ulid.Make(),
		SenderID:    1,
		RecipientID: 2,
		Text:        "test",
	}}
	suite.messageStore.EXPECT().GetMessageHistory(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("ulid.ULID"), "1:2").Return(entityMessages, nil)

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-user-id", "1")
	historyResponse, err := suite.messengerServiceClient.GetMessageHistory(ctx, &proto.HistoryRequest{ChatID: "1:2", MessageID: ulid.Make().Bytes()})
	require.NoErrorf(t, err, "Failed to get mesage history")

	message := historyResponse.Messages[0]
	require.Equal(t, entityMessages[0].MessageID.Bytes(), message.MessageID)
}

func TestGetMessageHistoryInternalError(t *testing.T) {
	config, err := configs.NewServerConfig("")
	require.NoError(t, err)

	suite, cleanup, err := InitializeSuite(t, config)
	require.NoError(t, err)
	defer cleanup()

	storeError := errors.New("store error")
	suite.messageStore.EXPECT().GetMessageHistory(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("ulid.ULID"), "1:2").Return(nil, storeError)

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-user-id", "1")
	_, err = suite.messengerServiceClient.GetMessageHistory(ctx, &proto.HistoryRequest{ChatID: "1:2", MessageID: ulid.Make().Bytes()})
	require.EqualError(t, err, "rpc error: code = Internal desc = Internal")
}

func TestGetMessageHistoryAuthError(t *testing.T) {
	config, err := configs.NewServerConfig("")
	require.NoError(t, err)

	suite, cleanup, err := InitializeSuite(t, config)
	require.NoError(t, err)
	defer cleanup()

	ctx := context.Background()
	_, err = suite.messengerServiceClient.GetMessageHistory(ctx, &proto.HistoryRequest{ChatID: "1:2", MessageID: ulid.Make().Bytes()})
	require.EqualError(t, err, "rpc error: code = Unauthenticated desc = Unauthenticated")

	st, ok := status.FromError(err)
	require.Equal(t, true, ok)

	details := st.Details()[0].(*errdetails.BadRequest)
	require.Equal(t, "failed to get user id from metadata", details.FieldViolations[0].Description)
}

func TestGetMessageHistoryArgumentError(t *testing.T) {
	config, err := configs.NewServerConfig("")
	require.NoError(t, err)

	suite, cleanup, err := InitializeSuite(t, config)
	require.NoError(t, err)
	defer cleanup()

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-user-id", "1")
	_, err = suite.messengerServiceClient.GetMessageHistory(ctx, &proto.HistoryRequest{ChatID: "", MessageID: ulid.Make().Bytes()})
	require.EqualError(t, err, "rpc error: code = InvalidArgument desc = InvalidArgument")

	st, ok := status.FromError(err)
	require.Equal(t, true, ok)

	details := st.Details()[0].(*errdetails.BadRequest)
	require.Equal(t, "failed to parse chat id from: ", details.FieldViolations[0].Description)

	ctx = metadata.AppendToOutgoingContext(ctx, "x-user-id", "")
	_, err = suite.messengerServiceClient.GetMessageHistory(ctx, &proto.HistoryRequest{ChatID: "1:2", MessageID: []byte{}})
	require.EqualError(t, err, "rpc error: code = InvalidArgument desc = InvalidArgument")

	st, ok = status.FromError(err)
	require.Equal(t, true, ok)

	details = st.Details()[0].(*errdetails.BadRequest)
	require.Equal(t, "failed to get message id if from: ", details.FieldViolations[0].Description)
}

func TestGetMessageHistoryPermissionError(t *testing.T) {
	config, err := configs.NewServerConfig("")
	require.NoError(t, err)

	suite, cleanup, err := InitializeSuite(t, config)
	require.NoError(t, err)
	defer cleanup()

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-user-id", "1")
	_, err = suite.messengerServiceClient.GetMessageHistory(ctx, &proto.HistoryRequest{ChatID: "2:3", MessageID: ulid.Make().Bytes()})
	require.EqualError(t, err, "rpc error: code = NotFound desc = Chat not found")
}
