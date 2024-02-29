package voteextensions

import (
	"context"

	"github.com/CosmWasm/wasmd/abci"
	"github.com/CosmWasm/wasmd/x/slpp/service"
	ve "github.com/CosmWasm/wasmd/x/slpp/vote_extensions"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MutltiOracleClient wraps multiple OracleClients, one per AVS-ID in state
type MultiOracleClient interface {
	VoteExtensionData(ctx context.Context, avsID uint64, req *service.VoteExtensionDataRequest) (*service.VoteExtensionDataResponse, error)
}

// NewExtendVoteHandler returns a handler for ExtendVote. This ExtendVoteHandler wraps an
// OracleClient
func NewExtendVoteHandler(oc MultiOracleClient, k abci.SLPPKeeper) sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, rev *cmtabci.RequestExtendVote) (*cmtabci.ResponseExtendVote, error) {
		ext := ve.VoteExtension{
			AvsData: make(map[uint64][]byte),
		}

		avses, err := k.GetAllAVSIDs(ctx)
		if err != nil {
			ctx.Logger().Error(
				"failed to get all AVSes",
				"err", err,
			)
			return &cmtabci.ResponseExtendVote{
				VoteExtension: []byte{},
			}, nil
		}

		for _, avs := range avses {
			veData, err := oc.VoteExtensionData(ctx, avs, &service.VoteExtensionDataRequest{})
			if err != nil {
				ctx.Logger().Error(
					"failed to get vote extension data",
					"err", err,
				)
			}

			ext.AvsData[avs] = veData.Data
		}

		extBytes, err := ext.Marshal()
		if err != nil {
			ctx.Logger().Error(
				"failed to marshal vote extension",
				"err", err,
			)
			return &cmtabci.ResponseExtendVote{
				VoteExtension: []byte{},
			}, nil
		}

		return &cmtabci.ResponseExtendVote{
			VoteExtension: extBytes,
		}, nil
	}
}

// NewVerifyVoteExtensionHandler is a no-op
func NewVerifyVoteExtensionHandler() sdk.VerifyVoteExtensionHandler {
	return func(ctx sdk.Context, rvve *cmtabci.RequestVerifyVoteExtension) (*cmtabci.ResponseVerifyVoteExtension, error) {
		return &cmtabci.ResponseVerifyVoteExtension{
			Status: cmtabci.ResponseVerifyVoteExtension_ACCEPT,
		}, nil
	}
}
