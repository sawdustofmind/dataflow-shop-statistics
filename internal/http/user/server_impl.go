package userhttp

import "github.com/sawdustofmind/dataflow-shop-statistics/internal/service"

type ServerImpl struct {
	dataService       *service.DataService
	statisticsService *service.StatisticsService
}

var _ ServerInterface = (*ServerImpl)(nil)

func NewServerImpl(
	dataService *service.DataService,
	statisticsService *service.StatisticsService,
) *ServerImpl {
	return &ServerImpl{
		dataService:       dataService,
		statisticsService: statisticsService,
	}
}
