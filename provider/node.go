package provider

import (
	"errors"

	"merkleapi/types"
)

type NodeProvider interface {
	GetNode(index uint64, version int) (types.Node, error)
	SetNode(index uint64, version int, node types.Node) error
}

var ErrNodeNotExist = errors.New("node does not exist")
