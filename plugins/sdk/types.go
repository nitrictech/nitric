package sdk

// NitricEvent - An event for asynchronous processing and reactive programming
type NitricEvent struct {
	ID          string                 `json:"id,omitempty"`
	PayloadType string                 `json:"payloadType,omitempty"`
	Payload     map[string]interface{} `json:"payload,omitempty"`
}

// NitricTask - A task for asynchronous processing
type NitricTask struct {
	ID          string                 `json:"id,omitempty"`
	LeaseID     string                 `json:"leaseId,omitempty"`
	PayloadType string                 `json:"payloadType,omitempty"`
	Payload     map[string]interface{} `json:"payload,omitempty"`
}
