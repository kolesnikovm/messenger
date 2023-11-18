package errors

import (
	"context"

	"github.com/kolesnikovm/messenger/server/grpc/messenger"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		userID, err := getUser(ctx)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, messenger.StringContextKey("userID"), userID)

		resp, err := handler(ctx, req)
		if err != nil {
			return nil, getError(err)
		}

		return resp, nil
	}
}

type MessengerServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (m *MessengerServerStream) Context() context.Context {
	return m.ctx
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		userID, err := getUser(ss.Context())
		if err != nil {
			return err
		}

		ctx := context.WithValue(ss.Context(), messenger.StringContextKey("userID"), userID)

		wrapper := &MessengerServerStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		err = handler(srv, wrapper)
		if err != nil {
			return getError(err)
		}
		return nil
	}
}

func getUser(ctx context.Context) (uint64, error) {
	const op = "getUser"

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, status.Error(codes.Unauthenticated, "No metadata in request")
	}

	userID, err := messenger.GetHeader(md, "x-user-id")
	if err != nil {
		log.Error().Err(err).Send()

		st := status.New(codes.Unauthenticated, codes.Unauthenticated.String())
		fv := &errdetails.BadRequest_FieldViolation{
			Field:       "x-user-id",
			Description: "failed to get user id from metadata",
		}
		br := &errdetails.BadRequest{}
		br.FieldViolations = append(br.FieldViolations, fv)

		st, err := st.WithDetails(br)
		if err != nil {
			log.Error().Err(err).Msgf("%s: error calling status.WithDetails", op)

			return 0, status.Error(codes.Internal, codes.Internal.String())
		}

		return 0, st.Err()
	}

	return userID, nil
}

func getError(err error) error {
	st, ok := status.FromError(err)

	if !ok || st.Code() == codes.Unknown {
		log.Error().Err(err).Send()

		return status.Error(codes.Internal, codes.Internal.String())
	}

	return st.Err()
}
