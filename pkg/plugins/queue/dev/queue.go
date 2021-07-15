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
	utils2 "github.com/nitric-dev/membrane/pkg/utils"
	"os"
	"strings"
	"time"

	"github.com/asdine/storm"
	"github.com/nitric-dev/membrane/pkg/sdk"
	"go.etcd.io/bbolt"
)

const DEFAULT_DIR = "nitric/queues/"

type DevQueueService struct {
	sdk.UnimplementedQueuePlugin
	dbDir string
}

type Item struct {
	ID   int `storm:"id,increment"` // primary key with auto increment
	Data []byte
}

func (s *DevQueueService) Send(queue string, task sdk.NitricTask) error {
	if queue == "" {
		return fmt.Errorf("provide non-blank queue")
	}

	db, err := s.createDb(queue)
	if err != nil {
		return err
	}
	defer db.Close()

	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	item := Item{
		Data: data,
	}

	err = db.Save(&item)
	if err != nil {
		return fmt.Errorf("Error sending %s : %v", task, err)
	}

	return nil
}

func (s *DevQueueService) SendBatch(queue string, tasks []sdk.NitricTask) (*sdk.SendBatchResponse, error) {
	if queue == "" {
		return nil, fmt.Errorf("provide non-blank queue")
	}
	if tasks == nil {
		return nil, fmt.Errorf("provide non-nil tasks")
	}

	db, err := s.createDb(queue)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	for _, task := range tasks {
		data, err := json.Marshal(task)
		if err != nil {
			return nil, err
		}

		item := Item{
			Data: data,
		}

		err = db.Save(&item)
		if err != nil {
			return nil, fmt.Errorf("Error sending %s : %v", task, err)
		}
	}

	return &sdk.SendBatchResponse{
		FailedTasks: make([]*sdk.FailedTask, 0),
	}, nil
}

func (s *DevQueueService) Receive(options sdk.ReceiveOptions) ([]sdk.NitricTask, error) {
	if options.QueueName == "" {
		return nil, fmt.Errorf("provide non-blank options.queue")
	}

	db, err := s.createDb(options.QueueName)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var items []Item
	err = db.All(&items, storm.Limit(int(*options.Depth)))

	poppedTasks := make([]sdk.NitricTask, 0)
	for _, item := range items {
		var task sdk.NitricTask
		err := json.Unmarshal(item.Data, &task)
		if err != nil {
			return nil, err
		}
		poppedTasks = append(poppedTasks, task)

		err = db.DeleteStruct(&item)
		if err != nil {
			return nil, err
		}
	}

	return poppedTasks, nil
}

// Completes a previously popped queue item
func (s *DevQueueService) Complete(queue string, leaseId string) error {
	if queue == "" {
		return fmt.Errorf("provide non-blank queue")
	}
	if leaseId == "" {
		return fmt.Errorf("provide non-blank leaseId")
	}
	return nil
}

func New() (sdk.QueueService, error) {
	dbDir := utils2.GetEnv("LOCAL_QUEUE_DIR", DEFAULT_DIR)

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
	dbPath := s.dbDir + strings.ToLower(queue) + ".db"

	options := storm.BoltOptions(0600, &bbolt.Options{Timeout: 1 * time.Second})
	db, err := storm.Open(dbPath, options)
	if err != nil {
		return nil, err
	}

	return db, nil
}
