package message

type MessageUseCase struct{}

func New() *MessageUseCase {
	return &MessageUseCase{}
}
