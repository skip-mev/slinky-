package proposals

import (
	"fmt"

	"github.com/CosmWasm/wasmd/abci"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ProposalHandler handles Prepare / ProcessProposal invocations.
type ProposalHandler struct {
	vs baseapp.ValidatorStore
}

func (ph ProposalHandler) PrepareProposalHandler() sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req *cmtabci.RequestPrepareProposal) (*cmtabci.ResponsePrepareProposal, error) {
		// skip logic if ves are not enabled
		if !abci.VoteExtensionsEnabled(ctx) {
			ctx.Logger().Info(
				"vote extensions are not enabled",
			)

			return &cmtabci.ResponsePrepareProposal{
				Txs: [][]byte{},
			}, nil
		}

		// validate the extended commit
		if err := baseapp.ValidateVoteExtensions(ctx, ph.vs, req.Height, ctx.ChainID(), req.LocalLastCommit); err != nil {
			ctx.Logger().Error(
				"failed to validate extended-commit",
				"err", err,
			)

			return &cmtabci.ResponsePrepareProposal{
				Txs: [][]byte{},
			}, err
		}

		// marshal
		extCommitInfoBz, err := req.LocalLastCommit.Marshal()
		if err != nil {
			ctx.Logger().Error(
				"failed to marshal extended-commit",
				"err", err,
			)

			return &cmtabci.ResponsePrepareProposal{
				Txs: [][]byte{},
			}, err
		}

		ctx.Logger().Info(
			"extended-commit prepared for proposal",
			"height", req.Height,
		)

		return &cmtabci.ResponsePrepareProposal{
			Txs: [][]byte{extCommitInfoBz},
		}, nil
	}
}

func (ph ProposalHandler) ProcessProposalHandler() sdk.ProcessProposalHandler {
	return func(ctx sdk.Context, req *cmtabci.RequestProcessProposal) (*cmtabci.ResponseProcessProposal, error) {
		// skip logic if ves are not enabled
		if !abci.VoteExtensionsEnabled(ctx) {
			ctx.Logger().Info(
				"vote extensions are not enabled",
			)

			return &cmtabci.ResponseProcessProposal{
				Status: cmtabci.ResponseProcessProposal_ACCEPT,
			}, nil
		}

		// check that the extended commit is in state
		if len(req.Txs) <= abci.InjectedTxs {
			ctx.Logger().Error(
				"no extended-commit found in the block",
			)

			return &cmtabci.ResponseProcessProposal{
				Status: cmtabci.ResponseProcessProposal_REJECT,
			}, fmt.Errorf("no extended-commit found in the block")
		}

		// decode the extended commit
		var extCommitInfo cmtabci.ExtendedCommitInfo
		if err := extCommitInfo.Unmarshal(req.Txs[abci.ExtCommitInfoIdx]); err != nil {
			ctx.Logger().Error(
				"failed to unmarshal extended-commit",
				"err", err,
			)

			return &cmtabci.ResponseProcessProposal{
				Status: cmtabci.ResponseProcessProposal_REJECT,
			}, err
		}

		// validate the extended commit
		if err := baseapp.ValidateVoteExtensions(ctx, ph.vs, req.Height, ctx.ChainID(), extCommitInfo); err != nil {
			ctx.Logger().Error(
				"failed to validate extended-commit",
				"err", err,
			)

			return &cmtabci.ResponseProcessProposal{
				Status: cmtabci.ResponseProcessProposal_REJECT,
			}, err
		}

		ctx.Logger().Info(
			"extended-commit processed for proposal",
			"height", req.Height,
		)
		return &cmtabci.ResponseProcessProposal{
			Status: cmtabci.ResponseProcessProposal_ACCEPT,
		}, nil
	}
}
