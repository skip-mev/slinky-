package keeper

import (
	"context"

	"github.com/CosmWasm/wasmd/x/slpp/types"

	"cosmossdk.io/collections"
	"cosmossdk.io/collections/indexes"
	"cosmossdk.io/core/store"
	"encoding/hex"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cometbft/cometbft/crypto"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type avsIndices struct {
	// idUnique is a uniqueness constraint on the IDs of CurrencyPairs. i.e id -> hex-encoded keccak-256 hash of contract bytes -> AVS
	idUnique *indexes.Unique[uint64, string, types.AVS]

	// idMulti is a multi-index on the IDs of CurrencyPairs, i.e. id -> hex-encoded keccak-256 hash of contract bytes -> AVS
	idMulti *indexes.Multi[uint64, string, types.AVS]
}

func (o *avsIndices) IndexesList() []collections.Index[string, types.AVS] {
	return []collections.Index[string, types.AVS]{
		o.idUnique,
		o.idMulti,
	}
}

func newAVSIndices(sb *collections.SchemaBuilder) *avsIndices {
	return &avsIndices{
		idUnique: indexes.NewUnique[uint64, string, types.AVS](
			sb, types.UniqueIndexAVSKeyPrefix, "avs_id_unique_idx", collections.Uint64Key, collections.StringKey,
			func(_ string, avs types.AVS) (uint64, error) {
				return avs.Id, nil
			},
		),
		idMulti: indexes.NewMulti[uint64, string, types.AVS](
			sb, types.IDIndexAVSKeyPrefix, "avs_id_idx", collections.Uint64Key, collections.StringKey,
			func(_ string, avs types.AVS) (uint64, error) {
				return avs.Id, nil
			},
		),
	}
}

type WasmMsgServer interface {
	StoreCode(ctx context.Context, msg *wasmtypes.MsgStoreCode) (*wasmtypes.MsgStoreCodeResponse, error)
	InstantiateContract(ctx context.Context, msg *wasmtypes.MsgInstantiateContract) (*wasmtypes.MsgInstantiateContractResponse, error)
}

// Keeper is the base keeper for the x/oracle module.
type Keeper struct {
	storeService store.KVStoreService
	cdc          codec.BinaryCodec

	wasmMsgServer WasmMsgServer

	// schema
	nextAVSID collections.Sequence
	avsMap    *collections.IndexedMap[string, types.AVS, *avsIndices]
	schema    collections.Schema

	// indexes
	idIndex *indexes.Multi[uint64, string, types.AVS]
}

func NewKeeper(
	ss store.KVStoreService,
	cdc codec.BinaryCodec,
	wasmMsgServer WasmMsgServer,
) Keeper {
	// create a new schema builder
	sb := collections.NewSchemaBuilder(ss)

	indices := newAVSIndices(sb)

	idMulti, ok := indices.IndexesList()[1].(*indexes.Multi[uint64, string, types.AVS])
	if !ok {
		panic("expected idMulti to be a *indexes.Multi[uint64, string, types.AVS]")
	}

	k := Keeper{
		storeService:  ss,
		cdc:           cdc,
		nextAVSID:     collections.NewSequence(sb, types.AVSIDKeyPrefix, "avs_id"),
		avsMap:        collections.NewIndexedMap(sb, types.AVSKeyPrefix, "avs", collections.StringKey, codec.CollValue[types.AVS](cdc), indices),
		idIndex:       idMulti,
		wasmMsgServer: wasmMsgServer,
	}

	// create the schema
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.schema = schema
	return k
}

func (k *Keeper) HasAVSContract(ctx sdk.Context, contractBytes []byte) bool {
	//calculate sha256 of contract bytes
	ok, err := k.avsMap.Has(ctx, hex.EncodeToString(crypto.Sha256(contractBytes)))
	if err != nil || !ok {
		return false
	}

	return true
}

func (k *Keeper) NextAVSID(ctx sdk.Context) (uint64, error) {
	return k.nextAVSID.Peek(ctx)
}

func (k *Keeper) RegisterAVS(ctx sdk.Context, cp types.MsgRegisterAVS) error {
	if k.HasAVSContract(ctx, cp.GetContractBin()) {
		return types.NewAVSContractAlreadyExistsError(hex.EncodeToString(cp.GetContractBin()))
	}

	id, err := k.nextAVSID.Next(ctx)
	if err != nil {
		return err
	}

	// TODO
	//StoreCode
	//InstantiateContract

	state := types.AVS{
		ContractAddress:    "0x000", // TODO: use actual contract address
		Id:                 id,
		SidecarDockerImage: cp.SidecarDockerImage,
	}
	return k.avsMap.Set(ctx, hex.EncodeToString(crypto.Sha256(cp.GetContractBin())), state)
}

func (k *Keeper) GetIDForAVSContract(ctx sdk.Context, contractBytes []byte) (uint64, bool) {
	cps, err := k.avsMap.Get(ctx, hex.EncodeToString(crypto.Sha256(contractBytes)))
	if err != nil {
		return 0, false
	}

	return cps.Id, true
}

func (k *Keeper) GetAVSByID(ctx sdk.Context, id uint64) (types.AVS, bool) {
	// use the ID index to match the given ID
	ids, err := k.idIndex.MatchExact(ctx, id)
	if err != nil {
		return types.AVS{}, false
	}
	// close the iterator
	defer ids.Close()
	if !ids.Valid() {
		return types.AVS{}, false
	}

	cps, err := ids.PrimaryKey()
	if err != nil {
		return types.AVS{}, false
	}

	avs, err := k.avsMap.Get(ctx, cps)
	if err != nil {
		return types.AVS{}, false
	}

	return avs, true
}
