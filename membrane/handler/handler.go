package handler

import (
	"fmt"
	"net/http"

	"github.com/nitric-dev/membrane/plugins/sdk/sources"
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
		Body: ioutils.,
	}
}
