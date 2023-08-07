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

package storage_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/nitrictech/nitric/core/pkg/plugins/storage"
)

var _ = Describe("Unimplemented Storage Plugin Tests", func() {
	uisp := &storage.UnimplementedStoragePlugin{}

	Context("Read", func() {
		When("Calling Read on UnimplementedStoragePlugin", func() {
			_, err := uisp.Read(context.TODO(), "test", "test")

			It("should return an unimplemented error", func() {
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("UNIMPLEMENTED"))
			})
		})
	})

	Context("Write", func() {
		When("Calling Write on UnimplementedStoragePlugin", func() {
			err := uisp.Write(context.TODO(), "test", "test", nil)

			It("should return an unimplemented error", func() {
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("UNIMPLEMENTED"))
			})
		})
	})

	Context("Delete", func() {
		When("Calling Delete on UnimplementedStoragePlugin", func() {
			err := uisp.Delete(context.TODO(), "test", "test")

			It("should return an unimplemented error", func() {
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("UNIMPLEMENTED"))
			})
		})
	})

	Context("ListFiles", func() {
		When("Calling ListFiles on UnimplementedStoragePlugin", func() {
			_, err := uisp.ListFiles(context.TODO(), "test", nil)

			It("should return an unimplemented error", func() {
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("UNIMPLEMENTED"))
			})
		})
	})

	Context("PreSignUrl", func() {
		When("Calling PreSignUrl on UnimplementedStoragePlugin", func() {
			_, err := uisp.PreSignUrl(context.TODO(), "test", "test", storage.READ, 300)

			It("should return an unimplemented error", func() {
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("UNIMPLEMENTED"))
			})
		})
	})

	Context("Exists", func() {
		When("Calling Exists on UnimplementedStoragePlugin", func() {
			_, err := uisp.Exists(context.TODO(), "test", "test")

			It("should return an unimplemented error", func() {
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("UNIMPLEMENTED"))
			})
		})
	})
})
