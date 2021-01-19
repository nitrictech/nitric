package storage_plugin

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nitric-dev/membrane/plugins/sdk"
	"github.com/nitric-dev/membrane/utils"
)

type LocalStoragePlugin struct {
	sdk.UnimplementedStoragePlugin
	storeDir string
}

func ensureDirExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func (s *LocalStoragePlugin) Put(bucket string, key string, payload []byte) error {
	bucketName := fmt.Sprintf("%s%s/", s.storeDir, bucket)

	if err := ensureDirExists(bucketName); err == nil {
		fileName := fmt.Sprintf("%s%s", bucketName, key)

		if err := ioutil.WriteFile(fileName, payload, os.ModePerm); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

// Retrieve an item from a bucket
func (s *LocalStoragePlugin) Get(bucket string, key string) ([]byte, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

// Create new DynamoDB documents server
// XXX: No External Args for function atm (currently the plugin loader does not pass any argument information)
func New() (sdk.StoragePlugin, error) {
	storeDir := utils.GetEnv("LOCAL_BLOB_DIR", "/nitric/buckets/")

	if err := ensureDirExists(storeDir); err != nil {
		return nil, err
	}

	return &LocalStoragePlugin{
		storeDir: storeDir,
	}, nil
}
