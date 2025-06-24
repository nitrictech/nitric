package gcpcloudrun

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/nitrictech/nitric/runtime/service"
)

type gcpcloudappService struct {
	proxy service.Proxy
}

func (a *gcpcloudappService) Start(proxy service.Proxy) error {
	fmt.Println("Starting Cloud Run service proxy")

	p := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Host = proxy.Host()
			req.URL.Scheme = "http"
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", p.ServeHTTP)

	port := os.Getenv("PORT")
	if port == "" {
		return fmt.Errorf("PORT environment variable not set")
	}

	fmt.Printf("Starting Cloud Run service proxy on port %s\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}

func Plugin() (service.Service, error) {
	return &gcpcloudappService{}, nil
}
