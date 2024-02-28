package keeper

import (
	"context"
	"fmt"

	"github.com/CosmWasm/wasmd/x/slpp/types"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServer returns the default implementation of the x/slpp message service.
func NewMsgServer(k Keeper) types.MsgServer {
	return &msgServer{
		keeper: k,
	}
}

// RegisterAVS does something..
func (m *msgServer) RegisterAVS(ctx context.Context, req *types.MsgRegisterAVS) (*types.MsgRegisterAVSResponse, error) {
	// check the validity of the message
	if req == nil {
		return nil, fmt.Errorf("message cannot be empty")
	}

	id, err := m.keeper.RegisterAVS(ctx, req)
	if err != nil {
		return nil, err
	}

	return &types.MsgRegisterAVSResponse{
		Id: id,
	}, nil
}
