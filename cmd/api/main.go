package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/sawdustofmind/dataflow-shop-statistics/internal/app"
	"github.com/sawdustofmind/dataflow-shop-statistics/internal/config"
	"github.com/sawdustofmind/dataflow-shop-statistics/internal/log"
)

func main() {
	if err := apiCmd.Execute(); err != nil {
		outputInitError(err)
	}
}

func outputInitError(err error) {
	fmt.Println(`{"error":"` + err.Error() + `"}`)
	os.Exit(-1)
}

func gracefulShutdown(ctx context.Context) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case s := <-signalChan:
		log.Info("signal os", zap.String("", s.String()))
	case <-ctx.Done():
		log.Error(ctx.Err().Error())
	}
}

var apiCmd = &cobra.Command{
	Use:           "api",
	Short:         "api",
	SilenceUsage:  true,
	SilenceErrors: true,
	Version:       "1.0.0",

	RunE: func(cmd *cobra.Command, args []string) error {

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		conf, err := config.Init()
		if err != nil {
			return err
		}

		apiApp, err := app.CreateApiApplication(conf)
		if err != nil {
			return err
		}
		defer apiApp.Close() // nolint: staticcheck

		apiApp.Run(ctx)

		gracefulShutdown(ctx)
		return nil
	},
}
