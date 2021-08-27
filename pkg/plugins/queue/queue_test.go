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
package queue_test

import (
	"github.com/nitric-dev/membrane/pkg/plugins/queue"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Queue Plugin", func() {

	When("NitricTask.String", func() {
		It("should print NitricTask", func() {
			task := &queue.NitricTask{
				ID:          "id",
				LeaseID:     "leaseId",
				PayloadType: "payloadType",
				Payload: map[string]interface{}{
					"key": "value",
				},
			}
			Expect(task.String()).To(BeEquivalentTo("{ID: id, LeaseID: leaseId, PayloadType: payloadType}"))
		})
	})

	When("ReceiveOptions.String", func() {
		It("should print ReceiveOptions", func() {
			// depth := uint32(10)
			options := &queue.ReceiveOptions{
				QueueName: "queue",
			}
			Expect(options.String()).To(BeEquivalentTo("{QueueName: queue, Depth: <nil>}"))
		})
	})
})
