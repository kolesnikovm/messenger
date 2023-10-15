package tests

import (
	"context"
	"testing"

	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/proto"
	"github.com/stretchr/testify/require"
)

func TestSendMessage(t *testing.T) {
	config, err := configs.NewServerConfig("")
	require.NoError(t, err)

	suite, err := InitializeSuite(t, config)
	require.NoError(t, err)
	defer suite.Stop()

	ctx := context.Background()
	stream, err := suite.messengerServiceClient.SendMessage(ctx)
	require.NoErrorf(t, err, "Failed to create stream")

	message := &proto.Message{Text: "test", RecipientID: 1}
	err = stream.Send(message)
	require.NoErrorf(t, err, "Error in %v.Send(%v)", stream, message)

	reply, err := stream.CloseAndRecv()
	require.NoErrorf(t, err, "Error in %v.CloseAndRecv()", stream)

	require.Equal(t, 0, int(reply.GetErrorCount()))
}
