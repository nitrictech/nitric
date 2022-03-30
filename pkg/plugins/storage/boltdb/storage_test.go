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

package boltdb_storage_service_test

import (
	"os"

	boltdb_storage_service "github.com/nitrictech/nitric/pkg/plugins/storage/boltdb"
	"github.com/nitrictech/nitric/pkg/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const BUCKET = "bucket"
const KEY = "key"
const DATA = "data"

var local_storage_directory = utils.GetRelativeDevPath(boltdb_storage_service.DEV_SUB_DIRECTORY)

var _ = Describe("Storage", func() {
	AfterSuite(func() {
		// Cleanup default secret directory
		os.RemoveAll(utils.GetDevVolumePath())
	})

	storagePlugin, err := boltdb_storage_service.New()
	if err != nil {
		panic(err)
	}

	AfterEach(func() {
		err := os.RemoveAll(local_storage_directory)
		if err != nil {
			panic(err)
		}

		_, err = os.Stat(local_storage_directory)
		if os.IsNotExist(err) {
			// Make directory if not present
			err := os.Mkdir(local_storage_directory, 0777)
			if err != nil {
				panic(err)
			}
		}
	})

	Context("Write", func() {
		Context("When object is nil", func() {
			It("Should return an error", func() {
				err := storagePlugin.Write(BUCKET, KEY, nil)
				Expect(err).Should(HaveOccurred())
			})
		})

		Context("When object is empty", func() {
			It("Should return an error", func() {
				err := storagePlugin.Write(BUCKET, KEY, []byte{})
				Expect(err).Should(HaveOccurred())
			})
		})

		Context("Valid write operation", func() {
			It("Should store the provided byte array", func() {
				err := storagePlugin.Write(BUCKET, KEY, []byte(DATA))
				Expect(err).To(BeNil())

				data, err := storagePlugin.Read(BUCKET, KEY)
				Expect(err).To(BeNil())
				Expect(data).NotTo(BeNil())
				Expect(data).To(BeEquivalentTo([]byte(DATA)))
			})
		})
	})

	Context("Read", func() {
		Context("Valid read operation", func() {
			It("Should read the provided byte array", func() {
				err := storagePlugin.Write(BUCKET, KEY, []byte(DATA))
				Expect(err).To(BeNil())

				data, err := storagePlugin.Read(BUCKET, KEY)
				Expect(err).To(BeNil())
				Expect(data).NotTo(BeNil())
				Expect(data).To(BeEquivalentTo([]byte(DATA)))
			})
		})

		Context("Read missing object operation", func() {
			It("Should return an error", func() {
				data, err := storagePlugin.Read(BUCKET, "not-found")
				Expect(data).To(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Context("Delete", func() {
		Context("Valid delete operation", func() {
			It("Should read the provided byte array", func() {
				err := storagePlugin.Write(BUCKET, KEY, []byte(DATA))
				Expect(err).To(BeNil())

				err = storagePlugin.Delete(BUCKET, KEY)
				Expect(err).To(BeNil())
			})
		})

		Context("Delete missing object operation", func() {
			It("Should return an error", func() {
				err := storagePlugin.Delete(BUCKET, "not-found")
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
