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

package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/golang/mock/gomock"
	"google.golang.org/protobuf/types/known/durationpb"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_azblob "github.com/nitrictech/nitric/cloud/azure/mocks/azblob"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
)

var _ = Describe("Azblob", func() {
	// Context("New", func() {
	//	When("", func() {

	//	})
	//})

	Context("Read", func() {
		When("Azure returns a successfully response", func() {
			crtl := gomock.NewController(GinkgoT())
			mockAzblob := mock_azblob.NewMockAzblobServiceUrlIface(crtl)
			mockContainer := mock_azblob.NewMockAzblobContainerUrlIface(crtl)
			mockBlob := mock_azblob.NewMockAzblobBlockBlobUrlIface(crtl)
			mockDown := mock_azblob.NewMockAzblobDownloadResponse(crtl)

			storagePlugin := &AzblobStorageService{
				client: mockAzblob,
			}

			It("should successfully return the read payload", func() {
				By("Retrieving the Container URL for the requested bucket")
				mockAzblob.EXPECT().NewContainerURL("my-bucket").Times(1).Return(mockContainer)

				By("Retrieving the blob url of the requested object")
				mockContainer.EXPECT().NewBlockBlobURL("my-blob").Times(1).Return(mockBlob)

				By("Calling Download once on the blob with the expected options")
				mockBlob.EXPECT().Download(
					gomock.Any(),
					int64(0),
					int64(0),
					azblob.BlobAccessConditions{},
					false,
					azblob.ClientProvidedKeyOptions{},
				).Times(1).Return(mockDown, nil)

				By("Reading from the download response")
				mockDown.EXPECT().Body(gomock.Any()).Times(1).Return(io.NopCloser(strings.NewReader("file-contents")))

				data, err := storagePlugin.Read(context.TODO(), &storagepb.StorageReadRequest{
					BucketName: "my-bucket",
					Key:        "my-blob",
				})

				By("Not returning an error")
				Expect(err).ToNot(HaveOccurred())

				By("Returning the read data")
				Expect(data.Body).To(BeEquivalentTo([]byte("file-contents")))

				crtl.Finish()
			})
		})

		When("Azure returns an error", func() {
			crtl := gomock.NewController(GinkgoT())
			mockAzblob := mock_azblob.NewMockAzblobServiceUrlIface(crtl)
			mockContainer := mock_azblob.NewMockAzblobContainerUrlIface(crtl)
			mockBlob := mock_azblob.NewMockAzblobBlockBlobUrlIface(crtl)

			storagePlugin := &AzblobStorageService{
				client: mockAzblob,
			}

			It("should return an error", func() {
				By("Retrieving the Container URL for the requested bucket")
				mockAzblob.EXPECT().NewContainerURL("my-bucket").Times(1).Return(mockContainer)

				By("Retrieving the blob url of the requested object")
				mockContainer.EXPECT().NewBlockBlobURL("my-blob").Times(1).Return(mockBlob)

				By("Calling Download once on the blob with the expected options")
				mockBlob.EXPECT().Download(
					gomock.Any(),
					int64(0),
					int64(0),
					azblob.BlobAccessConditions{},
					false,
					azblob.ClientProvidedKeyOptions{},
				).Times(1).Return(nil, fmt.Errorf("Failed to download"))

				_, err := storagePlugin.Read(context.TODO(), &storagepb.StorageReadRequest{
					BucketName: "my-bucket",
					Key:        "my-blob",
				})

				By("Returning an error")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("Write", func() {
		When("Azure returns a successful response", func() {
			crtl := gomock.NewController(GinkgoT())
			mockAzblob := mock_azblob.NewMockAzblobServiceUrlIface(crtl)
			mockContainer := mock_azblob.NewMockAzblobContainerUrlIface(crtl)
			mockBlob := mock_azblob.NewMockAzblobBlockBlobUrlIface(crtl)

			storagePlugin := &AzblobStorageService{
				client: mockAzblob,
			}

			It("should successfully write the blob", func() {
				By("Retrieving the Container URL for the requested bucket")
				mockAzblob.EXPECT().NewContainerURL("my-bucket").Times(1).Return(mockContainer)

				By("Retrieving the blob url of the requested object")
				mockContainer.EXPECT().NewBlockBlobURL("my-blob").Times(1).Return(mockBlob)

				By("Calling Upload once on the blob with the expected options")
				mockBlob.EXPECT().Upload(
					gomock.Any(),
					bytes.NewReader([]byte("test")),
					azblob.BlobHTTPHeaders{},
					azblob.Metadata{},
					azblob.BlobAccessConditions{},
					azblob.DefaultAccessTier,
					nil,
					azblob.ClientProvidedKeyOptions{},
				).Times(1).Return(&azblob.BlockBlobUploadResponse{}, nil)

				_, err := storagePlugin.Write(context.TODO(), &storagepb.StorageWriteRequest{
					BucketName: "my-bucket",
					Key:        "my-blob",
					Body:       []byte("test"),
				})

				By("Not returning an error")
				Expect(err).ToNot(HaveOccurred())

				crtl.Finish()
			})
		})

		When("Azure returns an error", func() {
			crtl := gomock.NewController(GinkgoT())
			mockAzblob := mock_azblob.NewMockAzblobServiceUrlIface(crtl)
			mockContainer := mock_azblob.NewMockAzblobContainerUrlIface(crtl)
			mockBlob := mock_azblob.NewMockAzblobBlockBlobUrlIface(crtl)

			storagePlugin := &AzblobStorageService{
				client: mockAzblob,
			}

			It("should return an error", func() {
				By("Retrieving the Container URL for the requested bucket")
				mockAzblob.EXPECT().NewContainerURL("my-bucket").Times(1).Return(mockContainer)

				By("Retrieving the blob url of the requested object")
				mockContainer.EXPECT().NewBlockBlobURL("my-blob").Times(1).Return(mockBlob)

				By("Calling Upload once on the blob with the expected options")
				mockBlob.EXPECT().Upload(
					gomock.Any(),
					bytes.NewReader([]byte("test")),
					azblob.BlobHTTPHeaders{},
					azblob.Metadata{},
					azblob.BlobAccessConditions{},
					azblob.DefaultAccessTier,
					nil,
					azblob.ClientProvidedKeyOptions{},
				).Times(1).Return(nil, fmt.Errorf("mock-error"))

				_, err := storagePlugin.Write(context.TODO(), &storagepb.StorageWriteRequest{
					BucketName: "my-bucket",
					Key:        "my-blob",
					Body:       []byte("test"),
				})

				By("returning an error")
				Expect(err).To(HaveOccurred())

				crtl.Finish()
			})
		})
	})

	Context("Delete", func() {
		When("Azure returns a successful response", func() {
			crtl := gomock.NewController(GinkgoT())
			mockAzblob := mock_azblob.NewMockAzblobServiceUrlIface(crtl)
			mockContainer := mock_azblob.NewMockAzblobContainerUrlIface(crtl)
			mockBlob := mock_azblob.NewMockAzblobBlockBlobUrlIface(crtl)

			storagePlugin := &AzblobStorageService{
				client: mockAzblob,
			}

			It("should successfully write the blob", func() {
				By("Retrieving the Container URL for the requested bucket")
				mockAzblob.EXPECT().NewContainerURL("my-bucket").Times(1).Return(mockContainer)

				By("Retrieving the blob url of the requested object")
				mockContainer.EXPECT().NewBlockBlobURL("my-blob").Times(1).Return(mockBlob)

				By("Calling Upload once on the blob with the expected options")
				mockBlob.EXPECT().Delete(
					gomock.Any(),
					azblob.DeleteSnapshotsOptionInclude,
					azblob.BlobAccessConditions{},
				).Times(1).Return(&azblob.BlobDeleteResponse{}, nil)

				_, err := storagePlugin.Delete(context.TODO(), &storagepb.StorageDeleteRequest{
					BucketName: "my-bucket",
					Key:        "my-blob",
				})

				By("Not returning an error")
				Expect(err).ToNot(HaveOccurred())

				crtl.Finish()
			})
		})

		When("Azure returns an error", func() {
			crtl := gomock.NewController(GinkgoT())
			mockAzblob := mock_azblob.NewMockAzblobServiceUrlIface(crtl)
			mockContainer := mock_azblob.NewMockAzblobContainerUrlIface(crtl)
			mockBlob := mock_azblob.NewMockAzblobBlockBlobUrlIface(crtl)

			storagePlugin := &AzblobStorageService{
				client: mockAzblob,
			}

			It("should successfully write the blob", func() {
				By("Retrieving the Container URL for the requested bucket")
				mockAzblob.EXPECT().NewContainerURL("my-bucket").Times(1).Return(mockContainer)

				By("Retrieving the blob url of the requested object")
				mockContainer.EXPECT().NewBlockBlobURL("my-blob").Times(1).Return(mockBlob)

				By("Calling Upload once on the blob with the expected options")
				mockBlob.EXPECT().Delete(
					gomock.Any(),
					azblob.DeleteSnapshotsOptionInclude,
					azblob.BlobAccessConditions{},
				).Times(1).Return(nil, fmt.Errorf("mock-error"))

				_, err := storagePlugin.Delete(context.TODO(), &storagepb.StorageDeleteRequest{
					BucketName: "my-bucket",
					Key:        "my-blob",
				})

				By("Not returning an error")
				Expect(err).To(HaveOccurred())

				crtl.Finish()
			})
		})
	})

	Context("ListFiles", func() {
		When("Azure returns a successful response", func() {
			ctrl := gomock.NewController(GinkgoT())
			doneMarker := ""
			mockAzblob := mock_azblob.NewMockAzblobServiceUrlIface(ctrl)
			mockContainer := mock_azblob.NewMockAzblobContainerUrlIface(ctrl)

			storagePlugin := &AzblobStorageService{
				client: mockAzblob,
			}

			It("should should return a list of files on the bucket", func() {
				By("Retrieving the Container URL for the requested bucket")
				mockAzblob.EXPECT().NewContainerURL("my-bucket").Times(1).Return(mockContainer)

				By("The container returning list of files")
				mockContainer.EXPECT().ListBlobsFlatSegment(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&azblob.ListBlobsFlatSegmentResponse{
					NextMarker: azblob.Marker{
						Val: &doneMarker,
					},
					Segment: azblob.BlobFlatListSegment{
						BlobItems: []azblob.BlobItemInternal{
							{
								Name: "/test/test.png",
							},
						},
					},
				}, nil)

				resp, err := storagePlugin.ListBlobs(context.TODO(), &storagepb.StorageListBlobsRequest{
					BucketName: "my-bucket",
					Prefix:     "/test/",
				})

				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("Returning a single file")
				Expect(resp.Blobs).To(HaveLen(1))

				By("Having the returned key")
				Expect(resp.Blobs[0].Key).To(Equal("/test/test.png"))

				ctrl.Finish()
			})
		})

		When("Azure returns an error", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockAzblob := mock_azblob.NewMockAzblobServiceUrlIface(ctrl)
			mockContainer := mock_azblob.NewMockAzblobContainerUrlIface(ctrl)

			storagePlugin := &AzblobStorageService{
				client: mockAzblob,
			}

			It("should return a wrapped error", func() {
				By("Retrieving the Container URL for the requested bucket")
				mockAzblob.EXPECT().NewContainerURL("my-bucket").Times(1).Return(mockContainer)

				By("Azure returning an error")
				mockContainer.EXPECT().ListBlobsFlatSegment(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil, fmt.Errorf("mock-error"))

				files, err := storagePlugin.ListBlobs(context.TODO(), &storagepb.StorageListBlobsRequest{
					BucketName: "my-bucket",
				})

				By("returning nil results")
				Expect(files).To(BeNil())

				By("returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Context("PresignUrl", func() {
		When("User delegation credentials are accessible", func() {
			crtl := gomock.NewController(GinkgoT())
			mockAzblob := mock_azblob.NewMockAzblobServiceUrlIface(crtl)
			mockContainer := mock_azblob.NewMockAzblobContainerUrlIface(crtl)
			mockBlob := mock_azblob.NewMockAzblobBlockBlobUrlIface(crtl)

			storagePlugin := &AzblobStorageService{
				client: mockAzblob,
			}

			It("should return a presigned url", func() {
				By("Retrieving the Container URL for the requested bucket")
				mockAzblob.EXPECT().NewContainerURL("my-bucket").Times(1).Return(mockContainer)

				By("Retrieving the blob url of the requested object")
				mockContainer.EXPECT().NewBlockBlobURL("my-blob").Times(1).Return(mockBlob)

				By("Retrieving user delegation credentials")
				mockAzblob.EXPECT().GetUserDelegationCredential(
					context.TODO(), gomock.Any(), gomock.Any(), nil,
				).Return(
					azblob.NewUserDelegationCredential("mock-account-name", azblob.UserDelegationKey{}),
					nil,
				)

				u, _ := url.Parse("https://fake-account.com/my-bucket/my-blob")
				By("Getting the URL")
				mockBlob.EXPECT().Url().Return(*u)

				resp, err := storagePlugin.PreSignUrl(context.TODO(), &storagepb.StoragePreSignUrlRequest{
					BucketName: "my-bucket",
					Key:        "my-blob",
					Operation:  storagepb.StoragePreSignUrlRequest_READ,
					Expiry:     durationpb.New(time.Second * 3600),
				})

				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("Returning a pre-signed URL from the computed blob URL")
				Expect(resp.Url).To(ContainSubstring("https://fake-account.com/my-bucket/my-blob"))
			})
		})

		When("retrieving user delegation credentials fails", func() {
			crtl := gomock.NewController(GinkgoT())
			mockAzblob := mock_azblob.NewMockAzblobServiceUrlIface(crtl)
			mockContainer := mock_azblob.NewMockAzblobContainerUrlIface(crtl)
			mockBlob := mock_azblob.NewMockAzblobBlockBlobUrlIface(crtl)

			storagePlugin := &AzblobStorageService{
				client: mockAzblob,
			}

			It("should return an error", func() {
				By("Retrieving the Container URL for the requested bucket")
				mockAzblob.EXPECT().NewContainerURL("my-bucket").Times(1).Return(mockContainer)

				By("Retrieving the blob url of the requested object")
				mockContainer.EXPECT().NewBlockBlobURL("my-blob").Times(1).Return(mockBlob)

				By("Failing to retrieve user delegation credentials")
				mockAzblob.EXPECT().GetUserDelegationCredential(
					context.TODO(), gomock.Any(), gomock.Any(), nil,
				).Return(
					nil,
					fmt.Errorf("mock-error"),
				)

				u, _ := url.Parse("https://fake-account.com/my-bucket/my-blob")
				By("Getting the URL")
				mockBlob.EXPECT().Url().Return(*u)

				resp, err := storagePlugin.PreSignUrl(context.TODO(), &storagepb.StoragePreSignUrlRequest{
					BucketName: "my-bucket",
					Key:        "my-blob",
					Operation:  storagepb.StoragePreSignUrlRequest_READ,
					Expiry:     durationpb.New(time.Second * 3600),
				})

				By("Returning nil")
				Expect(resp).To(BeNil())

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Context("Exists", func() {
		When("Azure returns a successful response", func() {
			Context("and the file exists", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockAzblob := mock_azblob.NewMockAzblobServiceUrlIface(ctrl)
				mockContainer := mock_azblob.NewMockAzblobContainerUrlIface(ctrl)
				mockBlob := mock_azblob.NewMockAzblobBlockBlobUrlIface(ctrl)

				storagePlugin := &AzblobStorageService{
					client: mockAzblob,
				}

				It("should should return true", func() {
					By("Retrieving the Container URL for the requested bucket")
					mockAzblob.EXPECT().NewContainerURL("my-bucket").Times(1).Return(mockContainer)

					By("Retrieving the Blob URL for the requested file")
					mockContainer.EXPECT().NewBlockBlobURL("test-file").Times(1).Return(mockBlob)

					By("Calling GetProperties on the file")
					mockBlob.EXPECT().GetProperties(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&azblob.BlobGetPropertiesResponse{}, nil)

					By("The file returning that it exists")
					resp, err := storagePlugin.Exists(context.TODO(), &storagepb.StorageExistsRequest{
						BucketName: "my-bucket",
						Key:        "test-file",
					})

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Returning true")
					Expect(resp.Exists).To(BeTrue())

					ctrl.Finish()
				})
			})

			Context("and the file does not exist", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockAzblob := mock_azblob.NewMockAzblobServiceUrlIface(ctrl)
				mockContainer := mock_azblob.NewMockAzblobContainerUrlIface(ctrl)
				mockBlob := mock_azblob.NewMockAzblobBlockBlobUrlIface(ctrl)
				mockError := mock_azblob.NewMockStorageError(ctrl)

				storagePlugin := &AzblobStorageService{
					client: mockAzblob,
				}

				It("should should return false", func() {
					By("Retrieving the Container URL for the requested bucket")
					mockAzblob.EXPECT().NewContainerURL("my-bucket").Times(1).Return(mockContainer)

					By("Retrieving the Blob URL for the requested file")
					mockContainer.EXPECT().NewBlockBlobURL("test-file").Times(1).Return(mockBlob)

					By("Producing a service code of azblob.ServiceCodeBlobNotFound")
					mockError.EXPECT().ServiceCode().Times(1).Return(azblob.ServiceCodeBlobNotFound)

					By("Calling GetProperties on the file")
					mockBlob.EXPECT().GetProperties(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil, mockError)

					By("The file returning that it exists")
					resp, err := storagePlugin.Exists(context.TODO(), &storagepb.StorageExistsRequest{
						BucketName: "my-bucket",
						Key:        "test-file",
					})

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Returning false")
					Expect(resp.Exists).To(BeFalse())

					ctrl.Finish()
				})
			})
		})

		When("Azure returns an error", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockAzblob := mock_azblob.NewMockAzblobServiceUrlIface(ctrl)
			mockContainer := mock_azblob.NewMockAzblobContainerUrlIface(ctrl)
			mockBlob := mock_azblob.NewMockAzblobBlockBlobUrlIface(ctrl)
			mockError := mock_azblob.NewMockStorageError(ctrl)

			storagePlugin := &AzblobStorageService{
				client: mockAzblob,
			}

			It("should should return an error", func() {
				By("Retrieving the Container URL for the requested bucket")
				mockAzblob.EXPECT().NewContainerURL("my-bucket").Times(1).Return(mockContainer)

				By("Retrieving the Blob URL for the requested file")
				mockContainer.EXPECT().NewBlockBlobURL("test-file").Times(1).Return(mockBlob)

				By("Producing a service code of azblob.ServiceCodeInternalError")
				mockError.EXPECT().ServiceCode().Times(1).Return(azblob.ServiceCodeInternalError)

				By("Inspecting the service error")
				mockError.EXPECT().Error().AnyTimes()

				By("Calling GetProperties on the file")
				mockBlob.EXPECT().GetProperties(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil, mockError)

				By("The file returning that it exists")
				resp, err := storagePlugin.Exists(context.TODO(), &storagepb.StorageExistsRequest{
					BucketName: "my-bucket",
					Key:        "test-file",
				})

				By("Not returning an error")
				Expect(err).Should(HaveOccurred())

				By("Returning nil")
				Expect(resp).To(BeNil())

				ctrl.Finish()
			})
		})
	})
})
