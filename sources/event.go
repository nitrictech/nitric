package sources

// Event - A nitric event that has come from a trigger source
type Event struct {
	ID      string
	Topic   string
	Payload []byte
}

func (*Event) GetSourceType() SourceType {
	return SourceType_Subscription
}
