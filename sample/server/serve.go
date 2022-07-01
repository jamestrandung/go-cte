package server

import (
	"github.com/jamestrandung/go-die/sample/service/costconfigs"
	"github.com/jamestrandung/go-die/sample/service/costconfigs/dummy"
)

func Serve() {
	costconfigs.InitComputer(dummy.CostConfigsFetcher{})
}
