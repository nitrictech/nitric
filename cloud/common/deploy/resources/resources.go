package resources

type ResourceType string

const (
	API           ResourceType = "api"
	Bucket                     = "bucket"
	Collection                 = "collection"
	ExecutionUnit              = "execution-unit"
	HttpProxy                  = "http-proxy"
	Policy                     = "policy"
	Queue                      = "queue"
	Schedule                   = "schedule"
	Secret                     = "secret"
	Stack                      = "stack"
	Topic                      = "topic"
	Websocket                  = "websocket"
)
