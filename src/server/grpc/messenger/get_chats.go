package messenger

import (
	"context"

	"github.com/kolesnikovm/messenger/proto"
)

func (s *Handler) GetChats(ctx context.Context, req *proto.ChatsRequest) (*proto.ChatsResponse, error) {
	userID := ctx.Value(StringContextKey("userID")).(uint64)

	chats, err := s.Usecase.GetChats(ctx, userID)
	if err != nil {
		return nil, err
	}

	protoChats := []*proto.Chat{}

	for _, chat := range chats {
		protoChats = append(protoChats, &proto.Chat{
			ChatID:         chat.ID,
			UnreadMessages: chat.UnreadMessages,
		})
	}

	response := &proto.ChatsResponse{}
	response.Chats = protoChats

	return response, nil
}
