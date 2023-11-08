package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	// Extract the directory name from the URL
	vars := mux.Vars(r)
	directoryName := vars["filename"]
	directoryPath := filepath.Join(UploadsDirectory, directoryName)

	// Remove the directory and its contents
	if err := os.RemoveAll(directoryPath); err != nil {
		http.Error(w, "Error deleting the directory", http.StatusInternalServerError)
		return
	}

	// Send a confirmation that the deletion was successful
	w.WriteHeader(http.StatusOK)
}
