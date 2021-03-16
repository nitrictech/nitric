package firestore_service_test

import (
	"context"
	"net"
	"os"

	"cloud.google.com/go/firestore"
	firestore_plugin "github.com/nitric-dev/membrane/plugins/gcp/documents/firestore"
	mocks "github.com/nitric-dev/membrane/plugins/gcp/mocks"
	"github.com/nitric-dev/membrane/plugins/sdk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pb "google.golang.org/genproto/googleapis/firestore/v1"
	"google.golang.org/grpc"
)

var _ = Describe("Firestore Documents Plugin", func() {
	// Setup mock environment...
	var opts []grpc.ServerOption
	var firestoreClient *firestore.Client
	var firestorePlugin sdk.DocumentService
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

	When("Creating a new document", func() {
		When("And the document does not already exist", func() {
			err := firestorePlugin.Create("Test", "Test", map[string]interface{}{
				"Test": "Test",
			})
			It("Should create and store the document", func() {
				Expect(err).To(BeNil())
			})
		})

		When("and the document already exists", func() {
			It("Should return an AlreadyExists error", func() {
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

				err := firestorePlugin.Create("Test", "Test", map[string]interface{}{
					"Test": "Test",
				})

				Expect(err).ToNot(BeNil())
			})
		})
	})

	When("Retrieving a document", func() {
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

	When("updating a document", func() {
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

				err := firestorePlugin.Update("Test", "Test", map[string]interface{}{
					"Test": "Test2",
				})

				Expect(err).To(BeNil())
			})
		})

		When("the document doesn't exist", func() {
			It("should return a not found error", func() {
				err := firestorePlugin.Update("Test", "Test", map[string]interface{}{
					"Test": "Test",
				})

				Expect(err).ToNot(BeNil())
			})
		})
	})

	When("deleting a document", func() {
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
