package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/kolesnikovm/messenger/configs"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

		config, err := configs.NewServerConfig(configFile)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to instantiate config")
		}

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ListenPort))
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to listen port %d", config.ListenPort)
		}

		ctx, cancel := context.WithCancel(context.Background())

		app, cleanup, err := InitializeApplication(config)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to initialize application")
		}
		defer func() {
			cancel()
			cleanup()
		}()

		app.notifier.Start(ctx)

		app.archiver.Start(ctx)

		go func() {
			log.Info().Msgf("Messenger server listening on %v", lis.Addr())
			if err := app.grpcServer.Serve(lis); err != nil {
				log.Fatal().Err(err).Msg("failed to start grpc server")
			}
		}()

		metricsMux := http.NewServeMux()
		metricsMux.Handle("/metrics", promhttp.Handler())
		metricsMux.HandleFunc("/debug/pprof/", pprof.Index)
		metricsMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		metricsMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		metricsMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		metricsMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		metricsSrv := http.Server{Addr: config.MetricsAddress, Handler: metricsMux}

		go func() {
			log.Info().Msgf("metrics server listening on %s", metricsSrv.Addr)
			if err := metricsSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal().Err(err).Msg("metrics server failure")
			}
		}()

		signal := <-stop
		log.Info().Str("signal", signal.String()).Msg("received os signal")

		log.Info().Msg("server stopped")
	},
}
