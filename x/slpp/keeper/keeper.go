package keeper

import (
	"context"

	"github.com/CosmWasm/wasmd/x/slpp/types"
)

// SLPPKeeper
type Keeper struct{}

func (k Keeper) RegisterAVS(ctx context.Context, avs *types.MsgRegisterAVS) (uint64, error) {
	panic("unimplemented")
}
