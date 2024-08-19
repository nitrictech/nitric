package membrane

import (
	"github.com/nitrictech/nitric/core/pkg/gateway"
	kvstorepb "github.com/nitrictech/nitric/core/pkg/proto/kvstore/v1"
	queuespb "github.com/nitrictech/nitric/core/pkg/proto/queues/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	secretspb "github.com/nitrictech/nitric/core/pkg/proto/secrets/v1"
	sqlpb "github.com/nitrictech/nitric/core/pkg/proto/sql/v1"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
	topicspb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
	websocketspb "github.com/nitrictech/nitric/core/pkg/proto/websockets/v1"
	"github.com/nitrictech/nitric/core/pkg/workers/apis"
	"github.com/nitrictech/nitric/core/pkg/workers/http"
	"github.com/nitrictech/nitric/core/pkg/workers/schedules"
	"github.com/nitrictech/nitric/core/pkg/workers/storage"
	"github.com/nitrictech/nitric/core/pkg/workers/topics"
	"github.com/nitrictech/nitric/core/pkg/workers/websockets"
)

type RuntimeServerOption func(opts *Membrane)

func WithResourcesPlugin(resources resourcespb.ResourcesServer) RuntimeServerOption {
	return func(opts *Membrane) {
		opts.ResourcesPlugin = resources
	}
}

func WithGatewayPlugin(gw gateway.GatewayService) RuntimeServerOption {
	return func(opts *Membrane) {
		opts.GatewayPlugin = gw
	}
}

func WithKeyValuePlugin(kv kvstorepb.KvStoreServer) RuntimeServerOption {
	return func(opts *Membrane) {
		opts.KeyValuePlugin = kv
	}
}

func WithTopicsPlugin(tp topicspb.TopicsServer) RuntimeServerOption {
	return func(opts *Membrane) {
		opts.TopicsPlugin = tp
	}
}

func WithStoragePlugin(sp storagepb.StorageServer) RuntimeServerOption {
	return func(opts *Membrane) {
		opts.StoragePlugin = sp
	}
}

func WithSecretManagerPlugin(sm secretspb.SecretManagerServer) RuntimeServerOption {
	return func(opts *Membrane) {
		opts.SecretManagerPlugin = sm
	}
}

func WithWebsocketPlugin(ws websocketspb.WebsocketServer) RuntimeServerOption {
	return func(opts *Membrane) {
		opts.WebsocketPlugin = ws
	}
}

func WithQueuesPlugin(qs queuespb.QueuesServer) RuntimeServerOption {
	return func(opts *Membrane) {
		opts.QueuesPlugin = qs
	}
}

func WithSqlPlugin(sql sqlpb.SqlServer) RuntimeServerOption {
	return func(opts *Membrane) {
		opts.SqlPlugin = sql
	}
}

func WithApiPlugin(api apis.ApiRequestHandler) RuntimeServerOption {
	return func(opts *Membrane) {
		opts.ApiPlugin = api
	}
}

func WithHttpPlugin(http http.HttpRequestHandler) RuntimeServerOption {
	return func(opts *Membrane) {
		opts.HttpPlugin = http
	}
}

func WithSchedulesPlugin(schedules schedules.ScheduleRequestHandler) RuntimeServerOption {
	return func(opts *Membrane) {
		opts.SchedulesPlugin = schedules
	}
}

func WithTopicsListenerPlugin(topics topics.SubscriptionRequestHandler) RuntimeServerOption {
	return func(opts *Membrane) {
		opts.TopicsListenerPlugin = topics
	}
}

func WithStorageListenerPlugin(storage storage.BucketRequestHandler) RuntimeServerOption {
	return func(opts *Membrane) {
		opts.StorageListenerPlugin = storage
	}
}

func WithWebsocketListenerPlugin(websockets websockets.WebsocketRequestHandler) RuntimeServerOption {
	return func(opts *Membrane) {
		opts.WebsocketListenerPlugin = websockets
	}
}

func WithServiceAddress(address string) RuntimeServerOption {
	return func(opts *Membrane) {
		opts.ServiceAddress = address
	}
}

// WithMinWorkers - Set the minimum number of workers that need to be available.
// this option is ignored if the MIN_WORKERS environment variable is set
func WithMinWorkers(minWorkers int) RuntimeServerOption {
	return func(opts *Membrane) {
		opts.MinWorkers = minWorkers
	}
}

func WithChildCommand(command []string) RuntimeServerOption {
	return func(opts *Membrane) {
		opts.ChildCommand = command
	}
}

func WithPreCommands(commands [][]string) RuntimeServerOption {
	return func(opts *Membrane) {
		opts.PreCommands = commands
	}
}

func WithChildTimeoutSeconds(timeout int) RuntimeServerOption {
	return func(opts *Membrane) {
		opts.ChildTimeoutSeconds = timeout
	}
}
