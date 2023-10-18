package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/entity"
	"github.com/kolesnikovm/messenger/proto"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSendMessage(t *testing.T) {
	config, err := configs.NewServerConfig("")
	require.NoError(t, err)

	suite, err := InitializeSuite(t, config)
	require.NoError(t, err)
	defer suite.Stop()

	messageID := ulid.Make()
	entityMessage := entity.Message{
		MessageID:   messageID,
		SenderID:    1,
		RecipientID: 2,
		Text:        "test",
	}
	suite.messageSender.EXPECT().Send(mock.AnythingOfType("*context.valueCtx"), entityMessage).Return(nil)

	ctx := context.Background()
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

func TestSendMessageError(t *testing.T) {
	config, err := configs.NewServerConfig("")
	require.NoError(t, err)

	suite, err := InitializeSuite(t, config)
	require.NoError(t, err)
	defer suite.Stop()

	messageID := ulid.Make()
	entityMessage := entity.Message{
		MessageID:   messageID,
		SenderID:    1,
		RecipientID: 2,
		Text:        "test",
	}
	notifierError := errors.New("notifier error")
	suite.messageSender.EXPECT().Send(mock.AnythingOfType("*context.valueCtx"), entityMessage).Return(notifierError)

	ctx := context.Background()
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
	require.EqualError(t, err, fmt.Sprintf("rpc error: code = Unknown desc = %s", notifierError.Error()))
}
