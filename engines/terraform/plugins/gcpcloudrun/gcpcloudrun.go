package gcpcloudrun

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/nitrictech/nitric/server/runtime/service"
)

type gcpcloudappService struct {
	proxy service.Proxy
}

func (a *gcpcloudappService) Start(proxy service.Proxy) error {
	fmt.Println("Starting Cloud Run service proxy")

	p := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			fmt.Println("Directing request:", req.URL.Path)
			req.URL.Host = proxy.Host()
			req.URL.Scheme = "http"
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", p.ServeHTTP)

	fmt.Println("Starting Cloud Run service proxy on port")
	return http.ListenAndServe(fmt.Sprintf(":8080"), mux)
}

func Plugin() (service.Service, error) {
	return &gcpcloudappService{}, nil
}
