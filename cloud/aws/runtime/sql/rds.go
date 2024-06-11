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

package sql_service

import (
	"context"
	"fmt"
	"os"

	sqlpb "github.com/nitrictech/nitric/core/pkg/proto/sql/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RdsSqlService struct {
	sqlpb.UnimplementedSqlServer
}

var _ sqlpb.SqlServer = (*RdsSqlService)(nil)

func (s *RdsSqlService) ConnectionString(ctx context.Context, req *sqlpb.SqlConnectionStringRequest) (*sqlpb.SqlConnectionStringResponse, error) {
	baseUrl := os.Getenv("NITRIC_DATABASE_BASE_URL")

	if baseUrl == "" {
		return nil, status.Error(codes.FailedPrecondition, "NITRIC_DATABASE_BASE_URL environment variable not set")
	}

	return &sqlpb.SqlConnectionStringResponse{
		ConnectionString: fmt.Sprintf("%s/%s", baseUrl, req.DatabaseName),
	}, nil
}

func NewRdsSqlService() *RdsSqlService {
	return &RdsSqlService{}
}
