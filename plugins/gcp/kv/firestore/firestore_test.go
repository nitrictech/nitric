package firestore_service_test

import (
	"context"
	"net"
	"os"

	"cloud.google.com/go/firestore"
	firestore_plugin "github.com/nitric-dev/membrane/plugins/gcp/kv/firestore"
	mocks "github.com/nitric-dev/membrane/plugins/gcp/mocks"
	"github.com/nitric-dev/membrane/plugins/sdk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pb "google.golang.org/genproto/googleapis/firestore/v1"
	"google.golang.org/grpc"
)

var _ = Describe("Firestore KeyValue Plugin", func() {
	// Setup mock environment...
	var opts []grpc.ServerOption
	var firestoreClient *firestore.Client
	var firestorePlugin sdk.KeyValueService
	grpcServer := grpc.NewServer(opts...)
	mockFirestoreServer := &mocks.MockFirestoreServer{
		Store: make(map[string]map[string]map[string]*pb.Value),
	}
	// Set the emulator host...
	os.Setenv("FIRESTORE_EMULATOR_HOST", "127.0.0.1:50051")
	pb.RegisterFirestoreServer(grpcServer, mockFirestoreServer)
	lis, _ := net.Listen("tcp", "127.0.0.1:50051")
	// Do not block on serve...
	go (func() {
		grpcServer.Serve(lis)
	})()

	// clientConn, _ := grpc.Dial("127.0.0.1:50051")
	firestoreClient, _ = firestore.NewClient(context.TODO(), "")
	firestorePlugin, _ = firestore_plugin.NewWithClient(firestoreClient)

	AfterSuite(func() {
		grpcServer.GracefulStop()
	})

	AfterEach(func() {
		mockFirestoreServer.ClearStore()
	})

	When("Get", func() {
		When("And the document already exists", func() {

			It("The stored document should be returned", func() {
				item := map[string]interface{}{
					"Test": "Test",
				}
				mockFirestoreServer.Store = map[string]map[string]map[string]*pb.Value{
					// Collection Test
					"Test": {
						// Resource Test
						"Test": {
							"Test": &pb.Value{
								ValueType: &pb.Value_StringValue{
									StringValue: "Test",
								},
							},
						},
					},
				}

				doc, err := firestorePlugin.Get("Test", "Test")

				Expect(err).To(BeNil())
				Expect(doc).To(BeEquivalentTo(item))
			})
		})

		When("And the document does not exist", func() {
			It("A not found error should be returned", func() {
				mockFirestoreServer.Store = map[string]map[string]map[string]*pb.Value{}

				_, err := firestorePlugin.Get("Test", "Test")

				Expect(err).ToNot(BeNil())
			})
		})
	})

	When("Put", func() {
		When("the document already exists", func() {
			It("should successfully update the document", func() {
				mockFirestoreServer.Store = map[string]map[string]map[string]*pb.Value{
					// Collection Test
					"Test": {
						// Resource Test
						"Test": {
							"Test": &pb.Value{
								ValueType: &pb.Value_StringValue{
									StringValue: "Test",
								},
							},
						},
					},
				}

				err := firestorePlugin.Put("Test", "Test", map[string]interface{}{
					"Test": "Test2",
				})

				Expect(err).To(BeNil())
			})
		})
	})

	When("Delete", func() {
		When("the document does not exist", func() {
			It("should return an error", func() {
				err := firestorePlugin.Delete("Test", "Test")

				Expect(err).ToNot(BeNil())
			})
		})

		When("the document exists", func() {
			It("should successfully delete the document", func() {
				mockFirestoreServer.Store = map[string]map[string]map[string]*pb.Value{
					// Collection Test
					"Test": {
						// Resource Test
						"Test": {
							"Test": &pb.Value{
								ValueType: &pb.Value_StringValue{
									StringValue: "Test",
								},
							},
						},
					},
				}

				err := firestorePlugin.Delete("Test", "Test")

				Expect(err).To(BeNil())
			})
		})
	})
})
