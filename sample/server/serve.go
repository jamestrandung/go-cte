package server

import (
	"github.com/jamestrandung/go-cte/sample/config"
	"github.com/jamestrandung/go-cte/sample/dependencies/configsfetcher"
	"github.com/jamestrandung/go-cte/sample/dependencies/mapservice"
	"github.com/jamestrandung/go-cte/sample/service/components/costconfigs"
	"github.com/jamestrandung/go-cte/sample/service/components/platformfee"
	"github.com/jamestrandung/go-cte/sample/service/components/quote"
	"github.com/jamestrandung/go-cte/sample/service/components/streaming"
	"github.com/jamestrandung/go-cte/sample/service/components/travelcost"
	"github.com/jamestrandung/go-cte/sample/service/components/travelplan"
	"github.com/jamestrandung/go-cte/sample/service/components/vat"
)

type dependencies struct {
	configsFetcher configsfetcher.Fetcher
	mapService     mapservice.Service
}

func (d dependencies) GetConfigsFetcher() configsfetcher.Fetcher {
	return d.configsFetcher
}

func (d dependencies) GetMapService() mapservice.Service {
	return d.mapService
}

var Dependencies dependencies

func Serve() {
	Dependencies = dependencies{
		configsFetcher: configsfetcher.Fetcher{},
		mapService:     mapservice.Service{},
	}

	config.Engine.RegisterComputers(
		costconfigs.Computer{},
		platformfee.Computer{},
		quote.Computer{},
		streaming.Computer{},
		travelcost.Computer{},
		travelplan.Computer{},
		vat.Computer{},
	)
}
