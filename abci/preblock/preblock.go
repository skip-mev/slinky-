package preblock

import (
	"sort"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	"github.com/CosmWasm/wasmd/abci"
	ve "github.com/CosmWasm/wasmd/x/slpp/vote_extensions"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
)

type VEWithVotingPower struct {
	ve.VoteExtension
	VotingPower uint64
}

// PreBlockHandler is a function that is called before each block is processed.
func PreBlocker(k abci.SLPPKeeper, wk *wasmkeeper.Keeper) sdk.PreBlocker {
	return func(ctx sdk.Context, req *cmtabci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
		// skip if ves are not enabled 
		if !abci.VoteExtensionsEnabled(ctx) {
			return &sdk.ResponsePreBlock{}, nil
		}
		
		// expect the extended-commit to be in the block's txs
		if len(req.Txs) < abci.InjectedTxs {
			ctx.Logger().Error(
				"no extended-commit found in the block",
			)
			return nil, fmt.Errorf("no extended-commit found in the block")
		}

		// decode the extended-commit
		var extCommitinfo cmtabci.ExtendedCommitInfo
		if err := extCommitinfo.Unmarshal(req.Txs[abci.ExtCommitInfoIdx]); err != nil {
			ctx.Logger().Error(
				"failed to unmarshal extended-commit",
				"err", err,
			)
			return nil, err
		}

		ctx.Logger().Info(
			"extended-commit found in the block",
			"votes", len(extCommitinfo.Votes),
			"round", extCommitinfo.Round,
			"height", ctx.BlockHeight(),
		)

		// extended-commit is valid
		ves, err := decodeVoteExtensions(extCommitinfo)
		if err != nil {
			return nil, err
		}

		// generate the messages in the order they should be executed
		msgs, err := generateMsgsFromVoteExtensions(ctx, k, ves)
		if err != nil {
			return nil, err
		}
		// sort the messages so they are determi nistically executed
		sort.Sort(sortedExecuteMsgs(msgs))

		// iterate over all execute messages, and execute them
		mk := wasmkeeper.NewMsgServerImpl(wk)
		for _, msg := range msgs {
			ctx.Logger().Info(
				"executing message",
				"msg", msg,
			)
			// cast to 
			res, err := mk.SudoContract(ctx, msg.(*wasmtypes.MsgSudoContract))
			if err != nil {
				ctx.Logger().Error(
					"failed to execute message",
					"msg", msg,
					"err", err,
				)		
				return nil, err
			}

			// log the result
			ctx.Logger().Info(
				"executed message",
				"msg", msg,
				"res", res,
			)
		}

		return &sdk.ResponsePreBlock{}, nil
	}
}

// decodeVoteExtensions decodes the vote extensions from the given commit
func decodeVoteExtensions(llc cmtabci.ExtendedCommitInfo) ([]VEWithVotingPower, error) {
	ves := make([]VEWithVotingPower, len(llc.Votes))

	for i, vote := range llc.Votes {
		// decode the vote extension
		var ve ve.VoteExtension
		if err := ve.Unmarshal(vote.VoteExtension); err != nil {
			return nil, err
		}

		ves[i] = VEWithVotingPower{
			VoteExtension: ve,
			VotingPower: uint64(vote.Validator.Power),
		}
	}

	return ves, nil
}

// generateMsgsFromVoteExtensions generates the messages from the given vote extensions.
// The msgs generated are SudoContractExecuteMsgs. Notice, the order in which the messages
// are returned is non-deterministic.
func generateMsgsFromVoteExtensions(ctx sdk.Context, k abci.SLPPKeeper, ves []VEWithVotingPower) ([]sdk.Msg, error) {
	avsDataPerID := make(map[uint64][]abci.DataWithVotingPower)

	// iterate over the vote extensions
	for _, ve := range ves {
		// per avs-id referenced in a vote-extension
		for avsId, avsData := range ve.AvsData {
			if _, ok := avsDataPerID[avsId]; !ok {
				avsDataPerID[avsId] = make([]abci.DataWithVotingPower, 0)
			}

			avsDataPerID[avsId] = append(avsDataPerID[avsId], abci.DataWithVotingPower{
				Data: avsData,
				VotingPower: ve.VotingPower,
			})
		}
	}

	// generate the messages
	msgs := make([]sdk.Msg, 0, len(avsDataPerID))
	for id, avsDataToAggregate := range avsDataPerID {
		avs, ok := k.GetAVSByID(ctx, id)
		if !ok {
			return nil, fmt.Errorf("avs with id %d not found", id)
		}

		payload, err := json.Marshal(abci.AggregationContractPayload{
			Data: avsDataToAggregate,
		})
		if err != nil {
			return nil, err
		}
		
		msgs = append(msgs, &wasmtypes.MsgSudoContract{
			Contract: avs.ContractAddress,
			Msg: payload,
		})
	}

	return msgs, nil
}

type sortedExecuteMsgs []sdk.Msg

func (s sortedExecuteMsgs) Len() int {
	return len(s)
}

func (s sortedExecuteMsgs) Less(i, j int) bool {
	msgi := s[i].(*wasmtypes.MsgSudoContract)
	msgj := s[j].(*wasmtypes.MsgSudoContract)

	return msgi.Contract < msgj.Contract
}

func (s sortedExecuteMsgs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
