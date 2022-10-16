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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/nitrictech/nitric/pkg/triggers"
)

type instrumentedWorker struct {
	Worker
	span        trace.Span
	setName     bool
	setHTTPAttr bool
}

var _ Worker = &instrumentedWorker{}

func InstrumentedWorkerFn(span trace.Span, setName, setHTTPAttr bool) func(Worker) Worker {
	return func(w Worker) Worker {
		return &instrumentedWorker{
			Worker:      w,
			span:        span,
			setName:     setName,
			setHTTPAttr: setHTTPAttr,
		}
	}
}

// HandleEvent implements worker.Adapter
func (a *instrumentedWorker) HandleEvent(trigger *triggers.Event) error {
	a.span.SetAttributes(
		semconv.CodeFunctionKey.String("HandleEvent"),
		semconv.FaaSTriggerKey.String("event"),
		semconv.MessagingDestinationKindTopic,
		semconv.MessagingDestinationKey.String(trigger.Topic),
		semconv.MessagingMessageIDKey.String(trigger.ID),
		semconv.MessagingOperationProcess,
	)

	if a.setName {
		a.span.SetName(trigger.Topic)
	}

	defer a.span.End()

	// nowhere to inject the traceID into here :-(

	err := a.Worker.HandleEvent(trigger)
	if err != nil {
		a.span.RecordError(err)
	}

	return err
}

// HandleHttpRequest implements worker.Adapter
func (a *instrumentedWorker) HandleHttpRequest(trigger *triggers.HttpRequest) (*triggers.HttpResponse, error) {
	a.span.SetAttributes(
		semconv.FaaSTriggerHTTP,
		semconv.CodeFunctionKey.String("HandleHttpRequest"),
	)

	if a.setHTTPAttr {
		a.span.SetAttributes(
			semconv.HTTPMethodKey.String(trigger.Method),
			semconv.HTTPURLKey.String(trigger.Path),
		)
	}

	if a.setName {
		a.span.SetName(trigger.Path)
	}

	defer a.span.End()

	// Inject the correct headers.
	var hc propagation.HeaderCarrier = trigger.Header

	otel.GetTextMapPropagator().Inject(trace.ContextWithSpan(context.TODO(), a.span), hc)

	trigger.Header = hc

	resp, err := a.Worker.HandleHttpRequest(trigger)
	if err != nil {
		a.span.RecordError(err)
	}

	if resp != nil {
		a.span.SetAttributes(semconv.HTTPStatusCodeKey.Int(resp.StatusCode))
	}

	return resp, err
}
