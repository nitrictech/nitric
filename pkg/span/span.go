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

package span

import (
	"context"
	"net/textproto"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	pb "github.com/nitrictech/nitric/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/pkg/utils"
)

var (
	FunctionName          = "not-set"
	UseFuncNameAsSpanName = true
	MembraneVersion       = utils.GetEnv("MEMBRANE_VERSION", "not-set")
)

func FromHeaders(ctx context.Context, spanName string, headers map[string][]string) context.Context {
	var hc propagation.HeaderCarrier = headers

	if UseFuncNameAsSpanName {
		spanName = FunctionName
	}

	// this extracts the traceID from the header and creates a parent span in the context.
	ctx, _ = otel.Tracer("membrane/pkg/span", trace.WithInstrumentationVersion(MembraneVersion)).
		Start(otel.GetTextMapPropagator().Extract(ctx, hc), spanName)

	return ctx
}

// simpleHeaderCarrier adapts map[string]string to satisfy the TextMapCarrier interface.
type simpleHeaderCarrier map[string]string

func (hc simpleHeaderCarrier) Get(key string) string {
	return hc[textproto.CanonicalMIMEHeaderKey(key)]
}

func (hc simpleHeaderCarrier) Set(key string, value string) {
	hc[textproto.CanonicalMIMEHeaderKey(key)] = value
}

func (hc simpleHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range hc {
		keys = append(keys, k)
	}

	return keys
}

func ToTraceContext(ctx context.Context) *pb.TraceContext {
	var hc simpleHeaderCarrier = make(simpleHeaderCarrier)

	// we want to inject cloud agnostic info here, so that the user process
	// can do the same.
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	prop.Inject(ctx, hc)

	return &pb.TraceContext{Values: hc}
}
