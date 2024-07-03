package job

import (
	queuespb "github.com/nitrictech/nitric/core/pkg/proto/queues/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	secretspb "github.com/nitrictech/nitric/core/pkg/proto/secrets/v1"
	sqlpb "github.com/nitrictech/nitric/core/pkg/proto/sql/v1"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
	topicspb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
)

type JobMembraneOption = func(*JobMembrane)

func WithTopicServer(srv topicspb.TopicsServer) JobMembraneOption {
	return func(o *JobMembrane) {
		o.topicServer = srv
	}
}

func WithStorageServer(srv storagepb.StorageServer) JobMembraneOption {
	return func(o *JobMembrane) {
		o.storageServer = srv
	}
}

func WithQueueServer(srv queuespb.QueuesServer) JobMembraneOption {
	return func(o *JobMembrane) {
		o.queueServer = srv
	}
}

func WithSecretsServer(srv secretspb.SecretManagerServer) JobMembraneOption {
	return func(o *JobMembrane) {
		o.secretServer = srv
	}
}

func WithSqlServer(srv sqlpb.SqlServer) JobMembraneOption {
	return func(o *JobMembrane) {
		o.sqlServer = srv
	}
}

func WithResourcesServer(srv resourcespb.ResourcesServer) JobMembraneOption {
	return func(o *JobMembrane) {
		o.resourcesServer = srv
	}
}
