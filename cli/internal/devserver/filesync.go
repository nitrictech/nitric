package devserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

type NitricFileSync struct {
	filePath         string
	file             *os.File
	debounce         time.Duration
	broadcast        BroadcastFunc
	lastSyncContents []byte
}

type FileSyncMessage Message[schema.Application]

type FileSyncOption func(*NitricFileSync)

func WithDebounce(debounce time.Duration) FileSyncOption {
	return func(fs *NitricFileSync) {
		fs.debounce = debounce
	}
}

func (fs *NitricFileSync) getFileContents() (*schema.Application, []byte, error) {
	fs.file.Seek(0, 0) // Seek to beginning
	contents, err := io.ReadAll(fs.file)
	if err != nil {
		return nil, nil, err
	}

	application, schemaResult, err := schema.ApplicationFromYaml(string(contents))
	if err != nil {
		fmt.Println("Error parsing application from yaml:", err)
		return nil, contents, err
	} else if schemaResult != nil && len(schemaResult.Errors()) > 0 {
		// Wrap the schema errors in a new error
		return nil, contents, fmt.Errorf("Errors parsing application from yaml: %v", schemaResult.Errors())
	}

	return application, contents, nil
}

func (fs *NitricFileSync) setFileContents(contents []byte) error {
	_, err := fs.file.Seek(0, 0)
	if err != nil {
		return err
	}
	err = fs.file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = fs.file.Write(contents)
	return err
}

func (fs *NitricFileSync) OnConnect(send SendFunc) {
	application, _, err := fs.getFileContents()
	if err != nil {
		return
	}
	// Send initial state to a newly connected client
	send(Message[any]{
		Type:    "nitricSync",
		Payload: *application,
	})
}

func (fs *NitricFileSync) OnMessage(message json.RawMessage) {
	var fileSyncMessage FileSyncMessage
	err := json.Unmarshal(message, &fileSyncMessage)
	if err != nil {
		return
	}

	// Not the right message type continue
	if fileSyncMessage.Type != "nitricSync" {
		return
	}

	// yaml.Marshal() defaults to 4 spaces for indentation
	// 2 is more common, so we use that here
	// TODO: In future if we can preserve the original indentation we should do that
	var buffer bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&buffer)
	yamlEncoder.SetIndent(2)

	err = yamlEncoder.Encode(fileSyncMessage.Payload)
	if err != nil {
		fmt.Println("Error marshalling application to yaml:", err)
		return
	}

	err = fs.setFileContents(buffer.Bytes())
	if err != nil {
		fmt.Println("Error setting file contents:", err)
		return
	}

	fs.lastSyncContents = buffer.Bytes()
}

func (fs *NitricFileSync) Close() error {
	return fs.file.Close()
}

func NewFileSync(filePath string, broadcast BroadcastFunc, options ...FileSyncOption) (*NitricFileSync, error) {
	var err error
	fs := &NitricFileSync{
		filePath:  filePath,
		broadcast: broadcast,
	}

	for _, option := range options {
		option(fs)
	}

	fs.file, err = os.OpenFile(fs.filePath, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return fs, nil
}

func (fs *NitricFileSync) Start() error {
	return fs.watchFile()
}

// watchFile watches the file for changes and broadcasts updates
func (fw *NitricFileSync) watchFile() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// Add the file to the watcher
	err = watcher.Add(fw.filePath)
	if err != nil {
		return err
	}

	var cancel, debounced func()
	for event := range watcher.Events {
		if event.Has(fsnotify.Write) {
			if cancel != nil {
				cancel()
			}

			var fileError error = nil
			debounced, cancel = lo.NewDebounce(fw.debounce, func() {
				application, contents, err := fw.getFileContents()
				if err != nil {
					fileError = err
					return
				}

				if bytes.Equal(fw.lastSyncContents, contents) {
					return
				}

				fw.broadcast(Message[any]{
					Type:    "nitricSync",
					Payload: *application,
				})
			})
			debounced()
			if fileError != nil {
				return err
			}
		}
	}

	return nil
}
