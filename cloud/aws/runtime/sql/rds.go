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
