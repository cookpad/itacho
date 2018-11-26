package server

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/handlers"

	"github.com/cookpad/itacho/storage"
)

// Opts for server options
type Opts struct {
	BindPort              uint
	BindAddr              string
	ObjectStorageEndpoint string
}

// Start new xDS API and admin HTTP server
func Start(opts Opts) error {
	endpoint := opts.ObjectStorageEndpoint
	port := opts.BindPort
	addr := opts.BindAddr
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("Failed to parse URL: %s : %s", endpoint, err)
	}
	str := storage.NewXdsResponseStorage(*endpointURL)

	h := http.NewServeMux()
	h.Handle("/hc", handlers.CombinedLoggingHandler(os.Stderr, &healthcheckHandler{}))
	h.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, NewXdsHandler(str)))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), h); err != nil {
		return fmt.Errorf("Failed to start server: %s", err)
	}

	return nil
}

type healthcheckHandler struct {
}

func (h *healthcheckHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
}
