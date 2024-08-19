package gateway

type lambdaGatewayOption func(*LambdaGateway)

// WithRuntime sets the lambda runtime handler for the LambdaGateway
func WithRuntime(runtime LambdaRuntimeHandler) lambdaGatewayOption {
	return func(g *LambdaGateway) {
		g.runtime = runtime
	}
}

// WithRouter sets the lambda event router for the LambdaGateway
func WithRouter(router LambdaEventRouter) lambdaGatewayOption {
	return func(g *LambdaGateway) {
		g.routeEvent = router
	}
}
