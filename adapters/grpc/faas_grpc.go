package grpc

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	pb "github.com/nitric-dev/membrane/interfaces/nitric/v1"
	"github.com/nitric-dev/membrane/sdk"
)

type FaasServer struct {
	pb.UnimplementedFaasServer
	eventPlugin sdk.EventService

	srv pb.Faas_TriggerStreamServer

	// The function is ready to go
	FunctionReady chan bool

	// Each trigger will get a channel back to wait for its response from the function
	// triggerQueue map[string]*pb.TriggerRequest
	// Add a write lock for triggers
	responseQueue sync.Map
	// Add a read lock for responses
}

// Push a trigger onto the queue
func (s *FaasServer) PushTrigger(request *pb.TriggerRequest) (chan *pb.TriggerResponse, error) {

	// Make a channel for this trigger and push it onto the response queue under its id
	// Generate an ID for this request response pair
	ID := uuid.New().String()
	message := &pb.Message{
		Id: ID,
		Content: &pb.Message_TriggerRequest{
			TriggerRequest: request,
		},
	}

	err := s.srv.Send(message)

	if err != nil {
		// There was an error enqueuing the message
		return nil, err
	}

	// Get a lock on the response queue
	returnChan := make(chan *pb.TriggerResponse)
	s.responseQueue.Store(ID, returnChan)

	return returnChan, nil
}

// Recieve messages from the function
func (s *FaasServer) recieveMessages(errch chan error) {
	for {
		var msg *pb.Message

		err := s.srv.RecvMsg(msg)

		if err != nil {
			// exit
			errch <- err
			break
		} else {

		}

		// Load the the response channel and delete its map key reference
		if val, ok := s.responseQueue.LoadAndDelete(msg.GetId()); ok {
			// For now assume this is a trigger response...
			response := msg.GetTriggerResponse()
			rChan := val.(chan *pb.TriggerResponse)
			// Write the response the the waiting recipient
			rChan <- response
		} else {
			errch <- fmt.Errorf("Fatal: FaaS server in base state exiting!!!")
		}
	}
}

// Start the stream
func (s *FaasServer) TriggerStream(srv pb.Faas_TriggerStreamServer) error {
	s.srv = srv

	errch := make(chan error)
	go s.recieveMessages(errch)

	return <-errch
}
