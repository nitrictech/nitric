package sdk

type NitricEvent struct {
	RequestId   string                 `json:"requestId,omitempty"`
	PayloadType string                 `json:"payloadType,omitempty"`
	Payload     map[string]interface{} `json:"payload,omitempty"`
}

type NitricQueueItem struct {
	Event NitricEvent					`json:"event,omitempty"`
	LeaseId string						`json:"leaseId,omitempty"`
}
