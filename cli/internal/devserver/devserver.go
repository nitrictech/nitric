package devserver

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type Message[T any] struct {
	Type    string `json:"type"`
	Payload T      `json:"payload"`
}

// DevWebsockerServer manages WebSocket connections and file watching
type DevWebsockerServer struct {
	clients     map[*websocket.Conn]bool
	broadcast   chan Message[any]
	register    chan *websocket.Conn
	unregister  chan *websocket.Conn
	incoming    chan json.RawMessage
	listener    *net.Listener
	mutex       sync.RWMutex
	subscribers map[string]Subscriber
}

// Implement additional constructor options
type WebsocketServerSyncOption func(*DevWebsockerServer)

func WithListener(listener net.Listener) WebsocketServerSyncOption {
	return func(ws *DevWebsockerServer) {
		ws.listener = &listener
	}
}

const defaultDebounce = time.Millisecond * 100

type BroadcastFunc func(Message[any])

// Broadcast a message to connected clients
func (fw *DevWebsockerServer) Broadcast(message Message[any]) {
	fw.broadcast <- message
}

func (fw *DevWebsockerServer) unsubscribe(subscriberId string) {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	delete(fw.subscribers, subscriberId)
}

func (fw *DevWebsockerServer) notify(message json.RawMessage) {
	fw.mutex.RLock()
	defer fw.mutex.RUnlock()
	for _, subscriber := range fw.subscribers {
		subscriber.OnMessage(message)
	}
}

type SendFunc func(message Message[any])

func (fw *DevWebsockerServer) sendToClient(client *websocket.Conn) SendFunc {
	return func(message Message[any]) {
		messageJSON, err := json.Marshal(message)
		if err != nil {
			log.Printf("Error marshaling message: %v", err)
			return
		}

		client.Write(messageJSON)
	}
}

type Subscriber interface {
	OnMessage(message json.RawMessage)
	// Provide a function reference that allows the subscriber to send messages to the newly connected client
	OnConnect(send SendFunc)
}

func (fw *DevWebsockerServer) Subscribe(subscriber Subscriber) func() {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()
	// generate a subscriber id
	subscriberId := uuid.New().String()
	fw.subscribers[subscriberId] = subscriber

	return func() {
		fw.unsubscribe(subscriberId)
	}
}

// NewDevWebsocketServer creates a Websocket file syncing service
func NewDevWebsocketServer(options ...WebsocketServerSyncOption) *DevWebsockerServer {
	ws := &DevWebsockerServer{
		clients:     make(map[*websocket.Conn]bool),
		broadcast:   make(chan Message[any]),
		register:    make(chan *websocket.Conn),
		unregister:  make(chan *websocket.Conn),
		incoming:    make(chan json.RawMessage),
		subscribers: make(map[string]Subscriber),
	}

	// Apply options
	for _, option := range options {
		option(ws)
	}

	return ws
}

// run starts the WebSocket server
func (fw *DevWebsockerServer) run() {
	for {
		select {
		case client := <-fw.register:
			fw.mutex.Lock()
			fw.clients[client] = true
			fw.mutex.Unlock()
			// Notify all subscribers that a new client has connected and allow them to message that client directly
			send := fw.sendToClient(client)
			for _, subscriber := range fw.subscribers {
				subscriber.OnConnect(send)
			}
			// log.Printf("Client connected. Total clients: %d", len(fw.clients))

		case client := <-fw.unregister:
			fw.mutex.Lock()
			delete(fw.clients, client)
			fw.mutex.Unlock()
			client.Close()
			// log.Printf("Client disconnected. Total clients: %d", len(fw.clients))

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
			// notify all subscribers
			fw.notify(message)
		}
	}
}

// handleWebSocket handles WebSocket connections
func (fw *DevWebsockerServer) handleWebSocket(ws *websocket.Conn) {
	fw.register <- ws

	// Handle client disconnection
	defer func() {
		fw.unregister <- ws
	}()

	// Handle incoming messages
	for {

		var message json.RawMessage
		err := websocket.JSON.Receive(ws, &message)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error receiving message: %v", err)
			}
			break
		}

		// Send the message to the incoming channel for processing
		fw.incoming <- message
	}
}

// Start starts the WebSocket server and file watcher
func (fw *DevWebsockerServer) Start() error {
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

	listener := *fw.listener

	// Start HTTP server in a goroutine
	go func() {
		if err := http.Serve(listener, nil); err != nil {
			log.Fatal("WebSocket server error:", err)
		}
	}()

	return http.Serve(listener, nil)
}
