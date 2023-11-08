package interceptors

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

func NewUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		userID, err := getUser(ctx)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, "userID", userID)

		resp, err := handler(ctx, req)
		if err != nil {
			st, ok := status.FromError(err)

			if !ok || st.Code() == codes.Unknown {
				log.Error().Err(err).Send()

				return nil, status.Error(codes.Internal, codes.Internal.String())
			}

			return nil, st.Err()
		}

		return resp, nil
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
			log.Error().Err(err).Str("op", op).Send()

			return 0, status.Error(codes.Internal, codes.Internal.String())
		}

		return 0, st.Err()
	}

	return userID, nil
}
