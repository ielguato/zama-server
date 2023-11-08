package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"zama-server/merkle"

	"github.com/gorilla/mux"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form in the request
	r.ParseMultipartForm(MaxUploadSize)
	vars := mux.Vars(r)
	fileName := vars["filename"]

	// Ensure the uploads directory exists
	if err := os.MkdirAll(UploadsDirectory, os.ModePerm); err != nil {
		http.Error(w, "Failed to create uploads directory", http.StatusInternalServerError)
		return
	}

	// Construct the file paths for the Merkle tree and the segment
	merkleTreePath := filepath.Join(UploadsDirectory, fileName, MerkleTreeFileName)
	segmentPath := filepath.Join(UploadsDirectory, fileName)

	// Initialize or load the Merkle tree
	merkleTree := &merkle.MerkleTree{}
	if _, err := os.Stat(merkleTreePath); err == nil {
		if err := merkleTree.LoadMerkleTree(merkleTreePath); err != nil {
			http.Error(w, "Error loading Merkle tree", http.StatusInternalServerError)
			return
		}
	}

	// Retrieve the file from the form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	segmentPath = filepath.Join(segmentPath, handler.Filename)
	// Create the directory for the segment if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(segmentPath), os.ModePerm); err != nil {
		http.Error(w, "Failed to create segment directory", http.StatusInternalServerError)
		return
	}

	// Check if the segment already exists
	if _, err := os.Stat(segmentPath); err == nil {
		// Directory already exists, return an error
		http.Error(w, "Segment already exists, rename the file you want to upload or delete it from the server before", http.StatusBadRequest)
		return
	}
	// Create the segment file and copy the uploaded content into it
	newFile, err := os.Create(segmentPath)
	if err != nil {
		http.Error(w, "Failed to create new file", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	if _, err = io.Copy(newFile, file); err != nil {
		http.Error(w, "Failed to copy file content", http.StatusInternalServerError)
		return
	}

	// Update the Merkle tree with the new file segment
	data, err := ioutil.ReadFile(segmentPath)
	if err != nil {
		http.Error(w, "Error reading segment data", http.StatusInternalServerError)
		return
	}
	merkleTree.AddLeaf(data)

	// Save the updated Merkle tree to disk
	if err := merkleTree.SaveMerkleTree(merkleTreePath); err != nil {
		http.Error(w, "Error saving Merkle tree", http.StatusInternalServerError)
		return
	}

	// Notify the client that the file was uploaded successfully
	fmt.Fprintf(w, "File uploaded successfully: %s\n", segmentPath)
}
