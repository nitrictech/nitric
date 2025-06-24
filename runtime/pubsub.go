package runtime

import (
	"context"

	pubsubpb "github.com/nitrictech/nitric/proto/pubsub/v2"
	"github.com/nitrictech/nitric/runtime/pubsub"
)

func (s *GrpcServer) Publish(ctx context.Context, req *pubsubpb.PubsubPublishRequest) (*pubsubpb.PubsubPublishResponse, error) {
	// 1. Resolve the pubsub plugin to use
	plugin := pubsub.GetPlugin("default")

	// 2. Call the plugin
	return plugin.Publish(ctx, req)
}
