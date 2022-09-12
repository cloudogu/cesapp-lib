package registry

import (
	"strings"
)

type Node struct {
	subnodes []*Node
	isDir    bool
	fullKey  string
	value    string
	parent   *Node
}

func (drn *Node) IsDir() bool {
	return drn.isDir
}

func (drn *Node) HasSubNodes() bool {
	return len(drn.GetSubNodes()) > 0
}

func (drn *Node) GetSubNodes() []*Node {
	return drn.subnodes
}

func (drn *Node) GetFullKey() string {
	return drn.fullKey
}

func (drn *Node) GetKey() string {
	if strings.Contains(drn.fullKey, "/") {
		splitted := strings.Split(drn.fullKey, "/")
		return splitted[len(splitted)-1]
	}

	return drn.fullKey
}

func (drn *Node) GetValue() string {
	return drn.value
}

func (drn *Node) GetParent() *Node {
	return drn.parent
}

func (drn *Node) GetSubNode(key string) *Node {
	for _, subnode := range drn.subnodes {
		if subnode.GetKey() == key {
			return subnode
		}
	}

	return nil
}
