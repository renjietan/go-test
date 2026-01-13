package tcp_client_enum

type ClientActionEnum string

const (
	ClientActionConnect   ClientActionEnum = "connect"
	ClientActionLogin     ClientActionEnum = "login"
	ClientActionHeartbeat ClientActionEnum = "heartbeat"
)
