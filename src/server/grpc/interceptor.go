package grpc

import (
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type MessengerServerStream struct {
	grpc.ServerStream
}

func (s *MessengerServerStream) RecvMsg(m interface{}) error {
	if err := s.ServerStream.RecvMsg(m); err != nil {
		log.Error().Err(err).Msg("")
		return err
	}
	return nil
}

func NewInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapper := &MessengerServerStream{
			ServerStream: ss,
		}
		return handler(srv, wrapper)
	}
}
