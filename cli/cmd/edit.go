package cmd

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
	"github.com/spf13/cobra"
	"golang.org/x/net/websocket"
)

// Message represents a file change notification with contents
type Message struct {
	Contents string `json:"contents"`
}

// WebSocketServer manages connected clients and broadcasts messages
type WebSocketServer struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan Message
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	incoming   chan Message
	filePath   string
	file       *os.File
	mutex      sync.RWMutex
}

// NewWebSocketServer creates a new WebSocket server instance
func NewWebSocketServer(filePath string) *WebSocketServer {
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

	return &WebSocketServer{
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
func (s *WebSocketServer) Close() {
	if s.file != nil {
		s.file.Close()
	}
}

// Run starts the WebSocket server
func (s *WebSocketServer) Run() {
	for {
		select {
		case client := <-s.register:
			s.mutex.Lock()
			s.clients[client] = true
			s.mutex.Unlock()
			log.Printf("Client connected. Total clients: %d", len(s.clients))

		case client := <-s.unregister:
			s.mutex.Lock()
			delete(s.clients, client)
			s.mutex.Unlock()
			client.Close()
			log.Printf("Client disconnected. Total clients: %d", len(s.clients))

		case message := <-s.broadcast:
			messageJSON, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				continue
			}

			s.mutex.RLock()
			for client := range s.clients {
				_, err := client.Write(messageJSON)
				if err != nil {
					log.Printf("Error broadcasting message: %v", err)
					client.Close()
					delete(s.clients, client)
				}
			}
			s.mutex.RUnlock()

		case message := <-s.incoming:
			// Write the incoming content to the file
			s.file.Seek(0, 0)  // Seek to beginning
			s.file.Truncate(0) // Clear the file
			_, err := s.file.Write([]byte(message.Contents))
			if err != nil {
				log.Printf("Error writing to file %s: %v", s.filePath, err)
			} else {
				log.Printf("File %s updated from client", s.filePath)
			}
		}
	}
}

// handleWebSocket handles WebSocket connections
func (s *WebSocketServer) handleWebSocket(ws *websocket.Conn) {
	s.register <- ws

	// Handle client disconnection
	defer func() {
		s.unregister <- ws
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
		s.incoming <- message
	}
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the nitric application",
	Long:  `Edits an application using the nitric.yaml application spec and referenced platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create WebSocket server
		wsServer := NewWebSocketServer("nitric.yaml")
		go wsServer.Run()

		httpServer := websocket.Server{
			Handler: websocket.Handler(wsServer.handleWebSocket),
			Config: websocket.Config{
				Origin: &url.URL{Scheme: "http", Host: "*"},
			},
		}

		// Set up HTTP server for WebSocket connections
		http.Handle("/ws", httpServer)

		listener, err := net.Listen("tcp", "localhost:0")
		cobra.CheckErr(err)
		// Start HTTP server in a goroutine
		go func() {
			log.Println("Starting WebSocket server on", listener.Addr().String())

			if err := http.Serve(listener, nil); err != nil {
				log.Fatal("WebSocket server error:", err)
			}
		}()

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		// Add sync wait for go routine
		done := make(chan bool)
		go func(doneCh chan bool) {
			defer close(doneCh)
			var cancel, debounced func()
			for event := range watcher.Events {
				if event.Has(fsnotify.Write) {
					if cancel != nil {
						cancel()
					}

					debounced, cancel = lo.NewDebounce(100*time.Millisecond, func() {
						// Read the current contents of the file
						wsServer.file.Seek(0, 0) // Seek to beginning
						contents, err := io.ReadAll(wsServer.file)
						if err != nil {
							log.Printf("Error reading file %s: %v", event.Name, err)
						}

						// Broadcast file change with contents to WebSocket clients
						message := Message{
							Contents: string(contents),
						}
						wsServer.broadcast <- message
					})
					debounced()
				}
			}
		}(done)

		// Add the nitric.yaml file to the watcher
		err = watcher.Add("nitric.yaml")
		if err != nil {
			log.Fatal(err)
		}

		// Block on done
		<-done

		// Close the file when done
		wsServer.Close()
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
