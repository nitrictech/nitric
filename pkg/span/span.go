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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/nitrictech/nitric/pkg/utils"
)

var (
	FunctionName          = "not-set"
	UseFuncNameAsSpanName = true
	MembraneVersion       = utils.GetEnv("MEMBRANE_VERSION", "not-set")
)

func FromHeaders(ctx context.Context, spanName string, headers map[string][]string) trace.Span {
	var hc propagation.HeaderCarrier = headers

	// this extracts the traceID from the header and creates a parent span in the context.
	_, sp := otel.Tracer("membrane/pkg/span", trace.WithInstrumentationVersion(MembraneVersion)).
		Start(otel.GetTextMapPropagator().Extract(ctx, hc), spanName)

	if UseFuncNameAsSpanName {
		sp.SetName(FunctionName)
	}

	return sp
}

func ToHeaders(ctx context.Context, headers map[string][]string) map[string][]string {
	var hc propagation.HeaderCarrier = headers

	otel.GetTextMapPropagator().Inject(ctx, hc)

	return hc
}
