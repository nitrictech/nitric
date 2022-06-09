package base_http

import (
	"github.com/fasthttp/router"
	"github.com/nitrictech/nitric/pkg/worker"
	"github.com/valyala/fasthttp"
)

type BaseHttpGatewayOption = func(*BaseHttpGateway)

func WithMiddleware(mw func(*fasthttp.RequestCtx, worker.WorkerPool) bool) BaseHttpGatewayOption {
	return func(bhg *BaseHttpGateway) {
		bhg.mw = mw
	}
}

func WithRouter(routeRegister func(*router.Router, worker.WorkerPool)) BaseHttpGatewayOption {
	return func(bhg *BaseHttpGateway) {
		bhg.routeReg = routeRegister
	}
}
