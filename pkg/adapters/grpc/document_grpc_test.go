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

package grpc_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/golang/mock/gomock"

	mock_document "github.com/nitrictech/nitric/mocks/document"
	"github.com/nitrictech/nitric/pkg/adapters/grpc"
	v1 "github.com/nitrictech/nitric/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/pkg/plugins/document"
	"github.com/nitrictech/protoutils"
)

var _ = Describe("GRPC Document", func() {
	Context("Get", func() {
		When("plugin not registered", func() {
			dss := &grpc.DocumentServiceServer{}
			resp, err := dss.Get(context.Background(), &v1.DocumentGetRequest{})
			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("Document plugin not registered"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request not valid", func() {
			g := gomock.NewController(GinkgoT())
			mockDS := mock_document.NewMockDocumentService(g)
			dss := grpc.NewDocumentServer(mockDS)
			resp, err := dss.Get(context.Background(), &v1.DocumentGetRequest{})

			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("invalid DocumentGetRequest.Key: value is required"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request is valid", func() {
			g := gomock.NewController(GinkgoT())

			mockDS := mock_document.NewMockDocumentService(g)
			key := &document.Key{
				Collection: &document.Collection{Name: "test"},
				Id:         "123456",
			}
			doc := &document.Document{
				Key: key,
				Content: map[string]interface{}{
					"x": "y",
				},
			}
			expect := &v1.Document{
				Key: &v1.Key{
					Collection: &v1.Collection{
						Name: key.Collection.Name,
					},
					Id: key.Id,
				},
			}
			var err error
			expect.Content, err = protoutils.NewStruct(doc.Content)
			Expect(err).Should(BeNil())

			mockDS.EXPECT().Get(context.TODO(), key).Return(doc, nil)

			dss := grpc.NewDocumentServer(mockDS)
			resp, err := dss.Get(context.Background(), &v1.DocumentGetRequest{Key: expect.Key})

			It("Should return a doc", func() {
				Expect(err).Should(BeNil())
				Expect(resp.Document.Content).Should(Equal(expect.Content))
			})
		})
	})

	Context("Query", func() {
		When("plugin not registered", func() {
			dss := &grpc.DocumentServiceServer{}
			resp, err := dss.Query(context.Background(), &v1.DocumentQueryRequest{})
			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("Document plugin not registered"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request not valid", func() {
			g := gomock.NewController(GinkgoT())
			mockDS := mock_document.NewMockDocumentService(g)
			dss := grpc.NewDocumentServer(mockDS)
			resp, err := dss.Query(context.Background(), &v1.DocumentQueryRequest{})

			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("invalid DocumentQueryRequest.Collection: value is required"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request is valid", func() {
			g := gomock.NewController(GinkgoT())

			mockDS := mock_document.NewMockDocumentService(g)
			key := &document.Key{
				Collection: &document.Collection{Name: "test"},
				Id:         "123456",
			}
			doc := &document.Document{
				Key: key,
				Content: map[string]interface{}{
					"x": "y",
				},
			}

			mockDS.EXPECT().Query(gomock.Any(), &document.Collection{Name: "zed"}, []document.QueryExpression{
				{
					Operand:  "count",
					Operator: ">",
					Value:    int64(5),
				},
			}, 3, map[string]string{}).Return(&document.QueryResult{
				Documents:   []document.Document{*doc},
				PagingToken: map[string]string{},
			}, nil)

			dss := grpc.NewDocumentServer(mockDS)
			resp, err := dss.Query(context.Background(), &v1.DocumentQueryRequest{
				Collection: &v1.Collection{
					Name: "zed",
				},
				Expressions: []*v1.Expression{
					{
						Operand:  "count",
						Operator: ">",
						Value:    &v1.ExpressionValue{Kind: &v1.ExpressionValue_IntValue{IntValue: int64(5)}},
					},
				},
				Limit:       3,
				PagingToken: map[string]string{},
			})

			It("Should return a doc", func() {
				Expect(err).Should(BeNil())
				Expect(resp.Documents[0].Key.Id).Should(Equal("123456"))
			})
		})
	})

	Context("Set", func() {
		When("plugin not registered", func() {
			dss := &grpc.DocumentServiceServer{}
			resp, err := dss.Set(context.Background(), &v1.DocumentSetRequest{})
			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("Document plugin not registered"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request not valid", func() {
			g := gomock.NewController(GinkgoT())
			mockDS := mock_document.NewMockDocumentService(g)
			dss := grpc.NewDocumentServer(mockDS)
			resp, err := dss.Set(context.Background(), &v1.DocumentSetRequest{})

			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("invalid DocumentSetRequest.Key: value is required"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request is valid", func() {
			g := gomock.NewController(GinkgoT())
			mockDS := mock_document.NewMockDocumentService(g)
			key := &document.Key{
				Collection: &document.Collection{Name: "test"},
				Id:         "123456",
			}
			doc := &document.Document{
				Key: key,
				Content: map[string]interface{}{
					"x": "y",
				},
			}
			expect := &v1.Document{
				Key: &v1.Key{
					Collection: &v1.Collection{
						Name: key.Collection.Name,
					},
					Id: key.Id,
				},
			}
			var err error
			expect.Content, err = protoutils.NewStruct(doc.Content)
			Expect(err).Should(BeNil())

			mockDS.EXPECT().Set(gomock.Any(), key, expect.Content.AsMap()).Return(nil)

			dss := grpc.NewDocumentServer(mockDS)
			resp, err := dss.Set(context.Background(), &v1.DocumentSetRequest{
				Key:     expect.Key,
				Content: expect.Content,
			})

			It("Should return a doc", func() {
				Expect(err).Should(BeNil())
				Expect(resp.String()).Should(Equal(""))
			})
		})
	})

	Context("Delete", func() {
		When("plugin not registered", func() {
			dss := &grpc.DocumentServiceServer{}
			resp, err := dss.Delete(context.Background(), &v1.DocumentDeleteRequest{})
			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("Document plugin not registered"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request not valid", func() {
			g := gomock.NewController(GinkgoT())
			mockDS := mock_document.NewMockDocumentService(g)
			dss := grpc.NewDocumentServer(mockDS)
			resp, err := dss.Delete(context.Background(), &v1.DocumentDeleteRequest{})

			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("invalid DocumentDeleteRequest.Key: value is required"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request is valid", func() {
			g := gomock.NewController(GinkgoT())
			mockDS := mock_document.NewMockDocumentService(g)
			key := &document.Key{
				Collection: &document.Collection{Name: "test"},
				Id:         "123456",
			}
			doc := &document.Document{
				Key: key,
				Content: map[string]interface{}{
					"x": "y",
				},
			}
			expect := &v1.Document{
				Key: &v1.Key{
					Collection: &v1.Collection{
						Name: key.Collection.Name,
					},
					Id: key.Id,
				},
			}
			var err error
			expect.Content, err = protoutils.NewStruct(doc.Content)
			Expect(err).Should(BeNil())

			mockDS.EXPECT().Delete(gomock.Any(), key).Return(nil)

			dss := grpc.NewDocumentServer(mockDS)
			resp, err := dss.Delete(context.Background(), &v1.DocumentDeleteRequest{
				Key: expect.Key,
			})

			It("Should delete a doc", func() {
				Expect(err).Should(BeNil())
				Expect(resp.String()).Should(Equal(""))
			})
		})
	})
})
