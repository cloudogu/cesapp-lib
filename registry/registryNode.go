package registry

import (
	"strings"
)

type RegistryNode interface {
	IsDir() bool
	HasSubNodes() bool
	GetSubNodes() []RegistryNode
	GetFullKey() string
	GetKey() string
	GetValue() string
	GetParent() RegistryNode
	GetSubNode(key string) RegistryNode
}

type defaultRegistryNode struct {
	subnodes []RegistryNode
	isDir    bool
	fullKey  string
	value    string
	parent   *defaultRegistryNode
}

func (drn *defaultRegistryNode) IsDir() bool {
	return drn.isDir
}

func (drn *defaultRegistryNode) HasSubNodes() bool {
	return len(drn.GetSubNodes()) > 0
}

func (drn *defaultRegistryNode) GetSubNodes() []RegistryNode {
	return drn.subnodes
}

func (drn *defaultRegistryNode) GetFullKey() string {
	return drn.fullKey
}

func (drn *defaultRegistryNode) GetKey() string {
	if strings.Contains(drn.fullKey, "/") {
		splitted := strings.Split(drn.fullKey, "/")
		return splitted[len(splitted)-1]
	}

	return drn.fullKey
}

func (drn *defaultRegistryNode) GetValue() string {
	return drn.value
}

func (drn *defaultRegistryNode) GetParent() RegistryNode {
	return drn.parent
}

func (drn *defaultRegistryNode) GetSubNode(key string) RegistryNode {
	for _, subnode := range drn.subnodes {
		if subnode.GetKey() == key {
			return subnode
		}
	}

	return nil
}
