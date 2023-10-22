package message

import (
	"context"
	"fmt"
	"strings"

	"github.com/kolesnikovm/messenger/entity"
	"github.com/oklog/ulid/v2"
)

func (m *MessageUseCase) Get(ctx context.Context, userID string, sessionID ulid.ULID, chatID string) (<-chan *entity.Message, func(), error) {
	var recepientID string

	if chatID != "" {
		chatParticipants := strings.Split(chatID, ":")
		if len(chatParticipants) != 2 {
			return nil, nil, fmt.Errorf("failed to get chat participants from chat id: %s", chatID)
		}
		if chatParticipants[0] == userID || chatParticipants[1] == userID {
			recepientID = chatID
		}
	} else {
		recepientID = userID
	}

	messageCh, cleanup := m.messageSender.Get(ctx, recepientID, sessionID)

	return messageCh, cleanup, nil
}
