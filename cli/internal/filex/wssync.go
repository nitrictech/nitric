package filex

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/samber/lo"
	"golang.org/x/net/websocket"
)

// Message represents a file change notification with contents
type Message struct {
	Contents string `json:"contents"`
}

// FileSyncer manages WebSocket connections and file watching
type FileSyncer struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan Message
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	incoming   chan Message
	filePath   string
	file       *os.File
	mutex      sync.RWMutex
}

// NewFileWatcher creates a new file watcher instance
func NewFileSyncer(filePath string) *FileSyncer {
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

	return &FileSyncer{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan Message),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		incoming:   make(chan Message),
		filePath:   filePath,
		file:       file,
	}
}

// Close closes the file handle
func (fw *FileSyncer) Close() {
	if fw.file != nil {
		fw.file.Close()
	}
}

// Run starts the WebSocket server
func (fw *FileSyncer) Run() {
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
			_, err := fw.file.Write([]byte(message.Contents))
			if err != nil {
				log.Printf("Error writing to file %s: %v", fw.filePath, err)
			} else {
				log.Printf("File %s updated from client", fw.filePath)
			}
		}
	}
}

// handleWebSocket handles WebSocket connections
func (fw *FileSyncer) handleWebSocket(ws *websocket.Conn) {
	fw.register <- ws

	// Handle client disconnection
	defer func() {
		fw.unregister <- ws
	}()

	// Handle incoming messages
	for {
		var message Message
		err := websocket.JSON.Receive(ws, &message)
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			break
		}

		// Send the message to the incoming channel for processing
		fw.incoming <- message
	}
}

// StartServer starts the WebSocket server and file watcher
func (fw *FileSyncer) StartServer() error {
	go fw.Run()

	httpServer := websocket.Server{
		Handler: websocket.Handler(fw.handleWebSocket),
		Config: websocket.Config{
			Origin: &url.URL{Scheme: "http", Host: "*"},
		},
	}

	// Set up HTTP server for WebSocket connections
	http.Handle("/ws", httpServer)

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return err
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Println("Starting WebSocket server on", listener.Addr().String())
		if err := http.Serve(listener, nil); err != nil {
			log.Fatal("WebSocket server error:", err)
		}
	}()

	// Start file watcher
	go fw.watchFile()

	return nil
}

// watchFile watches the file for changes and broadcasts updates
func (fw *FileSyncer) watchFile() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Add the file to the watcher
	err = watcher.Add(fw.filePath)
	if err != nil {
		log.Fatal(err)
	}

	var cancel, debounced func()
	for event := range watcher.Events {
		if event.Has(fsnotify.Write) {
			if cancel != nil {
				cancel()
			}

			debounced, cancel = lo.NewDebounce(100*time.Millisecond, func() {
				// Read the current contents of the file
				fw.file.Seek(0, 0) // Seek to beginning
				contents, err := io.ReadAll(fw.file)
				if err != nil {
					log.Printf("Error reading file %s: %v", event.Name, err)
					return
				}

				// Broadcast file change with contents to WebSocket clients
				message := Message{
					Contents: string(contents),
				}
				fw.broadcast <- message
			})
			debounced()
		}
	}
}
