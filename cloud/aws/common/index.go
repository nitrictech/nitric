package common

// AwsResourceName - Provides a type hint for the mapping of Nitric resource names to AWS resource names
type AwsResourceName = string

// AwsResourceArn - Provides a type hint for the mapping of Nitric resource names to AWS resource ARNs
type AwsResourceArn = string

type ApiGateway struct {
	Arn      string `json:"arn"`
	Endpoint string `json:"endpoint"`
}

// ResourceIndex - The resource index for a nitric stack
type ResourceIndex struct {
	// Buckets - The S3 Buckets
	// This is a map of Nitric Name to AWS Bucket name
	Buckets    map[string]AwsResourceName `json:"buckets"`
	Topics     map[string]string          `json:"topics"`
	KvStores   map[string]string          `json:"kvStores"`
	Queues     map[string]string          `json:"queues"`
	Secrets    map[string]string          `json:"secrets"`
	Apis       map[string]ApiGateway      `json:"apis"`
	Websockets map[string]ApiGateway      `json:"websockets"`
}

func NewResourceIndex() *ResourceIndex {
	return &ResourceIndex{
		Buckets:    make(map[string]AwsResourceName),
		Topics:     make(map[string]string),
		KvStores:   make(map[string]string),
		Queues:     make(map[string]string),
		Secrets:    make(map[string]string),
		Apis:       make(map[string]ApiGateway),
		Websockets: make(map[string]ApiGateway),
	}
}
