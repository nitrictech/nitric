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

	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/nitrictech/nitric/pkg/triggers"
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

// HandleEvent implements worker.Adapter
func (a *instrumentedWorker) HandleEvent(ctx context.Context, trigger *triggers.Event) error {
	s := trace.SpanFromContext(ctx)
	s.SetAttributes(
		semconv.CodeFunctionKey.String("HandleEvent"),
		semconv.MessagingDestinationKindTopic,
		semconv.MessagingDestinationKey.String(trigger.Topic),
		semconv.MessagingMessageIDKey.String(trigger.ID),
	)

	defer s.End()

	err := a.Worker.HandleEvent(trace.ContextWithSpan(ctx, s), trigger)
	if err != nil {
		s.SetStatus(codes.Error, "Event Handler returned an error")
		s.RecordError(err)
	} else {
		s.SetStatus(codes.Ok, "Event Handled Successfully")
	}

	return err
}

// HandleHttpRequest implements worker.Adapter
func (a *instrumentedWorker) HandleHttpRequest(ctx context.Context, trigger *triggers.HttpRequest) (*triggers.HttpResponse, error) {
	s := trace.SpanFromContext(ctx)
	s.SetAttributes(
		semconv.CodeFunctionKey.String("HandleHttpRequest"),
		semconv.HTTPMethodKey.String(trigger.Method),
		semconv.HTTPTargetKey.String(trigger.Path),
		semconv.HTTPURLKey.String(trigger.URL),
	)

	defer s.End()

	resp, err := a.Worker.HandleHttpRequest(trace.ContextWithSpan(ctx, s), trigger)
	if err != nil {
		s.SetStatus(codes.Error, "Request Handler returned an error")
		s.RecordError(err)
	} else {
		s.SetStatus(codes.Ok, "Request Handled Successfully")
	}

	if resp != nil {
		s.SetAttributes(semconv.HTTPStatusCodeKey.Int(resp.StatusCode))
	}

	return resp, err
}
