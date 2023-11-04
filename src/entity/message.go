package entity

import (
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
		stringSlice = append(stringSlice, strconv.Itoa(int(id)))
	}

	return strings.Join(stringSlice, ":")
}
