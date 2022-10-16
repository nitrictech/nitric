// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lambda_service

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/lambdacontext"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/nitrictech/nitric/pkg/span"
)

func spanFromContext(ctx context.Context) trace.Span {
	span := span.FromContext(ctx, span.FunctionName)

	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		log.Default().Println("failed to load lambda context from context, ensure tracing enabled in Lambda")
	}

	if lc != nil {
		ctxRequestID := lc.AwsRequestID
		span.SetAttributes(semconv.FaaSExecutionKey.String(ctxRequestID))

		// Some resource attrs added as span attrs because lambda
		// resource detectors are created before a lambda
		// invocation and therefore lack lambdacontext.
		// Create these attrs upon first invocation
		ctxFunctionArn := lc.InvokedFunctionArn
		span.SetAttributes(semconv.FaaSIDKey.String(ctxFunctionArn))
		arnParts := strings.Split(ctxFunctionArn, ":")
		if len(arnParts) >= 5 {
			span.SetAttributes(semconv.CloudAccountIDKey.String(arnParts[4]))
		}
	}

	return span
}
