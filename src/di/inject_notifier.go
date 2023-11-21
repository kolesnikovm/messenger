package di

import (
	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/notifier"
	"github.com/kolesnikovm/messenger/notifier/kafka"
)

func ProvideNotifier(conf configs.ServerConfig) (notifier.MessageSender, func(), error) {
	messageSender, err := kafka.New(conf.Kafka)

	cleanup := func() {
		messageSender.Close()
	}

	return messageSender, cleanup, err
}
