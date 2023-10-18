package grpc

import (
	"reflect"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MessengerServerStream struct {
	grpc.ServerStream
}

func (w *MessengerServerStream) RecvMsg(m any) error {
	log.Debug().Str("message type", reflect.TypeOf(m).String()).Msg("message received")
	return w.ServerStream.RecvMsg(m)
}

func (w *MessengerServerStream) SendMsg(m any) error {
	log.Debug().Str("message type", reflect.TypeOf(m).String()).Msg("message sent")
	return w.ServerStream.SendMsg(m)
}

func NewInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapper := &MessengerServerStream{
			ServerStream: ss,
		}

		err := handler(srv, wrapper)
		if err != nil {
			st, _ := status.FromError(err)
			log.Error().Err(st.Err()).Msg("")

			if st.Code() == codes.Unknown {
				return status.Error(codes.Internal, "Internal server error")
			}
		}
		return err
	}
}
