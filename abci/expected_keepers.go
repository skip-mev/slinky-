package abci

import (
	slpptypes "github.com/CosmWasm/wasmd/x/slpp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SLPPKeeper represents the expected interface for the slpp keeper
type SLPPKeeper interface {
	// GetAVSPerID returns the AVS for a given ID
	GetAVSPerID(ctx sdk.Context, id uint64) slpptypes.AVS
	// GetAllAVSes returns all AVSes
	GetAllAVSes(ctx sdk.Context) ([]slpptypes.AVS, error)
}
