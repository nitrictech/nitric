package sql_service

import (
	"context"
	"fmt"
	"os"

	sqlpb "github.com/nitrictech/nitric/core/pkg/proto/sql/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CloudSqlService struct {
	sqlpb.UnimplementedSqlServer
}

var _ sqlpb.SqlServer = (*CloudSqlService)(nil)

func (s *CloudSqlService) ConnectionString(ctx context.Context, req *sqlpb.SqlConnectionStringRequest) (*sqlpb.SqlConnectionStringResponse, error) {
	baseUrl := os.Getenv("NITRIC_DATABASE_BASE_URL")

	if baseUrl == "" {
		return nil, status.Error(codes.FailedPrecondition, "NITRIC_DATABASE_BASE_URL environment variable not set")
	}

	connectionString := fmt.Sprintf("%s/%s", baseUrl, req.DatabaseName)

	return &sqlpb.SqlConnectionStringResponse{
		ConnectionString: connectionString,
	}, nil
}

func New() *CloudSqlService {
	return &CloudSqlService{}
}
