// Copyright 2021 Nitric Pty Ltd.
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

package main

import (
	"context"
	"strconv"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/pkg/errors"
	lambdadetector "go.opentelemetry.io/contrib/detectors/aws/lambda"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/nitrictech/nitric/core/pkg/telemetry"

	commonenv "github.com/nitrictech/nitric/cloud/common/runtime/env"
)

// PercentFromIntString returns a float between 0.0 to 1 representing a percentage.
// this is converted from a string int in the range "0" to "100".
func decimalFromPercentIntString(in string) (float64, error) {
	intVar, err := strconv.Atoi(in)
	if err != nil {
		return 0, err
	}

	if intVar >= 100 {
		return 1, nil
	} else if intVar <= 0 {
		return 0, nil
	}

	return float64(intVar) / float64(100), nil
}

func newTracerProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	telemetry.FunctionName = lambdacontext.FunctionName
	telemetry.UseFuncNameAsSpanName = true

	exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithDetectors(lambdadetector.NewResourceDetector()),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.CloudProviderAWS,
			semconv.CloudPlatformAWSLambda,
			semconv.ServiceNameKey.String(telemetry.FunctionName),
			semconv.ServiceNamespaceKey.String(commonenv.NITRIC_STACK_ID.String()),
		),
	)
	if err != nil {
		return nil, err
	}

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			xray.Propagator{},
			propagation.TraceContext{},
		))

	rate, err := decimalFromPercentIntString(commonenv.NITRIC_TRACE_SAMPLE_PERCENT.String())
	if err != nil {
		return nil, errors.WithMessagef(err, "NITRIC_TRACE_SAMPLE_PERCENT should be an int")
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(rate))),
		sdktrace.WithBatcher(exp),
		sdktrace.WithIDGenerator(xray.NewIDGenerator()),
		sdktrace.WithResource(res),
	), nil
}
