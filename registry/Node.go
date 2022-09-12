package registry

import (
	"strings"
)

type Node struct {
	SubNodes []Node
	IsDir    bool
	FullKey  string
	Value    string
}

func (drn Node) HasSubNodes() bool {
	return len(drn.SubNodes) > 0
}

func (drn Node) Key() string {
	if strings.Contains(drn.FullKey, "/") {
		splitted := strings.Split(drn.FullKey, "/")
		return splitted[len(splitted)-1]
	}

	return drn.FullKey
}

func (drn Node) SubNodeByName(key string) Node {
	for _, subnode := range drn.SubNodes {
		if subnode.Key() == key {
			return subnode
		}
	}

	return Node{}
}
