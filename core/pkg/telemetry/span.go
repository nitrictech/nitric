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

package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/nitrictech/nitric/core/pkg/env"
)

var (
	FunctionName          = "not-set"
	UseFuncNameAsSpanName = true
	MembraneVersion       = env.MEMBRANE_VERSION.String()
)

func Name(n string) string {
	if !UseFuncNameAsSpanName {
		return n
	}

	return FunctionName
}

func FromHeaders(ctx context.Context, headers map[string][]string) context.Context {
	var hc propagation.HeaderCarrier = headers

	return otel.GetTextMapPropagator().Extract(ctx, hc)
}

// func ToTraceContext(ctx context.Context) *pb.TraceContext {
// 	hc := propagation.MapCarrier{}

// 	// we want to inject cloud agnostic info here, so that the user process
// 	// can do the same.
// 	propagation.TraceContext{}.Inject(ctx, hc)

// 	return &pb.TraceContext{Values: hc}
// }
