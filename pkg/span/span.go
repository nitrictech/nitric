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
	"net/http"
	"os"

	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	FunctionName          = "unknown"
	UseFuncNameAsSpanName = true
)

func init() {
	// support for aws, gcp and azure
	for _, en := range []string{"AWS_LAMBDA_FUNCTION_NAME", "K_SERVICE", "WEBSITE_SITE_NAME"} {
		name := os.Getenv(en)
		if name != "" {
			FunctionName = name
			break
		}
	}
}

func FromFastHttp(fhCtx *fasthttp.RequestCtx) trace.Span {
	headerCopy := propagation.HeaderCarrier{}

	fhCtx.Request.Header.VisitAll(func(key []byte, val []byte) {
		http.Header(headerCopy).Add(string(key), string(val))
	})

	// this extracts the traceID from the header and creates a parent span in the context.
	ctx := otel.GetTextMapPropagator().Extract(context.TODO(), headerCopy)

	_, span := otel.Tracer("membrane/pkg/span").Start(ctx, string(fhCtx.URI().PathOriginal()))

	if UseFuncNameAsSpanName {
		span.SetName(FunctionName)
	}

	span.SetAttributes(
		semconv.HTTPMethodKey.String(string(fhCtx.Method())),
		semconv.HTTPTargetKey.String(string(fhCtx.URI().PathOriginal())),
		semconv.HTTPURLKey.String(fhCtx.URI().String()),
	)

	return span
}

func FromContext(ctx context.Context, spanName string) trace.Span {
	_, span := otel.Tracer("membrane/pkg/span").Start(ctx, spanName)

	if UseFuncNameAsSpanName {
		span.SetName(FunctionName)
	}

	return span
}
