package api

import (
	"github.com/DeathHand/gateway/router"
	"github.com/DeathHand/gateway/storage"
)

type Api struct {
	Router         *router.Router
	GatewayStorage *storage.Gateway
}
