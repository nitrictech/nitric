package snsiface

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNSAPI interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}
