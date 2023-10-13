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
	"os"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/nitrictech/nitric/core/pkg/span"
	"github.com/nitrictech/nitric/core/pkg/utils"
)

func newTraceProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	span.FunctionName = os.Getenv("K_SERVICE")
	span.UseFuncNameAsSpanName = false

	exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithDetectors(gcp.NewDetector()),
		resource.WithAttributes(
			semconv.CloudProviderGCP,
			semconv.CloudPlatformGCPCloudRun,
			attribute.Key("component").String("Nitric membrane"),
			semconv.ServiceNameKey.String(span.FunctionName),
			semconv.ServiceNamespaceKey.String(utils.GetEnv("NITRIC_STACK_ID", "")),
		),
	)
	if err != nil {
		return nil, err
	}

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagator.CloudTraceFormatPropagator{},
			propagation.TraceContext{},
		))

	rate, err := utils.PercentFromIntString(utils.GetEnv("NITRIC_TRACE_SAMPLE_PERCENT", "10"))
	if err != nil {
		return nil, errors.WithMessagef(err, "NITRIC_TRACE_SAMPLE_PERCENT should be an int")
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(rate))),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exp),
	), nil
}
