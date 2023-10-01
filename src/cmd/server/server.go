package server

import (
	"fmt"
	"net"

	"github.com/kolesnikovm/messenger/configs"
	messageGrpc "github.com/kolesnikovm/messenger/internal/controller/message/grpc"
	messageUseCase "github.com/kolesnikovm/messenger/internal/usecase/message"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
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

		s := grpc.NewServer()
		messageUseCase := messageUseCase.New()
		messageGrpc.NewMessageServerGrpc(s, messageUseCase)

		log.Info().Msgf("Messenger server listening on %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("failed to start grpc server")
		}
	},
}
