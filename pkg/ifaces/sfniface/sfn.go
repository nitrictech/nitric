package sfniface

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sfn"
)

type SFNAPI interface {
	StartExecution(ctx context.Context, params *sfn.StartExecutionInput, optFns ...func(*sfn.Options)) (*sfn.StartExecutionOutput, error)
}
