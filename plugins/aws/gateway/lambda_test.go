package main_test

import (
	"context"
	"encoding/json"
	"fmt"

	events "github.com/aws/aws-lambda-go/events"
	"github.com/nitric-dev/membrane/plugins/sdk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	plugin "github.com/nitric-dev/membrane/plugins/aws/gateway"
)

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
				}},
			}

			client, _ := plugin.NewWithRuntime(runtime.Start)

			capturedRequests := make([]*sdk.NitricRequest, 0)
			// This function will block which means we don't need to wait on processing,
			// the function will unblock once processing has finished, this is due to our mock
			// handler only looping once over each request

			It("The gateway should translate into a standard NitricRequest", func() {
				client.Start(func(request *sdk.NitricRequest) *sdk.NitricResponse {
					capturedRequests = append(capturedRequests, request)
					// store the Nitric requests and prepare them for comparison with our expectations
					return &sdk.NitricResponse{
						Headers: map[string]string{
							"Content-Type": "text/plain",
						},
						Status: 200,
						Body:   []byte("Test Response"),
					}
				})

				Expect(len(capturedRequests)).To(Equal(1))

				request := capturedRequests[0]
				ctx := request.Context

				Expect(request.Payload).To(BeEquivalentTo([]byte("Test Payload")))
				Expect(request.ContentType).To(Equal("text/plain"))
				Expect(ctx.RequestId).To(Equal("test-request-id"))
				Expect(ctx.PayloadType).To(Equal("TestPayload"))
				Expect(ctx.SourceType).To(Equal(sdk.Request))
				Expect(ctx.Source).To(Equal("Test"))
			})
		})
	})

	Context("SNS Events", func() {
		When("The Lambda Gateway recieves SNS events", func() {
			topicName := "MyTopic"
			eventPayload := map[string]interface{}{
				"test": "test",
			}
			eventBytes, _ := json.Marshal(&eventPayload)

			event := sdk.NitricEvent{
				RequestId:   "test-request-id",
				PayloadType: "test-payload",
				Payload:     eventPayload,
			}

			messageBytes, _ := json.Marshal(&event)

			runtime := MockLambdaRuntime{
				// Setup mock events for our runtime to process...
				eventQueue: []interface{}{&events.SNSEvent{
					Records: []events.SNSEventRecord{
						events.SNSEventRecord{
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
				capturedRequests := make([]*sdk.NitricRequest, 0)
				// This function will block which means we don't need to wait on processing,
				// the function will unblock once processing has finished, this is due to our mock
				// handler only looping once over each request
				client.Start(func(request *sdk.NitricRequest) *sdk.NitricResponse {
					capturedRequests = append(capturedRequests, request)
					// store the Nitric requests and prepare them for comparison with our expectations
					return &sdk.NitricResponse{
						Headers: map[string]string{
							"Content-Type": "text/plain",
						},
						Status: 200,
						Body:   []byte("Test Response"),
					}
				})

				Expect(len(capturedRequests)).To(Equal(1))

				request := capturedRequests[0]
				context := request.Context

				Expect(request.ContentType).To(Equal("application/json"))
				Expect(eventBytes).To(BeEquivalentTo(request.Payload))
				Expect(context.PayloadType).To(Equal("test-payload"))
				Expect(context.RequestId).To(Equal("test-request-id"))
				Expect(context.SourceType).To(Equal(sdk.Subscription))
				Expect(context.Source).To(Equal(topicName))
			})
		})
	})
})
