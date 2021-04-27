package lambda_service_test

import (
	"context"
	"encoding/json"
	"fmt"

	events "github.com/aws/aws-lambda-go/events"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/triggers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	plugin "github.com/nitric-dev/membrane/plugins/gateway/lambda"
)

type MockTriggerHandler struct {
	httpRequests []*triggers.HttpRequest
	events       []*triggers.Event
}

func (m *MockTriggerHandler) HandleEvent(trigger *triggers.Event) error {
	if m.events == nil {
		m.events = make([]*triggers.Event, 0)
	}
	m.events = append(m.events, trigger)

	return nil
}

func (m *MockTriggerHandler) HandleHttpRequest(trigger *triggers.HttpRequest) (*triggers.HttpResponse, error) {
	if m.httpRequests == nil {
		m.httpRequests = make([]*triggers.HttpRequest, 0)
	}
	m.httpRequests = append(m.httpRequests, trigger)

	return &triggers.HttpResponse{
		StatusCode: 200,
		Body:       []byte("Mock Handled!"),
	}, nil
}

func (m *MockTriggerHandler) reset() {
	m.httpRequests = make([]*triggers.HttpRequest, 0)
	m.events = make([]*triggers.Event, 0)
}

type MockLambdaRuntime struct {
	plugin.LambdaRuntimeHandler
	// FIXME: Make this a union array of stuff to send....
	eventQueue []interface{}
}

func (m *MockLambdaRuntime) Start(handler interface{}) {
	// cast the function type to what we know it will be
	typedFunc := handler.(func(ctx context.Context, event plugin.Event) (interface{}, error))
	for _, event := range m.eventQueue {

		bytes, _ := json.Marshal(event)
		evt := plugin.Event{}

		json.Unmarshal(bytes, &evt)
		// Unmarshal the thing into the event type we expect...
		// TODO: Do something with out results here...
		_, err := typedFunc(context.TODO(), evt)

		if err != nil {
			// Print the error?
		}
	}
}

var _ = Describe("Lambda", func() {
	mockHandler := &MockTriggerHandler{}
	AfterEach(func() {
		mockHandler.reset()
	})

	Context("Http Events", func() {
		When("Sending a compliant HTTP Event", func() {

			runtime := MockLambdaRuntime{
				// Setup mock events for our runtime to process...
				eventQueue: []interface{}{&events.APIGatewayV2HTTPRequest{
					Headers: map[string]string{
						"User-Agent":            "Test",
						"x-nitric-payload-type": "TestPayload",
						"x-nitric-request-id":   "test-request-id",
						"Content-Type":          "text/plain",
					},
					RawPath: "/test/test",
					Body:    "Test Payload",
					RequestContext: events.APIGatewayV2HTTPRequestContext{
						HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
							Method: "GET",
						},
					},
				}},
			}

			client, _ := plugin.NewWithRuntime(runtime.Start)

			// This function will block which means we don't need to wait on processing,
			// the function will unblock once processing has finished, this is due to our mock
			// handler only looping once over each request
			It("The gateway should translate into a standard NitricRequest", func() {
				client.Start(mockHandler)

				By("Handling a single HTTP request")
				Expect(len(mockHandler.httpRequests)).To(Equal(1))

				request := mockHandler.httpRequests[0]

				By("Retaining the body")
				Expect(string(request.Body)).To(BeEquivalentTo("Test Payload"))
				By("Retaining the Headers")
				Expect(request.Header["User-Agent"]).To(Equal("Test"))
				Expect(request.Header["x-nitric-payload-type"]).To(Equal("TestPayload"))
				Expect(request.Header["x-nitric-request-id"]).To(Equal("test-request-id"))
				Expect(request.Header["Content-Type"]).To(Equal("text/plain"))
				By("Retaining the method")
				Expect(request.Method).To(Equal("GET"))
				By("Retaining the path")
				Expect(request.Path).To(Equal("/test/test"))
			})
		})
	})

	Context("SNS Events", func() {
		When("The Lambda Gateway recieves SNS events", func() {
			topicName := "MyTopic"
			eventPayload := map[string]interface{}{
				"test": "test",
			}
			// eventBytes, _ := json.Marshal(&eventPayload)

			event := sdk.NitricEvent{
				ID:          "test-request-id",
				PayloadType: "test-payload",
				Payload:     eventPayload,
			}

			messageBytes, _ := json.Marshal(&event)

			runtime := MockLambdaRuntime{
				// Setup mock events for our runtime to process...
				eventQueue: []interface{}{&events.SNSEvent{
					Records: []events.SNSEventRecord{
						{
							EventVersion:         "",
							EventSource:          "aws:sns",
							EventSubscriptionArn: "some:arbitrary:subscription:arn:MySubscription",
							SNS: events.SNSEntity{
								TopicArn: fmt.Sprintf("some:arbitrary:topic:arn:%s", topicName),
								Message:  string(messageBytes),
							},
						},
					},
				}},
			}

			client, _ := plugin.NewWithRuntime(runtime.Start)

			It("The gateway should translate into a standard NitricRequest", func() {
				// This function will block which means we don't need to wait on processing,
				// the function will unblock once processing has finished, this is due to our mock
				// handler only looping once over each request
				client.Start(mockHandler)

				By("Handling a single event")
				Expect(len(mockHandler.events)).To(Equal(1))

				request := mockHandler.events[0]

				By("Containing the Source Topic")
				Expect(request.Topic).To(Equal("MyTopic"))
			})
		})
	})
})
