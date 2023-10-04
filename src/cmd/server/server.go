package server

import (
	"fmt"
	"net"

	"github.com/kolesnikovm/messenger/configs"
	"github.com/kolesnikovm/messenger/server/grpc"
	"github.com/kolesnikovm/messenger/server/grpc/messenger"
	messageUseCase "github.com/kolesnikovm/messenger/usecase/message"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "server",
	Short: "Start messenger in server mode",
	Run: func(cmd *cobra.Command, args []string) {
		configFile, err := cmd.InheritedFlags().GetString("config")
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}

		// check config is ok
		config, err := configs.NewServerConfig(configFile)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to instantiate config")
		}

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ListenPort))
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to listen port %d", config.ListenPort)
		}

		messageUseCase := messageUseCase.New()

		serverBuilder := grpc.ServerBuilder{
			MessengerServer: &messenger.Handler{
				Usecase: messageUseCase,
			},
		}
		server := serverBuilder.Build()

		log.Info().Msgf("Messenger server listening on %v", lis.Addr())
		if err := server.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("failed to start grpc server")
		}
	},
}
