package queue_plugin

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

func (s *DefaultQueueDriver) ReadFile(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

type LocalQueuePlugin struct {
	sdk.UnimplementedQueuePlugin
	driver   ifaces.StorageDriver
	queueDir string
}

func (s *LocalQueuePlugin) Push(queue string, events []*sdk.NitricEvent) (*sdk.PushResponse, error) {
	if err := s.storageDriver.EnsureDirExists(s.queueDir); err == nil {
		fileName := fmt.Sprintf("%s%s", s.queueDir, queue)

		var existingQueue []sdk.NitricEvent
		// See if the queue exists first...
		if os.Stat(fileName); !os.IsNotExist(err) {
			// Read the file first
			if b, err := s.driver.ReadFile(fileName); err == nil {
				if err := json.Unmarshal(b, &existingQueue); err {
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
			newQueue := append(newQueue, *evt)
		}

		if queueByte, err := json.Marshal(&newQueue); err == nil {
			// Write the new queue, to a file named after the queue
			if err := s.storageDriver.WriteFile(fileName, payload, os.ModePerm); err != nil {
				return nil, err
			}
		}
	} else {
		return nil, err
	}

	return &sdk.PushResponse{
		FailedMessages: make([]*sdk.NitricEvent, 0),
	}, nil
}

func New() (sdk.QueuePlugin, error) {
	queueDir := utils.GetEnv("LOCAL_QUEUE_DIR", "/nitric/queues/")

	return &LocalQueuePlugin{
		driver:   &DefaultQueueDriver{},
		queueDir: queueDir,
	}
}
