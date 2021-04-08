package kv_service

import (
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/nitric-dev/membrane/plugins/dev/ifaces"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/utils"
)

type DevKVService struct {
	sdk.UnimplementedKeyValuePlugin
	db ifaces.ScribbleIface
}

func (s *DevKVService) Get(collection string, key string) (map[string]interface{}, error) {
	value := make(map[string]interface{})
	err := s.db.Read(collection, key, &value)

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (s *DevKVService) Put(collection string, key string, value map[string]interface{}) error {
	return s.db.Write(collection, key, value)
}

func (s *DevKVService) Delete(collection string, key string) error {
	error := s.db.Delete(collection, key)

	if error != nil {
		return error
	}

	return nil
}

// New - Create a new dev KV plugin
func New() (sdk.KeyValueService, error) {
	dbDir := utils.GetEnv("LOCAL_DB_DIR", "/nitric/")
	db, err := scribble.New(dbDir, nil)

	if err != nil {
		return nil, err
	}

	return &DevKVService{
		db: db,
	}, nil
}

func NewWithDB(db ifaces.ScribbleIface) (sdk.KeyValueService, error) {
	return &DevKVService{
		db: db,
	}, nil
}
