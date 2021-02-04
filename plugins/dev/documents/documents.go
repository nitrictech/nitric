package documents_service

import (
	"fmt"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/nitric-dev/membrane/plugins/dev/ifaces"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"github.com/nitric-dev/membrane/utils"
)

type LocalDocumentService struct {
	sdk.UnimplementedDocumentsPlugin
	db ifaces.ScribbleIface
}

type NitricDocument struct {
	Key   string
	Value map[string]interface{}
}

func (s *LocalDocumentService) CreateDocument(collection string, key string, document map[string]interface{}) error {
	existingDocument := make(map[string]interface{})
	err := s.db.Read(collection, key, &existingDocument)

	// There was an error reading the existing document we'll assume this means that it doesn't exist
	// So we can go ahead with creation of the new document
	if err != nil {
		err := s.db.Write(collection, key, document)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("Document already exists!")
}

func (s *LocalDocumentService) GetDocument(collection string, key string) (map[string]interface{}, error) {
	document := make(map[string]interface{})
	err := s.db.Read(collection, key, &document)

	if err != nil {
		return nil, err
	}

	return document, nil
}

func (s *LocalDocumentService) UpdateDocument(collection string, key string, document map[string]interface{}) error {
	existingDocument := make(map[string]interface{})
	err := s.db.Read(collection, key, &existingDocument)

	// There was an error reading the existing document we'll assume this means that it doesn't exist
	// So we can go ahead with creation of the new document
	if err == nil {
		err := s.db.Write(collection, key, document)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("Document does not exist!")
}

func (s *LocalDocumentService) DeleteDocument(collection string, key string) error {
	error := s.db.Delete(collection, key)

	if error != nil {
		return error
	}

	return nil
}

// Create new DynamoDB documents server
// XXX: No External Args for function atm (currently the plugin loader does not pass any argument information)
func New() (sdk.DocumentService, error) {
	dbDir := utils.GetEnv("LOCAL_DB_DIR", "/nitric/")
	db, err := scribble.New(dbDir, nil)

	if err != nil {
		return nil, err
	}

	return &LocalDocumentService{
		db: db,
	}, nil
}

func NewWithDB(db ifaces.ScribbleIface) (sdk.DocumentService, error) {
	return &LocalDocumentService{
		db: db,
	}, nil
}
