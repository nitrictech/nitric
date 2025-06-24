package awslambda

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/nitrictech/nitric/runtime/service"
)

type awsfargateService struct{}

func (a *awsfargateService) Start(proxy service.Proxy) error {
	fmt.Println("Starting Fargate service proxy")
	// get the container port from the environmentW
	containerPort := os.Getenv("CONTAINER_PORT")
	if containerPort == "" {
		return fmt.Errorf("CONTAINER_PORT is not set")
	}

	p := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			fmt.Println("Directing request:", req.URL.Path)
			// TODO: Do additional analysis of the request in order to perform event subscription routing
			req.URL.Host = proxy.Host()
			req.URL.Scheme = "http"
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", p.ServeHTTP)
	mux.HandleFunc("/x-nitric-health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Recieved Health check")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	fmt.Println("Starting Fargate service proxy on port", containerPort)
	return http.ListenAndServe(fmt.Sprintf(":%s", containerPort), mux)
}

func Plugin() (service.Service, error) {
	return &awsfargateService{}, nil
}
