package gcpcloudrun

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/nitrictech/nitric/server/runtime/service"
)

type gcpcloudappService struct {
	proxy service.Proxy
}

func (a *gcpcloudappService) Start(proxy service.Proxy) error {
	mux := http.NewServeMux()

	a.proxy = proxy
	mux.HandleFunc("/", a.handler)

	err := http.ListenAndServe(":"+os.Getenv("INGRESS_PORT"), mux)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (a *gcpcloudappService) handler(w http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	resp, err := a.proxy.Forward(ctx, request)
	if err != nil {
		http.Error(w, "failed to forward request", http.StatusBadGateway)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to forward request", http.StatusInternalServerError)
		return
	}

	// Translate response headers to a map
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func Plugin() (service.Service, error) {
	return &gcpcloudappService{}, nil
}
