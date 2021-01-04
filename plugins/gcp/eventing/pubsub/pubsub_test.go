package pubsub_plugin_test

import (
	"context"
	"fmt"
	"net"
	"os"

	"cloud.google.com/go/pubsub"
	pubsub_plugin "github.com/nitric-dev/membrane/plugins/gcp/eventing/pubsub"
	mocks "github.com/nitric-dev/membrane/plugins/gcp/mocks"
	"github.com/nitric-dev/membrane/plugins/sdk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pb "google.golang.org/genproto/googleapis/pubsub/v1"
	"google.golang.org/grpc"
)

const GCP_PROJECT_NAME = "fake-project"

func createPubsubTopicId(name string) string {
	return fmt.Sprintf("projects/%s/topics/%s", GCP_PROJECT_NAME, name)
}

var _ = Describe("Pubsub Plugin", func() {
	var opts []grpc.ServerOption
	var pubsubClient *pubsub.Client
	var pubsubPlugin sdk.EventingPlugin
	grpcServer := grpc.NewServer(opts...)
	mockPubsubServer := mocks.NewPubsubPublisherServer([]string{})
	// Set the emulator host...
	os.Setenv("PUBSUB_EMULATOR_HOST", "127.0.0.1:50051")
	pb.RegisterPublisherServer(grpcServer, mockPubsubServer)
	lis, _ := net.Listen("tcp", "127.0.0.1:50051")
	// Do not block on serve...
	go (func() {
		grpcServer.Serve(lis)
	})()

	AfterEach(func() {
		mockPubsubServer.ClearMessages()
	})

	pubsubClient, _ = pubsub.NewClient(context.TODO(), GCP_PROJECT_NAME)
	pubsubPlugin, _ = pubsub_plugin.NewWithClient(pubsubClient)

	When("Listing Available Topics", func() {
		When("There are no topics available", func() {

			It("Should return an empty list of topics", func() {
				topics, err := pubsubPlugin.GetTopics()
				Expect(err).To(BeNil())
				Expect(topics).To(BeEmpty())
			})
		})

		When("There are topics available", func() {
			BeforeEach(func() {
				mockPubsubServer.SetTopics([]string{
					createPubsubTopicId("Test"),
				})
			})

			It("Should return all available topics", func() {
				topics, err := pubsubPlugin.GetTopics()
				Expect(err).To(BeNil())
				Expect(topics).To(ContainElement("Test"))
			})
		})
	})

	When("Publishing Messages", func() {
		event := &sdk.NitricEvent{
			RequestId:   "Test",
			PayloadType: "Test",
			Payload: map[string]interface{}{
				"Test": "Test",
			},
		}

		When("To a topic that does not exist", func() {
			BeforeEach(func() {
				mockPubsubServer.SetTopics([]string{})
			})

			It("should return an error", func() {
				err := pubsubPlugin.Publish("Test", event)
				Expect(err).ToNot(BeNil())
			})
		})

		When("To a topic that does exist", func() {
			pubsubTopicName := createPubsubTopicId("Test")
			BeforeEach(func() {
				mockPubsubServer.SetTopics([]string{
					pubsubTopicName,
				})
			})

			It("should successfully publish the message", func() {
				err := pubsubPlugin.Publish("Test", event)
				Expect(err).To(BeNil())
				Expect(mockPubsubServer.GetMessages()[pubsubTopicName]).To(HaveLen(1))
			})
		})
	})
})
