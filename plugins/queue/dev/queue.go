package queue_service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/utils"
)

type DefaultQueueDriver struct{}

// EnsureDirExists - Recurively creates directories for the given path
func (s *DefaultQueueDriver) EnsureDirExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

// WriteFile - Writes the given byte array to the given path
func (s *DefaultQueueDriver) WriteFile(file string, contents []byte, fileMode os.FileMode) error {
	return ioutil.WriteFile(file, contents, fileMode)
}

func (s *DefaultQueueDriver) ExistsOrFail(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	return nil
}

func (s *DefaultQueueDriver) ReadFile(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

func (s *DefaultQueueDriver) DeleteFile(file string) error {
	return os.Remove(file)
}

type DevQueueService struct {
	sdk.UnimplementedQueuePlugin
	driver   StorageDriver
	queueDir string
}

func (s *DevQueueService) Send(queue string, task sdk.NitricTask) error {
	if err := s.driver.EnsureDirExists(s.queueDir); err == nil {
		fileName := fmt.Sprintf("%s%s", s.queueDir, queue)

		var existingQueue []sdk.NitricTask
		// See if the queue exists first...
		if err := s.driver.ExistsOrFail(fileName); err == nil {
			// Read the file first
			if b, err := s.driver.ReadFile(fileName); err == nil {
				if err := json.Unmarshal(b, &existingQueue); err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			existingQueue = make([]sdk.NitricTask, 0)
		}

		newQueue := append(existingQueue, task)

		if queueByte, err := json.Marshal(&newQueue); err == nil {
			// Write the new queue, to a file named after the queue
			if err := s.driver.WriteFile(fileName, queueByte, os.ModePerm); err != nil {
				return err
			}
		}
	} else {
		return err
	}

	return nil
}

func (s *DevQueueService) SendBatch(queue string, tasks []sdk.NitricTask) (*sdk.SendBatchResponse, error) {
	if err := s.driver.EnsureDirExists(s.queueDir); err == nil {
		fileName := fmt.Sprintf("%s%s", s.queueDir, queue)

		var existingQueue []sdk.NitricTask
		// See if the queue exists first...
		if err := s.driver.ExistsOrFail(fileName); err == nil {
			// Read the file first
			if b, err := s.driver.ReadFile(fileName); err == nil {
				if err := json.Unmarshal(b, &existingQueue); err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			existingQueue = make([]sdk.NitricTask, 0)
		}

		newQueue := existingQueue
		for _, task := range tasks {
			// Add indirected task references to the new queue...
			newQueue = append(newQueue, task)
		}

		if queueByte, err := json.Marshal(&newQueue); err == nil {
			// Write the new queue, to a file named after the queue
			if err := s.driver.WriteFile(fileName, queueByte, os.ModePerm); err != nil {
				return nil, err
			}
		}
	} else {
		return nil, err
	}

	return &sdk.SendBatchResponse{
		FailedTasks: make([]*sdk.FailedTask, 0),
	}, nil
}

func (s *DevQueueService) Receive(options sdk.ReceiveOptions) ([]sdk.NitricTask, error) {
	if err := s.driver.EnsureDirExists(s.queueDir); err == nil {
		fileName := fmt.Sprintf("%s%s", s.queueDir, options.QueueName)

		var existingQueue []sdk.NitricTask
		// See if the queue exists first...
		if err := s.driver.ExistsOrFail(fileName); err == nil {
			// Read the file first
			if b, err := s.driver.ReadFile(fileName); err == nil {
				if err := json.Unmarshal(b, &existingQueue); err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("queue not found")
		}

		if len(existingQueue) == 0 {
			return []sdk.NitricTask{}, nil
		}

		poppedTasks := make([]sdk.NitricTask, 0)
		remainingItems := make([]sdk.NitricTask, 0)
		for i, task := range existingQueue {
			if uint32(i) < *options.Depth {
				poppedTasks = append(poppedTasks, sdk.NitricTask{
					ID:          task.ID,
					Payload:     task.Payload,
					PayloadType: task.PayloadType,
					LeaseID:     task.LeaseID,
				})
			} else {
				remainingItems = append(remainingItems, task)
			}
		}

		// Store the remaining items back to the queue file.
		if queueByte, err := json.Marshal(&remainingItems); err == nil {
			// Write the new queue, to a file named after the queue
			if err := s.driver.WriteFile(fileName, queueByte, os.ModePerm); err != nil {
				return nil, err
			}
		}
		return poppedTasks, nil
	} else {
		return nil, err
	}
}

// Completes a previously popped queue item
func (s *DevQueueService) Complete(queue string, leaseId string) error {
	return nil
}

func New() (sdk.QueueService, error) {
	queueDir := utils.GetEnv("LOCAL_QUEUE_DIR", "/nitric/queues/")

	return &DevQueueService{
		driver:   &DefaultQueueDriver{},
		queueDir: queueDir,
	}, nil
}

func NewWithStorageDriver(driver StorageDriver) (sdk.QueueService, error) {
	queueDir := utils.GetEnv("LOCAL_QUEUE_DIR", "/nitric/queues/")

	return &DevQueueService{
		driver:   driver,
		queueDir: queueDir,
	}, nil
}
