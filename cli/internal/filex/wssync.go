package filex

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/invopop/yaml"
	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/samber/lo"
	"golang.org/x/net/websocket"
)

type MessageType string

const (
	MessageTypeNitricSync MessageType = "nitricSync"
)

// Message represents a file change notification with contents
type Message struct {
	Type string `json:"type"`
}

// TODO: Possibly handle multiple event types to sync with the dashboard
type NitricSyncMessage struct {
	Message
	Payload schema.Application `json:"payload"`
}

// WebsocketServerSync manages WebSocket connections and file watching
type WebsocketServerSync struct {
	listener *net.Listener
	clients  map[*websocket.Conn]bool
	// TODO: Update to interface
	broadcast  chan NitricSyncMessage
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	// TODO: Update to interface
	incoming chan NitricSyncMessage
	filePath string
	file     *os.File
	mutex    sync.RWMutex
	debounce time.Duration
}

// Implement additional constructor options
type WebsocketServerSyncOption func(*WebsocketServerSync)

func WithDebounce(debounce time.Duration) WebsocketServerSyncOption {
	return func(ws *WebsocketServerSync) {
		ws.debounce = debounce
	}
}

func WithListener(listener net.Listener) WebsocketServerSyncOption {
	return func(ws *WebsocketServerSync) {
		ws.listener = &listener
	}
}

const defaultDebounce = time.Millisecond * 100

// NewWebsocketServerSync creates a Websocket file syncing service
func NewWebsocketServerSync(filePath string, options ...WebsocketServerSyncOption) *WebsocketServerSync {
	// Open the file and keep it open
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		log.Printf("Error opening file %s: %v", filePath, err)
		// Create the file if it doesn't exist
		file, err = os.Create(filePath)
		if err != nil {
			log.Fatal("Error creating file:", err)
		}
	}

	ws := &WebsocketServerSync{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan NitricSyncMessage),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		incoming:   make(chan NitricSyncMessage),
		filePath:   filePath,
		file:       file,
		debounce:   defaultDebounce,
	}

	// Apply options
	for _, option := range options {
		option(ws)
	}

	return ws
}

// Close closes the file handle
func (fw *WebsocketServerSync) Close() {
	if fw.file != nil {
		fw.file.Close()
	}
}

// run starts the WebSocket server
func (fw *WebsocketServerSync) run() {
	for {
		select {
		case client := <-fw.register:
			fw.mutex.Lock()
			fw.clients[client] = true
			fw.mutex.Unlock()
			log.Printf("Client connected. Total clients: %d", len(fw.clients))

		case client := <-fw.unregister:
			fw.mutex.Lock()
			delete(fw.clients, client)
			fw.mutex.Unlock()
			client.Close()
			log.Printf("Client disconnected. Total clients: %d", len(fw.clients))

		case message := <-fw.broadcast:
			messageJSON, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				continue
			}

			fw.mutex.RLock()
			for client := range fw.clients {
				_, err := client.Write(messageJSON)
				if err != nil {
					log.Printf("Error broadcasting message: %v", err)
					client.Close()
					delete(fw.clients, client)
				}
			}
			fw.mutex.RUnlock()

		case message := <-fw.incoming:
			// Write the incoming content to the file
			fw.file.Seek(0, 0)  // Seek to beginning
			fw.file.Truncate(0) // Clear the file
			yamlContents, err := yaml.Marshal(message.Payload)
			_, err = fw.file.Write(yamlContents)
			if err != nil {
				log.Printf("Error writing to file %s: %v", fw.filePath, err)
			} else {
				log.Printf("File %s updated from client", fw.filePath)
			}
		}
	}
}

// handleWebSocket handles WebSocket connections
func (fw *WebsocketServerSync) handleWebSocket(ws *websocket.Conn) {
	fw.register <- ws

	// Handle client disconnection
	defer func() {
		fw.unregister <- ws
	}()

	// Handle incoming messages
	for {
		var message NitricSyncMessage
		err := websocket.JSON.Receive(ws, &message)
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			break
		}

		// Send the message to the incoming channel for processing
		fw.incoming <- message
	}
}

// Start starts the WebSocket server and file watcher
func (fw *WebsocketServerSync) Start() error {
	go fw.run()

	httpServer := websocket.Server{
		Handler: websocket.Handler(fw.handleWebSocket),
		Config: websocket.Config{
			Origin: &url.URL{Scheme: "http", Host: "*"},
		},
	}

	// Set up HTTP server for WebSocket connections
	http.Handle("/ws", httpServer)

	if fw.listener == nil {
		listener, err := net.Listen("tcp", "localhost:0")
		if err != nil {
			return err
		}
		fw.listener = &listener
	}

	address := (*fw.listener).Addr().String()
	// Start HTTP server in a goroutine
	go func() {
		log.Println("Starting WebSocket server on", address)
		if err := http.Serve(*fw.listener, nil); err != nil {
			log.Fatal("WebSocket server error:", err)
		}
	}()

	// Block on watch file and return in case of errors
	return fw.watchFile()
}

// watchFile watches the file for changes and broadcasts updates
func (fw *WebsocketServerSync) watchFile() error {
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
				// Read the current contents of the file
				fw.file.Seek(0, 0) // Seek to beginning
				contents, err := io.ReadAll(fw.file)
				if err != nil {
					fileError = err
					return
				}

				application, schemaResult, err := schema.ApplicationFromYaml(string(contents))
				if err != nil {
					fmt.Println("Error parsing application from yaml:", err)
					return
				} else if schemaResult != nil && len(schemaResult.Errors()) > 0 {
					fmt.Println("Errors parsing application from yaml:", schemaResult.Errors())
					return
				}

				// Broadcast file change with contents to WebSocket clients
				message := NitricSyncMessage{
					Message: Message{
						Type: string(MessageTypeNitricSync),
					},
					Payload: *application,
				}
				fw.broadcast <- message
			})
			debounced()
			if fileError != nil {
				return err
			}
		}
	}

	return nil
}
