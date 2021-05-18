// Copyright 2021 Nitric Pty Ltd.
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

package firestore_service_test

import (
	"context"
	"errors"
	"net"
	"os"

	"cloud.google.com/go/firestore"
	mocks "github.com/nitric-dev/membrane/mocks/firestore"
	firestore_plugin "github.com/nitric-dev/membrane/plugins/kv/firestore"
	"github.com/nitric-dev/membrane/sdk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pb "google.golang.org/genproto/googleapis/firestore/v1"
	"google.golang.org/grpc"
)

var _ = Describe("Firestore KeyValue Plugin", func() {
	defer GinkgoRecover()

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

	key := map[string]interface{}{
		"key": "Test",
	}

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

				doc, err := firestorePlugin.Get("Test", key)

				Expect(err).To(BeNil())
				Expect(doc).To(BeEquivalentTo(item))
			})
		})

		When("And the document does not exist", func() {
			It("A not found error should be returned", func() {
				mockFirestoreServer.Store = map[string]map[string]map[string]*pb.Value{}

				_, err := firestorePlugin.Get("Test", key)

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

				err := firestorePlugin.Put("Test", key, map[string]interface{}{
					"Test": "Test2",
				})

				Expect(err).To(BeNil())
			})
		})
	})

	When("Delete", func() {
		When("the collection does not exist", func() {
			It("should return an error", func() {
				// key := map[string]interface{}{
				// 	"key": "Not Found",
				// }
				err := firestorePlugin.Delete("collection ?", key)
				Expect(err).NotTo(BeNil())
			})
		})
		When("the document does not exist", func() {
			It("should return not error", func() {
				mockFirestoreServer.Store = map[string]map[string]map[string]*pb.Value{
					// Collection Test
					"Collection": {
						// Resource Test
						"Key": {
							"Value": &pb.Value{
								ValueType: &pb.Value_StringValue{
									StringValue: "user@server.com",
								},
							},
						},
					},
				}
				key := map[string]interface{}{
					"key": "Not Found",
				}
				err := firestorePlugin.Delete("Collection", key)

				Expect(err).To(BeNil())
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

				err := firestorePlugin.Delete("Test", key)

				Expect(err).To(BeNil())
			})
		})
	})

	When("Query", func() {
		When("collection is empty", func() {
			It("should return an error", func() {
				result, err := firestorePlugin.Query("", nil, 0)

				Expect(err).ToNot(BeNil())
				Expect(result).To(BeNil())
			})
		})
		When("expressions is nil", func() {
			It("should return an error", func() {
				result, err := firestorePlugin.Query("collection", nil, 0)

				Expect(err).To(BeEquivalentTo(errors.New("provide non-nil expressions")))
				Expect(result).To(BeNil())
			})
		})
		// TODO: create a query mocking facility
		// When("empty result", func() {
		// 	It("should return empty list", func() {
		// 		exps := []sdk.QueryExpression{
		// 			{Operand: "Pk", Operator: "==", Value: "123"},
		// 		}
		// 		result, err := firestorePlugin.Query("collection", exps, 10)
		// 		Expect(result).NotTo(BeNil())
		// 		Expect(err).To(BeNil())
		// 	})
		// })
	})
})
