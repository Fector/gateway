package event

type GatewayConnected struct {
	*Event
}

type GatewayConnectionFailed struct {
	*Event
}

type GatewayDisconnected struct {
	*Event
}

type GatewayAuthenticated struct {
	*Event
}

type GatewayAuthenticationFailed struct {
	*Event
}

type GatewayGenericNack struct {
	*Event
}
