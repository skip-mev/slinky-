package keeper

import (
	"context"

	"github.com/CosmWasm/wasmd/x/slpp/types"
)

// queryServer is the default implementation of the x/slpp QueryService.
type queryServer struct {
	keeper Keeper
}

// NewQueryServer creates a new x/slpp QueryServer
func NewQueryServer(keeper Keeper) types.QueryServer {
	return &queryServer{keeper: keeper}
}

// GetAVS implements types.QueryServer.
func (q *queryServer) GetAVS(context.Context, *types.GetAVSRequest) (*types.GetAVSResponse, error) {
	panic("unimplemented")
}
