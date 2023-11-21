package entity

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/oklog/ulid/v2"
)

type ChatType string

const (
	Group   ChatType = "Group"
	Channel ChatType = "Channel"
	P2P     ChatType = "P2P"
)

type Message struct {
	MessageID   ulid.ULID
	SenderID    uint64
	RecipientID uint64
	chatType    ChatType
	Text        string
}

func NewMessage(messageID ulid.ULID, senderID, recipientID uint64, text string) *Message {
	var chatType ChatType

	switch {
	case isGroup(recipientID):
		chatType = Group
	case isChannel(recipientID):
		chatType = Channel
	default:
		chatType = P2P
	}

	return &Message{
		MessageID:   messageID,
		SenderID:    senderID,
		RecipientID: recipientID,
		chatType:    chatType,
		Text:        text,
	}
}

func (m *Message) GetChatID() string {
	var chatID string

	switch m.chatType {
	case Group:
		chatID = m.getGroupID()
	case Channel:
		chatID = m.getChannelID()
	default:
		chatID = m.getP2PID()
	}

	return chatID
}

// TODO implement for groups
func isGroup(id uint64) bool {
	return false
}

// TODO implement for groups
func isChannel(id uint64) bool {
	return false
}

func (m *Message) getP2PID() string {
	slice := []uint64{m.SenderID, m.RecipientID}
	sort.SliceStable(slice, func(i, j int) bool { return slice[i] < slice[j] })

	var stringSlice []string
	for _, id := range slice {
		stringSlice = append(stringSlice, strconv.FormatUint(id, 10))
	}

	return strings.Join(stringSlice, ":")
}

func (m *Message) getGroupID() string {
	return fmt.Sprintf("g:%d", m.RecipientID)
}

func (m *Message) getChannelID() string {
	return fmt.Sprintf("c:%d", m.RecipientID)
}

func GetChatType(chatID string) ChatType {
	switch {
	case strings.HasPrefix(chatID, "g:"):
		return Group
	case strings.HasPrefix(chatID, "c:"):
		return Channel
	default:
		return P2P
	}
}

func GetGroupID(chatID string) (uint64, error) {
	const op = "GetGroupID"

	tokens := strings.Split(chatID, ":")
	if len(tokens) != 2 {
		return 0, fmt.Errorf("%s: failed to get group id from chat id: %s", op, chatID)
	}

	groupID, err := strconv.ParseUint(tokens[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return groupID, nil
}

func GetChannelID(chatID string) (uint64, error) {
	const op = "GetChannelID"

	tokens := strings.Split(chatID, ":")
	if len(tokens) != 2 {
		return 0, fmt.Errorf("%s: failed to get channel id from chat id: %s", op, chatID)
	}

	channelID, err := strconv.ParseUint(tokens[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return channelID, nil
}

func GetUserIDs(chatID string) (uint64, uint64, error) {
	const op = "GetUserIDs"

	tokens := strings.Split(chatID, ":")
	if len(tokens) != 2 {
		return 0, 0, fmt.Errorf("%s: failed to get user ids from chat id: %s", op, chatID)
	}

	userID1, err := strconv.ParseUint(tokens[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("%s: %w", op, err)
	}

	userID2, err := strconv.ParseUint(tokens[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID1, userID2, nil
}
