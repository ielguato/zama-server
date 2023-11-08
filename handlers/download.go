package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	// Extract file and segment names from the URL
	vars := mux.Vars(r)
	filename := vars["filename"]
	segmentName := vars["segmentName"]
	filePath := filepath.Join(UploadsDirectory, filename, segmentName)

	// Serve the requested file segment for download
	http.ServeFile(w, r, filePath)
}
