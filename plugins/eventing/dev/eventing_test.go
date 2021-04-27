package eventing_service_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	eventing_plugin "github.com/nitric-dev/membrane/plugins/eventing/dev"
	"github.com/nitric-dev/membrane/sdk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockHttpClient struct {
	eventing_plugin.LocalHttpEventingClient
	capturedRequests []*http.Request
}

func (m *MockHttpClient) reset() {
	m.capturedRequests = make([]*http.Request, 0)
}

func (m *MockHttpClient) Do(request *http.Request) (*http.Response, error) {
	if m.capturedRequests == nil {
		m.capturedRequests = make([]*http.Request, 0)
	}

	// Capture the request for assertion
	m.capturedRequests = append(m.capturedRequests, request)

	// Our dev handler currently doesn't care about failure...
	// or even look at the response...
	return &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
	}, nil
}

var _ = Describe("Eventing", func() {
	mockHttpClient := &MockHttpClient{}

	AfterEach(func() {
		mockHttpClient.reset()
	})

	When("Getting available topics", func() {

		When("topics exist", func() {
			subs := map[string][]string{
				"test": {"http://test-endpoint/"},
			}

			pubsubClient, _ := eventing_plugin.NewWithClientAndSubs(mockHttpClient, subs)

			It("Should return the available topics", func() {
				topics, err := pubsubClient.ListTopics()
				Expect(err).To(BeNil())
				Expect(topics).To(ContainElement("test"))
			})
		})

		When("no topics exist", func() {
			subs := map[string][]string{}
			pubsubClient, _ := eventing_plugin.NewWithClientAndSubs(mockHttpClient, subs)

			It("Should return the no topics", func() {
				topics, err := pubsubClient.ListTopics()

				Expect(err).To(BeNil())
				Expect(topics).To(HaveLen(0))
			})
		})
	})

	When("Publishing an event", func() {
		testPayload := map[string]interface{}{
			"Test": "test",
		}
		testEvent := &sdk.NitricEvent{
			ID:          "1234",
			PayloadType: "Test-Payload",
			Payload:     testPayload,
		}

		When("The target topic is not available", func() {
			subs := map[string][]string{}
			pubsubClient, _ := eventing_plugin.NewWithClientAndSubs(mockHttpClient, subs)

			It("should return an error", func() {
				err := pubsubClient.Publish("test", testEvent)
				Expect(err).ToNot(BeNil())
			})
		})

		When("The target topic is available", func() {
			subs := map[string][]string{
				"test": {"http://test-endpoint/"},
			}

			pubsubClient, _ := eventing_plugin.NewWithClientAndSubs(mockHttpClient, subs)

			It("should successfully publish", func() {
				err := pubsubClient.Publish("test", testEvent)

				By("Not returning an error")
				Expect(err).To(BeNil())

				By("Publishing to the only configured endpoint")
				Expect(mockHttpClient.capturedRequests).To(HaveLen(1))

				capturedRequest := mockHttpClient.capturedRequests[0]
				By("Publishing to the given endpoint from subs")
				Expect(capturedRequest.Host).To(Equal("test-endpoint"))

				By("Providing the event RequestId in headers")
				Expect(capturedRequest.Header.Get("x-nitric-request-id")).To(Equal("1234"))

				By("Providing the event PayloadType in headers")
				Expect(capturedRequest.Header.Get("x-nitric-payload-type")).To(Equal("Test-Payload"))

				By("Providing the sourceType in header as SUBSCRIPTION")
				Expect(capturedRequest.Header.Get("x-nitric-source-type")).To(Equal("SUBSCRIPTION"))

				By("Providing the source in header as the name of the topic")
				Expect(capturedRequest.Header.Get("x-nitric-source")).To(Equal("test"))

				By("Providing the payload in the Body")
				bodyBytes, _ := ioutil.ReadAll(capturedRequest.Body)
				bodyMap := make(map[string]interface{})
				json.Unmarshal(bodyBytes, &bodyMap)
				Expect(bodyMap).To(BeEquivalentTo(testPayload))

			})
		})
	})
})
