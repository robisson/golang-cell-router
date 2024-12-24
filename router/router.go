package router

import (
	"cell-router/config"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Router struct {
	config *config.Config
}

func NewRouter() *Router {
	cfg := config.LoadConfig()
	log.Println("Configuration loaded successfully")
	return &Router{config: cfg}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	clientIDStr := req.URL.Query().Get("client_id")
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		log.Printf("Invalid client_id: %s", clientIDStr)
		http.Error(w, "Invalid client_id", http.StatusBadRequest)
		return
	}

	log.Printf("Received request for client_id: %d", clientID)

	for _, cell := range r.config.Cells {
		if clientID >= cell.RangeFrom && clientID <= cell.RangeTo {
			log.Printf("Routing to cell: %s, endpoint: %s", cell.Name, cell.Endpoint)

			proxyURL, err := url.Parse(cell.Endpoint)
			if err != nil {
				log.Printf("Error parsing endpoint URL: %v", err)
				http.Error(w, "Invalid endpoint URL", http.StatusInternalServerError)
				return
			}

			proxyReq, err := http.NewRequest(req.Method, proxyURL.String(), req.Body)
			if err != nil {
				log.Printf("Error creating proxy request: %v", err)
				http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
				return
			}
			proxyReq.Header = req.Header

			client := &http.Client{}
			resp, err := client.Do(proxyReq)
			if err != nil {
				log.Printf("Error making request to cell: %v", err)
				http.Error(w, "Error making request to cell", http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			log.Printf("Received response from cell: %s, status: %d", cell.Name, resp.StatusCode)

			for key, values := range resp.Header {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
			return
		}
	}

	log.Printf("No matching cell found for client_id: %d", clientID)
	http.Error(w, "No matching cell found", http.StatusNotFound)
}
