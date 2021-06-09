package worker

import (
	"fmt"
	"sync"

	pb "github.com/nitric-dev/membrane/interfaces/nitric/v1"
	"github.com/nitric-dev/membrane/triggers"
)

// FaasWorker
// Worker representation for a Nitric FaaS functon
type FaasWorker struct {
	// gRPC Stream for this worker
	stream pb.Faas_TriggerStreamServer
	// Response channels for this worker
	responseQueue sync.Map	
}

func (s *FaasWorker) HandleHttpRequest(trigger *triggers.HttpRequest) (*triggers.HttpResponse, error) {
	// Generate an ID here
	ID := uuid.New().String()

	triggerRequest := &pb.TriggerRequest{
		Data: trigger.Body,
		Context: &pb.TriggerRequest_Http{
			Http: &pb.HttpTriggerContext{
				Method: trigger.Method,
				QueryParams: trigger.Query,
				// TODO: Populate path params
				PathParams: make(map[string]string),
				// TODO: Update contract to provide original path as well???
			}
		},
	}

	// construct the message
	message := &pb.Message{
		Id: ID,
		Content: &pb.Message_TriggerRequest{
			TriggerRequest: triggerRequest,
		},
	}

	// send the message
	error := s.stream.Send(message)

	if err != nil {
		// There was an error enqueuing the message
		return nil, err
	}

	// Get a lock on the response queue
	returnChan := make(chan *pb.TriggerResponse)

	// Let the worker know where to return the results
	s.responseQueue.Store(ID, returnChan)

	// wait for the response
	response := <-returnChan

	// translate the response to a Http response trigger

	return returnChan, nil
}

func (s *FaasWorker) HandleEvent(trigger *triggers.Event) error {
	// Generate an ID here
	ID := uuid.New().String()
	triggerRequest := &pb.TriggerRequest{
		Data: trigger.Payload,
		Context: &pb.TriggerRequest_Topic{
			Topic: &pb.TopicTriggerContext{
				Topic: trigger.Topic,
				// FIXME: Add missing fields here...
			},
		},
	}

	// construct the message
	message := &pb.Message{
		Id: ID,
		Content: &pb.Message_TriggerRequest{
			TriggerRequest: triggerRequest,
		},
	}

	// send the message
	error := s.stream.Send(message)

	if err != nil {
		// There was an error enqueuing the message
		return nil, err
	}

	// Get a lock on the response queue
	returnChan := make(chan *pb.TriggerResponse)

	// Let the worker know where to return the results
	s.responseQueue.Store(ID, returnChan)

	// wait for the response
	// FIXME: Need to handle timeouts here...
	response := <-returnChan

	topic := response.GetTopic()

	if topic == nil {
		// Fatal error in this case
		// We don't have the correct response type for this handler
		return fmt.Errorf("Fatal: Error handling event, incorrect response recieved from function")
	}

	if topic.GetSuccess() {
		return nil
	}

	return fmt.Errorf("Error ocurred handling the event")
}

// listen
func (s *FaasWorker) listen(chan error) {
	// Listen for responses
	for {
		var msg *pb.Message

		// Blocking read here...
		err := s.srv.RecvMsg(msg)

		if err != nil {
			// exit
			errch <- err
			break
		}

		// Load the the response channel and delete its map key reference
		if val, ok := s.responseQueue.LoadAndDelete(msg.GetId()); ok {
			// For now assume this is a trigger response...
			response := msg.GetTriggerResponse()
			rChan := val.(chan *pb.TriggerResponse)
			// Write the response the the waiting recipient
			rChan <- response
		} else {
			errch <- fmt.Errorf("Fatal: FaaS Worker in bad state exiting!!!")
			break
		}
	}
}

// Package private method
// Only a pool may create a new faas worker
func newFaasWorker(stream pb.Faas_TriggerStreamServer) TriggerHandler {
	return &FaasWorker{
		stream: stream
	}
}