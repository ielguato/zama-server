package merkle

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// MerkleNode represents a node in the Merkle tree.
type MerkleNode struct {
	Hash []byte `json:"hash"`
}

// MerkleTree represents a Merkle tree.
type MerkleTree struct {
	MerkleNodes [][]*MerkleNode `json:"merkleNodes"`
	IsTreeBuilt bool            `json:"isTreeBuilt"`
}

type ProofPart struct {
	Hash    []byte
	IsRight bool
}

// NewMerkleNode creates a new MerkleNode.
func newMerkleNode(hash []byte) *MerkleNode {
	node := MerkleNode{
		Hash: hash,
	}
	return &node
}

// AddLeaf adds a new leaf node to the MerkleTree.
func (mt *MerkleTree) AddLeaf(data []byte) {
	hash := sha256.Sum256(data)
	leaf := newMerkleNode(hash[:])

	// Check if the level slice already exists for leaves (level 0)
	if len(mt.MerkleNodes) == 0 {
		mt.MerkleNodes = append(mt.MerkleNodes, []*MerkleNode{})
	}

	// Add the leaf to the leaves slice (level 0)
	mt.MerkleNodes[0] = append(mt.MerkleNodes[0], leaf)
}

// GetRootHash returns the root hash of the MerkleTree.
func (mt *MerkleTree) GetRootHash() []byte {
	// The root hash is in the last level of the Merkle tree
	return mt.MerkleNodes[len(mt.MerkleNodes)-1][0].Hash
}

func (mt *MerkleTree) BuildMerkleTree() error {
	if len(mt.MerkleNodes) == 0 || len(mt.MerkleNodes[0]) == 0 {
		return fmt.Errorf("No leaves provided")
	}

	// Initialize the tree with the level containing the leaves
	newMerkleNodes := make([][]*MerkleNode, 1)
	newMerkleNodes[0] = mt.MerkleNodes[0]
	mt.MerkleNodes = newMerkleNodes

	for level := 0; len(mt.MerkleNodes[level]) > 1; level++ {
		currentLevel := mt.MerkleNodes[level]

		var newLevel []*MerkleNode
		for i := 0; i < len(currentLevel); i += 2 {
			left := currentLevel[i]
			var right *MerkleNode
			if i+1 < len(currentLevel) {
				right = currentLevel[i+1]
			}

			// Combine hashes if necessary
			combined := left.Hash
			if right != nil {
				combined = append(left.Hash, right.Hash...)
			} else {
				combined = append(left.Hash, left.Hash...)
			}

			// Calculate the combined hash of left and right children
			hash := sha256.Sum256(combined)

			// Create a new node with the hash
			newNode := &MerkleNode{Hash: hash[:]}
			newLevel = append(newLevel, newNode)
		}
		mt.MerkleNodes = append(mt.MerkleNodes, newLevel)
	}

	mt.IsTreeBuilt = true // Set the flag to true when the tree is built

	return nil
}

func (mt *MerkleTree) GenerateMerkleProofForSegment(segmentIndex int) ([]ProofPart, error) {

	// Check if Merkle tree is not empty and is properly constructed
	if mt == nil || len(mt.MerkleNodes) == 0 || len(mt.MerkleNodes[0]) == 0 {
		return nil, fmt.Errorf("Merkle tree is not properly constructed")
	}

	// Validate the segment index is within the bounds of the leaf nodes
	if segmentIndex < 0 || segmentIndex >= len(mt.MerkleNodes[0]) {
		return nil, fmt.Errorf("Invalid segment index: index out of bounds")
	}
	proof := []ProofPart{}

	// Include the leaf hash itself (although typically in Merkle Proofs, the leaf itself isn't included as proof part)
	leafHash := mt.MerkleNodes[0][segmentIndex].Hash
	proof = append(proof, ProofPart{Hash: leafHash, IsRight: false})

	for level := 0; level < len(mt.MerkleNodes)-1; level++ {
		nodesAtLevel := mt.MerkleNodes[level]
		isRight := false
		var siblingIndex int

		// Check if the current index is even or odd to determine if the proof part should be a left or right sibling
		if segmentIndex%2 == 0 {
			siblingIndex = segmentIndex + 1
			isRight = true // If our node is the left child, then the sibling (proof part) is the right child
		} else {
			siblingIndex = segmentIndex - 1
			// isRight is already set to false
		}

		// Check if the siblingIndex is within the bounds of the tree at this level
		if siblingIndex < len(nodesAtLevel) {
			siblingHash := nodesAtLevel[siblingIndex].Hash
			proof = append(proof, ProofPart{Hash: siblingHash, IsRight: isRight})
		} else if segmentIndex == len(nodesAtLevel)-1 {
			// This is the case where we have an odd number of nodes at this level,
			// and the current node is the last one (without a sibling),
			// we will need to duplicate it in the proof.
			proof = append(proof, ProofPart{Hash: nodesAtLevel[segmentIndex].Hash, IsRight: false})
		}

		// Move to the parent index for the next level of the tree
		segmentIndex /= 2
	}

	return proof, nil
}

func VerifyMerkleProof(rootHash []byte, proof []ProofPart) (bool, error) {
	if len(proof) == 0 || len(proof[0].Hash) != sha256.Size {
		return false, fmt.Errorf("Invalid proof")
	}

	currentHash := proof[0].Hash

	for _, part := range proof[1:] {
		if len(part.Hash) != sha256.Size {
			return false, fmt.Errorf("Invalid hash size in proof part")
		}

		var combinedHash [sha256.Size]byte
		if part.IsRight {
			// If the proof part is the right node, append it to the right of the current hash.
			combinedHash = sha256.Sum256(append(currentHash, part.Hash...))
		} else {
			// If the proof part is the left node, append it to the left of the current hash.
			combinedHash = sha256.Sum256(append(part.Hash, currentHash...))
		}

		currentHash = combinedHash[:]
	}

	return bytes.Equal(currentHash, rootHash), nil
}

// MerkleNode represents a node i// SaveMerkleTree saves a Merkle tree to a JSON file using the path.
func (mt *MerkleTree) SaveMerkleTree(filePath string) error {
	// Serialize the MerkleTree to JSON
	data, err := json.Marshal(mt)
	if err != nil {
		return err
	}

	// Write the JSON data to a file using the path
	err = ioutil.WriteFile(filePath, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// LoadMerkleTree loads a Merkle tree from a JSON file using the path.
func (mt *MerkleTree) LoadMerkleTree(filePath string) error {
	// Read the JSON data from the file using the path
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Deserialize the JSON data into the MerkleTree
	err = json.Unmarshal(data, mt)
	if err != nil {
		return err
	}

	return nil
}
