// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package batch

import (
	"context"
	"fmt"
	"io"

	batch "cloud.google.com/go/batch/apiv1"
	gcpbatchpb "cloud.google.com/go/batch/apiv1/batchpb"
	"cloud.google.com/go/storage"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/env"
	batchpb "github.com/nitrictech/nitric/core/pkg/proto/batch/v1"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/encoding/protojson"
)

type GcpBatchService struct {
	projectId      string
	region         string
	batchClient    *batch.Client
	storageClient  *storage.Client
	jobsBucketName string
	batchpb.UnimplementedBatchServer
}

func (a *GcpBatchService) SubmitJob(ctx context.Context, request *batchpb.JobSubmitRequest) (*batchpb.JobSubmitResponse, error) {
	// read the job definition of the GCP jobs bucket
	fmt.Println("submitting job")
	jobDefReader, err := a.storageClient.Bucket(a.jobsBucketName).Object(fmt.Sprintf("%s.json", request.JobName)).NewReader(ctx)
	if err != nil {
		fmt.Printf("Error reading job definition from bucket: %+v\n", err)
		return nil, err
	}

	jobContents, err := io.ReadAll(jobDefReader)
	if err != nil {
		return nil, err
	}

	jobDefinition := &gcpbatchpb.Job{}
	err = protojson.Unmarshal(jobContents, jobDefinition)
	if err != nil {
		return nil, err
	}

	jobData, err := protojson.Marshal(request.Data)
	if err != nil {
		return nil, err
	}

	// Add job data to environment variables
	jobDefinition.TaskGroups[0].TaskSpec.Environment.Variables["NITRIC_JOB_DATA"] = string(jobData)

	_, err = a.batchClient.CreateJob(ctx, &gcpbatchpb.CreateJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", a.projectId, a.region),
		Job:    jobDefinition,
	})
	if err != nil {
		fmt.Printf("Error creating job: %+v\n", err)
		return nil, err
	}

	return &batchpb.JobSubmitResponse{}, nil
}

func New() (*GcpBatchService, error) {
	credentials, credentialsError := google.FindDefaultCredentials(context.TODO(),
		storage.ScopeReadWrite,
		// required for signing blob urls
		iamcredentials.CloudPlatformScope,
	)
	if credentialsError != nil {
		return nil, fmt.Errorf("GCP credentials error: %w", credentialsError)
	}

	batchClient, err := batch.NewClient(context.TODO(), option.WithCredentials(credentials))
	if err != nil {
		return nil, err
	}

	storageClient, err := storage.NewClient(context.TODO(), option.WithCredentials(credentials))
	if err != nil {
		return nil, err
	}

	projectId := env.GOOGLE_PROJECT_ID.String()
	region := env.GCP_REGION.String()
	jobsBucket := env.JOBS_BUCKET_NAME.String()

	return &GcpBatchService{
		jobsBucketName: jobsBucket,
		projectId:      projectId,
		region:         region,
		batchClient:    batchClient,
		storageClient:  storageClient,
	}, nil
}
