package provider

import "merkleapi/types"

type StateProvider interface {
	GetState() (*types.State, error)
	SetState(state *types.State) error
}
