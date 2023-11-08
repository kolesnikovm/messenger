package interceptors

import (
	"context"
	"reflect"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MessengerServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (m *MessengerServerStream) Context() context.Context {
	return m.ctx
}

func (w *MessengerServerStream) RecvMsg(m any) error {
	log.Debug().Str("message type", reflect.TypeOf(m).String()).Msg("message received")
	return w.ServerStream.RecvMsg(m)
}

func (w *MessengerServerStream) SendMsg(m any) error {
	log.Debug().Str("message type", reflect.TypeOf(m).String()).Msg("message sent")
	return w.ServerStream.SendMsg(m)
}

func NewStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		userID, err := getUser(ss.Context())
		if err != nil {
			return err
		}

		ctx := context.WithValue(ss.Context(), "userID", userID)

		wrapper := &MessengerServerStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		err = handler(srv, wrapper)
		if err != nil {
			st, ok := status.FromError(err)

			if !ok || st.Code() == codes.Unknown {
				log.Error().Err(err).Send()

				return status.Error(codes.Internal, codes.Internal.String())
			}

			return st.Err()
		}
		return nil
	}
}
