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

package worker

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/span"
)

type instrumentedWorker struct {
	Worker
}

var _ Worker = &instrumentedWorker{}

func InstrumentedWorkerFn(w Worker) Worker {
	return &instrumentedWorker{
		Worker: w,
	}
}

func (a *instrumentedWorker) tracerFromTrigger(ctx context.Context, trigger *v1.TriggerRequest) (context.Context, trace.Span, error) {
	if http := trigger.GetHttp(); http != nil {
		ctx, s := otel.Tracer("membrane/pkg/worker", trace.WithInstrumentationVersion(span.MembraneVersion)).
			Start(ctx, span.Name(http.Path))

		s.SetAttributes(
			semconv.CodeFunctionKey.String("HandleHttp"),
			semconv.HTTPMethodKey.String(http.Method),
			semconv.HTTPTargetKey.String(http.Path),
		)

		return ctx, s, nil
	} else if topic := trigger.GetTopic(); topic != nil {
		ctx, s := otel.Tracer("membrane/pkg/worker", trace.WithInstrumentationVersion(span.MembraneVersion)).
			Start(ctx, span.Name("topic-"+topic.Topic))

		s.SetAttributes(
			semconv.CodeFunctionKey.String("HandleEvent"),
			semconv.MessagingSystemKey.String("nitric"),
			semconv.MessagingDestinationKindTopic,
			semconv.MessagingDestinationKey.String(topic.Topic),
		)

		return ctx, s, nil
	}

	return nil, nil, fmt.Errorf("invalid trigger provided")
}

func (a *instrumentedWorker) HandleTrigger(ctx context.Context, trigger *v1.TriggerRequest) (*v1.TriggerResponse, error) {
	var s trace.Span

	ctx, s, err := a.tracerFromTrigger(ctx, trigger)
	if err != nil {
		return nil, err
	}

	defer s.End()

	resp, err := a.Worker.HandleTrigger(ctx, trigger)
	if err != nil {
		s.SetStatus(codes.Error, "Request Handler returned an error")
		s.RecordError(err)
	} else {
		s.SetStatus(codes.Ok, "Request Handled Successfully")
	}

	if http := resp.GetHttp(); http != nil {
		s.SetAttributes(semconv.HTTPStatusCodeKey.Int(int(http.Status)))
	}

	return resp, err
}
