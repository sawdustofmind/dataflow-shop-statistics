package app

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/sawdustofmind/dataflow-shop-statistics/internal/config"
	userhttp "github.com/sawdustofmind/dataflow-shop-statistics/internal/http/user"
	"github.com/sawdustofmind/dataflow-shop-statistics/internal/log"
	"github.com/sawdustofmind/dataflow-shop-statistics/internal/service"
	"github.com/sawdustofmind/dataflow-shop-statistics/internal/storage"
)

type ApiApplication struct {
	servers struct {
		userApi *http.Server
	}

	Config *config.Config
}

func CreateApiApplication(conf *config.Config) (*ApiApplication, error) {
	log.Info("Starting p2p version", zap.Any("config", conf))

	app := &ApiApplication{
		Config: conf,
	}

	s := storage.NewInMemoryStorage(conf.Storage.InitialCapacity)

	dataService := service.NewDataService(s)
	statisticsService := service.NewStatisticsService(s)

	userRouter := gin.Default()
	userServerImpl := userhttp.NewServerImpl(dataService, statisticsService)
	userhttp.RegisterHandlers(userRouter, userServerImpl)
	app.servers.userApi = &http.Server{
		Addr:           conf.UserServer.Address,
		ReadTimeout:    conf.UserServer.ReadTimeout,
		WriteTimeout:   conf.UserServer.WriteTimeout,
		MaxHeaderBytes: conf.UserServer.MaxHeaderBytes,
		Handler:        userRouter,
	}
	return app, nil
}

func (app *ApiApplication) Run(ctx context.Context) {
	go func() {
		err := app.servers.userApi.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start user server", zap.Error(err))
		}
	}()

	log.Info("ready to work")
}

func (app *ApiApplication) Close() error {
	log.Info("Closing application")
	var err error

	err = app.servers.userApi.Close()
	if err != nil {
		err = multierr.Append(err, err)
	}

	return multierr.Combine(
		err,
		log.Sync(),
	)
}
