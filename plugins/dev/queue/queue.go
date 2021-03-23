package queue_service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nitric-dev/membrane/plugins/dev/ifaces"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"github.com/nitric-dev/membrane/utils"
)

type DefaultQueueDriver struct {
	ifaces.UnimplementedStorageDriver
}

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

type DevQueueService struct {
	sdk.UnimplementedQueuePlugin
	driver   ifaces.StorageDriver
	queueDir string
}

func (s *DevQueueService) Send(queue string, event sdk.NitricEvent) error {
	if err := s.driver.EnsureDirExists(s.queueDir); err == nil {
		fileName := fmt.Sprintf("%s%s", s.queueDir, queue)

		var existingQueue []sdk.NitricEvent
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
			existingQueue = make([]sdk.NitricEvent, 0)
		}

		newQueue := append(existingQueue, event)

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

func (s *DevQueueService) SendBatch(queue string, events []sdk.NitricEvent) (*sdk.SendBatchResponse, error) {
	if err := s.driver.EnsureDirExists(s.queueDir); err == nil {
		fileName := fmt.Sprintf("%s%s", s.queueDir, queue)

		var existingQueue []sdk.NitricEvent
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
			existingQueue = make([]sdk.NitricEvent, 0)
		}

		newQueue := existingQueue
		for _, evt := range events {
			// Add indirected event references to the new queue...
			newQueue = append(newQueue, evt)
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
		FailedMessages: make([]*sdk.FailedMessage, 0),
	}, nil
}

func (s *DevQueueService) Receive(options sdk.ReceiveOptions) ([]sdk.NitricQueueItem, error) {
	if err := s.driver.EnsureDirExists(s.queueDir); err == nil {
		fileName := fmt.Sprintf("%s%s", s.queueDir, options.QueueName)

		var existingQueue []sdk.NitricEvent
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
			return []sdk.NitricQueueItem{}, nil
		}

		poppedItems := make([]sdk.NitricQueueItem, 0)
		remainingItems := make([]sdk.NitricEvent, 0)
		for i, evt := range existingQueue {
			if uint32(i) < *options.Depth {
				poppedItems = append(poppedItems, sdk.NitricQueueItem{
					Event:   evt,
					LeaseId: evt.RequestId,
				})
			} else {
				remainingItems = append(remainingItems, evt)
			}
		}

		// Store the remaining items back to the queue file.
		if queueByte, err := json.Marshal(&remainingItems); err == nil {
			// Write the new queue, to a file named after the queue
			if err := s.driver.WriteFile(fileName, queueByte, os.ModePerm); err != nil {
				return nil, err
			}
		}
		return poppedItems, nil
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

func NewWithStorageDriver(driver ifaces.StorageDriver) (sdk.QueueService, error) {
	queueDir := utils.GetEnv("LOCAL_QUEUE_DIR", "/nitric/queues/")

	return &DevQueueService{
		driver:   driver,
		queueDir: queueDir,
	}, nil
}
