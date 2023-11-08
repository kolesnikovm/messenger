package entity

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/oklog/ulid/v2"
)

type Message struct {
	MessageID   ulid.ULID
	SenderID    uint64
	RecipientID uint64
	Text        string
}

func (m *Message) GetChatID() string {
	slice := []uint64{m.SenderID, m.RecipientID}
	sort.SliceStable(slice, func(i, j int) bool { return slice[i] < slice[j] })

	var stringSlice []string
	for _, id := range slice {
		stringSlice = append(stringSlice, strconv.FormatUint(id, 10))
	}

	return strings.Join(stringSlice, ":")
}

func ParseChatID(chatID string) (user1 uint64, user2 uint64, err error) {
	const op = "entity.ParseChatID"

	chatParticipants := strings.Split(chatID, ":")

	if len(chatParticipants) != 2 {
		return 0, 0, fmt.Errorf("%s: failed to get chat participants from chat id: %s", op, chatID)
	}

	user1, err = strconv.ParseUint(chatParticipants[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("%s: failed to parse user id from: %s", op, chatParticipants[0])
	}

	user2, err = strconv.ParseUint(chatParticipants[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("%s: failed to parse user id from: %s", op, chatParticipants[1])
	}

	return user1, user2, nil
}
