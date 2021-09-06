package azblob_service

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_azblob "github.com/nitric-dev/membrane/mocks/azblob"
)

var _ = Describe("Azblob", func() {
	//Context("New", func() {
	//	When("", func() {

	//	})
	//})

	Context("Read", func() {
		When("Azure returns a successfuly response", func() {
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
				mockDown.EXPECT().Body(gomock.Any()).Times(1).Return(ioutil.NopCloser(strings.NewReader("file-contents")))

				data, err := storagePlugin.Read("my-bucket", "my-blob")

				By("Not returning an error")
				Expect(err).ToNot(HaveOccurred())

				By("Returning the read data")
				Expect(data).To(BeEquivalentTo([]byte("file-contents")))

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

				_, err := storagePlugin.Read("my-bucket", "my-blob")

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

				err := storagePlugin.Write("my-bucket", "my-blob", []byte("test"))

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

				err := storagePlugin.Write("my-bucket", "my-blob", []byte("test"))

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

				err := storagePlugin.Delete("my-bucket", "my-blob")

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

				err := storagePlugin.Delete("my-bucket", "my-blob")

				By("Not returning an error")
				Expect(err).To(HaveOccurred())

				crtl.Finish()
			})
		})
	})
})
