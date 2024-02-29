package keeper

import (
	"context"
	"fmt"

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

// GetAVS returns the AVS with the given id
func (q *queryServer) GetAVS(ctx context.Context, req *types.GetAVSRequest) (*types.GetAVSResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	avs, err := q.keeper.GetAVS(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.GetAVSResponse{
		Avs: avs,
	}, nil
}
