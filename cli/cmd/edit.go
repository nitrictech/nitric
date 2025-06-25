package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
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
	mutex      sync.RWMutex
}

// NewWebSocketServer creates a new WebSocket server instance
func NewWebSocketServer(filePath string) *WebSocketServer {
	return &WebSocketServer{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan Message),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		incoming:   make(chan Message),
		filePath:   filePath,
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
			err := os.WriteFile(s.filePath, []byte(message.Contents), 0644)
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

		// Start HTTP server in a goroutine
		go func() {
			log.Println("Starting WebSocket server on :8080")
			if err := http.ListenAndServe("localhost:8080", nil); err != nil {
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
			for event := range watcher.Events {
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("File modified:", event.Name)

					// Read the current contents of the file
					file, err := os.Open(event.Name)
					if err != nil {
						log.Printf("Error opening file %s: %v", event.Name, err)
						continue
					}

					contents, err := io.ReadAll(file)
					file.Close()
					if err != nil {
						log.Printf("Error reading file %s: %v", event.Name, err)
						continue
					}

					// Broadcast file change with contents to WebSocket clients
					message := Message{
						Contents: string(contents),
					}
					wsServer.broadcast <- message
				}
			}
		}(done)

		// Add the nitric.yaml file to the watcher
		err = watcher.Add("nitric.yaml")
		if err != nil {
			log.Fatal(err)
		}

		// I want to start a websocket server here that will broadcast yaml file changes to connected clients

		// Block on done
		<-done
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
