package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	cellName := os.Getenv("CELL_NAME")
	log.Printf("Received request at %s", cellName)
	fmt.Fprintf(w, "Hello from %s!", cellName)
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Starting mock cell server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
