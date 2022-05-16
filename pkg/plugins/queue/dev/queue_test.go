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

package queue_service_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/asdine/storm"
	"go.etcd.io/bbolt"

	"github.com/nitrictech/nitric/pkg/plugins/queue"
	queue_service "github.com/nitrictech/nitric/pkg/plugins/queue/dev"
	"github.com/nitrictech/nitric/pkg/utils"
)

var task1 = queue.NitricTask{
	ID:          "1234",
	PayloadType: "test-payload",
	Payload: map[string]interface{}{
		"Test": "Test 1",
	},
}
var task2 = queue.NitricTask{
	ID:          "2345",
	PayloadType: "test-payload",
	Payload: map[string]interface{}{
		"Test": "Test 3",
	},
}
var task3 = queue.NitricTask{
	ID:          "3456",
	PayloadType: "test-payload",
	Payload: map[string]interface{}{
		"Test": "Test 3",
	},
}
var task4 = queue.NitricTask{
	ID:          "4567",
	PayloadType: "test-payload",
	Payload: map[string]interface{}{
		"Test": "Test 4",
	},
}

var local_queue_directory = utils.GetRelativeDevPath(queue_service.DEV_SUB_DIRECTORY)

var _ = Describe("Queue", func() {

	queuePlugin, err := queue_service.New()
	if err != nil {
		Expect(err).ShouldNot(HaveOccurred())
	}

	AfterEach(func() {
		err := os.RemoveAll(local_queue_directory)
		if err != nil {
			Expect(err).ShouldNot(HaveOccurred())
		}

		_, err = os.Stat(local_queue_directory)
		if os.IsNotExist(err) {
			// Make diretory if not present
			err := os.Mkdir(local_queue_directory, 0777)
			if err != nil {
				Expect(err).ShouldNot(HaveOccurred())
			}
		}
	})

	AfterSuite(func() {
		err := os.RemoveAll(local_queue_directory)
		if err == nil {
			os.Remove(local_queue_directory)
			os.Remove(utils.GetDevVolumePath())
		}
	})

	Context("Send", func() {
		When("The queue is empty", func() {
			It("Should store the events in the queue", func() {
				err := queuePlugin.Send("test", task1)
				Expect(err).ShouldNot(HaveOccurred())

				storedTasks := GetAllTasks("test")
				Expect(storedTasks).NotTo(BeNil())
				Expect(storedTasks).To(HaveLen(1))
				Expect(storedTasks[0]).To(BeEquivalentTo(task1))
			})
		})

		When("The queue is not empty", func() {
			It("Should append to the existing queue", func() {
				err := queuePlugin.Send("test", task1)
				Expect(err).ShouldNot(HaveOccurred())

				storedTasks := GetAllTasks("test")
				Expect(storedTasks).NotTo(BeNil())
				Expect(storedTasks).To(HaveLen(1))
				Expect(storedTasks[0]).To(BeEquivalentTo(task1))

				err = queuePlugin.Send("test", task2)
				Expect(err).ShouldNot(HaveOccurred())

				storedTasks = GetAllTasks("test")
				Expect(storedTasks).NotTo(BeNil())
				Expect(storedTasks).To(HaveLen(2))
				Expect(storedTasks[0]).To(BeEquivalentTo(task1))
				Expect(storedTasks[1]).To(BeEquivalentTo(task2))
			})
		})
	})

	Context("SendBatch", func() {
		When("The queue is empty", func() {
			tasks := []queue.NitricTask{task1, task2}
			It("Should store the events in the queue", func() {
				resp, err := queuePlugin.SendBatch("test", tasks)
				Expect(resp.FailedTasks).To(BeEmpty())
				Expect(err).ShouldNot(HaveOccurred())

				storedTasks := GetAllTasks("test")
				Expect(storedTasks).NotTo(BeNil())
				Expect(storedTasks).To(HaveLen(2))
				Expect(storedTasks[0]).To(BeEquivalentTo(task1))
				Expect(storedTasks[1]).To(BeEquivalentTo(task2))
			})
		})

		When("The queue is not empty", func() {
			It("Should append to the existing queue", func() {
				batch1 := []queue.NitricTask{task1, task2}
				resp, err := queuePlugin.SendBatch("test", batch1)
				Expect(resp.FailedTasks).To(BeEmpty())
				Expect(err).ShouldNot(HaveOccurred())

				storedTasks := GetAllTasks("test")
				Expect(storedTasks).NotTo(BeNil())
				Expect(storedTasks).To(HaveLen(2))

				batch2 := []queue.NitricTask{task3, task4}
				resp, err = queuePlugin.SendBatch("test", batch2)
				Expect(resp.FailedTasks).To(BeEmpty())
				Expect(err).ShouldNot(HaveOccurred())

				storedTasks = GetAllTasks("test")
				Expect(storedTasks).NotTo(BeNil())
				Expect(storedTasks).To(HaveLen(4))
				Expect(storedTasks[2]).To(BeEquivalentTo(task3))
				Expect(storedTasks[3]).To(BeEquivalentTo(task4))
			})
		})
	})

	Context("Receive", func() {
		When("The queue is empty", func() {
			It("Should return an empty slice of queue items", func() {
				depth := uint32(10)
				items, err := queuePlugin.Receive(queue.ReceiveOptions{
					QueueName: "test",
					Depth:     &depth,
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(items).To(HaveLen(0))
			})
		})

		When("The queue is not empty", func() {
			It("Should append to the existing queue", func() {
				err := queuePlugin.Send("test", task1)
				Expect(err).ShouldNot(HaveOccurred())

				depth := uint32(10)
				items, err := queuePlugin.Receive(queue.ReceiveOptions{
					QueueName: "test",
					Depth:     &depth,
				})
				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("Returning 1 item")
				Expect(items).To(HaveLen(1))

				storedTasks := GetAllTasks("test")
				Expect(storedTasks).NotTo(BeNil())
				Expect(storedTasks).To(HaveLen(0))
			})
		})

		When("The queue depth is 15", func() {
			It("Should return 10 items", func() {

				task := queue.NitricTask{
					ID:          "1234",
					PayloadType: "test-payload",
					Payload: map[string]interface{}{
						"Test": "Test",
					},
				}
				tasks := []queue.NitricTask{}

				// Add 15 items to the queue
				for i := 0; i < 15; i++ {
					tasks = append(tasks, task)
				}

				_, err = queuePlugin.SendBatch("test", tasks)
				Expect(err).ShouldNot(HaveOccurred())
				storedTasks := GetAllTasks("test")
				Expect(storedTasks).NotTo(BeNil())
				Expect(storedTasks).To(HaveLen(15))

				depth := uint32(10)
				items, err := queuePlugin.Receive(queue.ReceiveOptions{
					QueueName: "test",
					Depth:     &depth,
				})
				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("Returning 10 item")
				Expect(items).To(HaveLen(10))

				storedTasks = GetAllTasks("test")
				Expect(storedTasks).NotTo(BeNil())
				Expect(storedTasks).To(HaveLen(5))
			})
		})
	})

	Context("Complete", func() {
		// Currently the local queue complete method is a stub that always returns successfully.
		// We may consider adding more realistic behavior if that is useful in future.
		When("it always returns successfully", func() {
			It("Should retnot return an error", func() {
				err := queuePlugin.Complete("test-queue", "test-id")
				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})

func GetAllTasks(q string) []queue.NitricTask {
	dbPath := filepath.Join(local_queue_directory, strings.ToLower(q)+".db")

	options := storm.BoltOptions(0600, &bbolt.Options{Timeout: 1 * time.Second})
	db, err := storm.Open(dbPath, options)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	var items []queue_service.Item
	err = db.All(&items)
	if err != nil {
		panic(err)
	}

	tasks := make([]queue.NitricTask, 0)
	for _, item := range items {
		var task queue.NitricTask
		err := json.Unmarshal(item.Data, &task)
		if err != nil {
			panic(err)
		}
		tasks = append(tasks, task)
	}

	return tasks
}
