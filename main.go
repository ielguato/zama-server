package main

import (
	"log"
	"net/http"
	"zama-server/handlers"

	"github.com/gorilla/mux"
)

const (
	UploadsDirectory = "uploads" // Directory where files will be stored
	MaxUploadSize    = 10 << 20  // Max file upload size (10 MB)
)

func main() {
	r := mux.NewRouter()
	// Define route handlers for different endpoints and HTTP methods
	r.HandleFunc("/upload/{filename}", handlers.UploadFile).Methods("POST")
	r.HandleFunc("/download/{filename}/{segmentName}", handlers.DownloadFile).Methods("GET")
	r.HandleFunc("/delete/{filename}", handlers.DeleteFile).Methods("DELETE")
	r.HandleFunc("/requestProof/{filename}/{segmentname}", handlers.RequestProof).Methods("GET")
	r.HandleFunc("/list", handlers.ListFiles).Methods("GET")

	// Start the HTTP server on port 8080
	log.Fatal(http.ListenAndServe(":8080", r))
}
