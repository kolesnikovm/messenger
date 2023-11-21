package grpc

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/kolesnikovm/messenger/proto"
	"github.com/kolesnikovm/messenger/server/grpc/interceptors/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type ServerBuilder struct {
	MessengerServer proto.MessengerServer
}

func interceptorLogger(l zerolog.Logger) logging.Logger {
	const op = "interceptorLogger"

	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l := l.With().Fields(fields).Logger()

		switch lvl {
		case logging.LevelDebug:
			l.Debug().Msg(msg)
		case logging.LevelInfo:
			l.Info().Msg(msg)
		case logging.LevelWarn:
			l.Warn().Msg(msg)
		case logging.LevelError:
			l.Error().Msg(msg)
		default:
			log.Error().Msgf("%s: unknown level %v", op, lvl)
		}
	})
}

func (s *ServerBuilder) Build() *grpc.Server {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptorLogger(log.Logger)),
			errors.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(interceptorLogger(log.Logger)),
			errors.StreamServerInterceptor(),
		),
	)
	proto.RegisterMessengerServer(srv, s.MessengerServer)

	return srv
}
