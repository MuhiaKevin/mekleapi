package types

import "merkleapi/address"

type State struct {
	LastIndex uint64
	Version   int
	Root      Node
	Address   *address.Address
}
