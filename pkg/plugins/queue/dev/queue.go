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

package queue_service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nitric-dev/membrane/pkg/utils"

	"github.com/asdine/storm"
	"github.com/nitric-dev/membrane/pkg/plugins/errors"
	"github.com/nitric-dev/membrane/pkg/plugins/errors/codes"
	"github.com/nitric-dev/membrane/pkg/plugins/queue"
	"go.etcd.io/bbolt"
)

const DEV_SUB_DIRECTORY = "./queues/"

type DevQueueService struct {
	queue.UnimplementedQueuePlugin
	dbDir string
}

type Item struct {
	ID   int `storm:"id,increment"` // primary key with auto increment
	Data []byte
}

func (s *DevQueueService) Send(queue string, task queue.NitricTask) error {
	newErr := errors.ErrorsWithScope(
		"DevQueueService.Send",
		fmt.Sprintf("queue=%s", queue),
	)

	if queue == "" {
		return newErr(
			codes.InvalidArgument,
			"provide non-blank queue",
			nil,
		)
	}

	db, err := s.createDb(queue)
	if err != nil {
		return newErr(
			codes.FailedPrecondition,
			"createDb error",
			err,
		)
	}
	defer db.Close()

	data, err := json.Marshal(task)
	if err != nil {
		return newErr(
			codes.Internal,
			"error marshalling task",
			err,
		)
	}

	item := Item{
		Data: data,
	}

	err = db.Save(&item)
	if err != nil {
		return newErr(
			codes.Internal,
			"error sending task",
			err,
		)
	}

	return nil
}

func (s *DevQueueService) SendBatch(q string, tasks []queue.NitricTask) (*queue.SendBatchResponse, error) {
	newErr := errors.ErrorsWithScope(
		"DevQueueService.SendBatch",
		fmt.Sprintf("queue=%s", q),
	)

	if q == "" {
		return nil, newErr(
			codes.InvalidArgument,
			"provide non-blank queue",
			nil,
		)
	}
	if tasks == nil {
		return nil, newErr(
			codes.InvalidArgument,
			"provide non-nil tasks",
			nil,
		)
	}

	db, err := s.createDb(q)
	if err != nil {
		return nil, newErr(
			codes.FailedPrecondition,
			"createDb error",
			err,
		)
	}
	defer db.Close()

	for _, task := range tasks {
		data, err := json.Marshal(task)
		if err != nil {
			return nil, newErr(
				codes.Internal,
				fmt.Sprintf("error marshalling task: %v", task),
				err,
			)
		}

		item := Item{
			Data: data,
		}

		err = db.Save(&item)
		if err != nil {
			return nil, newErr(
				codes.Internal,
				fmt.Sprintf("error sending task: %v", task),
				err,
			)
		}
	}

	return &queue.SendBatchResponse{
		FailedTasks: make([]*queue.FailedTask, 0),
	}, nil
}

func (s *DevQueueService) Receive(options queue.ReceiveOptions) ([]queue.NitricTask, error) {
	newErr := errors.ErrorsWithScope(
		"DevQueueService.Receive",
		fmt.Sprintf("options=%v", options),
	)

	if options.QueueName == "" {
		return nil, newErr(
			codes.InvalidArgument,
			"provide non-blank options.queue",
			nil,
		)
	}

	db, err := s.createDb(options.QueueName)
	if err != nil {
		return nil, newErr(
			codes.FailedPrecondition,
			"createDb error",
			err,
		)
	}
	defer db.Close()

	var items []Item
	err = db.All(&items, storm.Limit(int(*options.Depth)))

	poppedTasks := make([]queue.NitricTask, 0)
	for _, item := range items {
		var task queue.NitricTask
		err := json.Unmarshal(item.Data, &task)
		if err != nil {
			return nil, newErr(
				codes.Internal,
				"error marshalling task",
				err,
			)
		}
		task.LeaseID = uuid.New().String()
		poppedTasks = append(poppedTasks, task)

		err = db.DeleteStruct(&item)
		if err != nil {
			return nil, newErr(
				codes.Internal,
				"error de-queuing task",
				err,
			)
		}
	}

	return poppedTasks, nil
}

// Completes a previously popped queue item
func (s *DevQueueService) Complete(queue string, leaseId string) error {
	newErr := errors.ErrorsWithScope(
		"DevQueueService.Complete",
		fmt.Sprintf("queue=%s", queue),
	)

	if queue == "" {
		return newErr(
			codes.InvalidArgument,
			"provide non-blank queue",
			nil,
		)
	}
	if leaseId == "" {
		return newErr(
			codes.InvalidArgument,
			"provide non-blank leaseId",
			nil,
		)
	}
	return nil
}

func New() (queue.QueueService, error) {
	dbDir := utils.GetEnv("LOCAL_QUEUE_DIR", utils.GetRelativeDevPath(DEV_SUB_DIRECTORY))

	// Check whether file exists
	_, err := os.Stat(dbDir)
	if os.IsNotExist(err) {
		// Make diretory if not present
		err := os.MkdirAll(dbDir, 0777)
		if err != nil {
			return nil, err
		}
	}

	return &DevQueueService{
		dbDir: dbDir,
	}, nil
}

func (s *DevQueueService) createDb(queue string) (*storm.DB, error) {
	dbPath := filepath.Join(s.dbDir, strings.ToLower(queue)+".db")

	options := storm.BoltOptions(0600, &bbolt.Options{Timeout: 1 * time.Second})
	db, err := storm.Open(dbPath, options)
	if err != nil {
		return nil, err
	}

	return db, nil
}
