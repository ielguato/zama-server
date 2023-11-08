package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ListFiles responds with a list of directories inside the UploadsDirectory.
func ListFiles(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests for this handler
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	files, err := ioutil.ReadDir(UploadsDirectory)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading directory: %v", err), http.StatusInternalServerError)
		return
	}

	directories := make([]string, 0)
	for _, f := range files {
		if f.IsDir() {
			directories = append(directories, f.Name())
		}
	}

	// Set the content type as application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the directories to JSON and send the response
	if err := json.NewEncoder(w).Encode(directories); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}
