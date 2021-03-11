package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nitric-dev/membrane/sources"
)

type SourceHandler interface {
	HandleEvent(source *sources.Event) error
	HandleHttpRequest(source *sources.HttpRequest) *http.Response
}

// UnimplementedSourceHandler
type UnimplementedSourceHandler struct{}

func (*UnimplementedSourceHandler) HandleEvent(source *sources.Event) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

func (*UnimplementedSourceHandler) HandleHttpRequest(source *sources.HttpRequest) *http.Response {
	return &http.Response{
		Status:     "Internal Server Error",
		StatusCode: 500,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("HTTP Handler Unimplemented"))),
	}
}
