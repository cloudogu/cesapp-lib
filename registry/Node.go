package registry

import (
	"strings"
)

// Node Represents the structure (key/value/dir) of the etcd
type Node struct {
	// SubNodes contains all direct children of the current node.
	SubNodes []Node
	// IsDir is true if the current node is a dir, false if the node is a key
	IsDir bool
	// FullKey is the full path of a key including the path of all its parents.
	FullKey string
	// Value is the value of the node - empty string if the node is a directory
	Value string
}

// HasSubNodes returns true if the SubNodes attribute has values, otherwise false
func (drn Node) HasSubNodes() bool {
	return len(drn.SubNodes) > 0
}

// Key Returns the direct key of the node - parent paths are excluded
func (drn Node) Key() string {
	if strings.Contains(drn.FullKey, "/") {
		splitted := strings.Split(drn.FullKey, "/")
		return splitted[len(splitted)-1]
	}

	return drn.FullKey
}

// SubNodeByName Searches in the direct sub nodes for a node with the given key and returns it if found. Returns empty node otherwise.
func (drn Node) SubNodeByName(key string) Node {
	for _, subnode := range drn.SubNodes {
		if subnode.Key() == key {
			return subnode
		}
	}

	return Node{}
}
