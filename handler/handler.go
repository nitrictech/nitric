package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nitric-dev/membrane/triggers"
)

type TriggerHandler interface {
	HandleEvent(trigger *triggers.Event) error
	HandleHttpRequest(trigger *triggers.HttpRequest) (*triggers.HttpResponse, error)
}

type UnimplementedTriggerHandler struct{}

func (*UnimplementedTriggerHandler) HandleEvent(trigger *triggers.Event) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

func (*UnimplementedTriggerHandler) HandleHttpRequest(trigger *triggers.HttpRequest) *http.Response {
	return &http.Response{
		Status:     "Unimplemented",
		StatusCode: 501,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("HTTP Handler Unimplemented"))),
	}
}
