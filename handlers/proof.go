package handlers

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"zama-server/merkle"

	"github.com/gorilla/mux"
)

func RequestProof(w http.ResponseWriter, r *http.Request) {
	// Extract the filename and segment index from the URL
	vars := mux.Vars(r)
	filename := vars["filename"]
	segmentName := vars["segmentname"]
	segmentIndex, err := strconv.Atoi(segmentName)
	if err != nil {
		http.Error(w, "Invalid segment index", http.StatusBadRequest)
		return
	}

	// Load the Merkle tree from the associated JSON file
	merkleTree := &merkle.MerkleTree{}
	merkleTreePath := filepath.Join(UploadsDirectory, filename, MerkleTreeFileName)
	if err := merkleTree.LoadMerkleTree(merkleTreePath); err != nil {
		http.Error(w, "Error loading Merkle tree", http.StatusInternalServerError)
		return
	}

	// Ensure the Merkle tree is built
	if !merkleTree.IsTreeBuilt {
		merkleTree.BuildMerkleTree()
	}

	// Generate a proof for the specified segment
	proof, err := merkleTree.GenerateMerkleProofForSegment(segmentIndex)
	if err != nil {
		http.Error(w, "Error generating proof", http.StatusInternalServerError)
		return
	}

	// Verify the proof (not strictly necessary to respond to the request)
	if verified, err := merkle.VerifyMerkleProof(merkleTree.GetRootHash(), proof); err != nil || !verified {
		http.Error(w, "Proof verification failed", http.StatusInternalServerError)
		return
	}

	// Marshal the proof into JSON for the response
	proofBytes, err := json.Marshal(proof)
	if err != nil {
		http.Error(w, "Failed to serialize proof", http.StatusInternalServerError)
		return
	}

	// Send the proof to the client
	w.Header().Set("Content-Type", "application/json")
	w.Write(proofBytes)
}
