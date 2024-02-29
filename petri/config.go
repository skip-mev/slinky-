package petri

import (
	"context"
	"fmt"
	"math/big"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/CosmWasm/wasmd/x/slpp"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/skip-mev/petri/core/v2/provider"
	"github.com/skip-mev/petri/core/v2/provider/docker"
	petritypes "github.com/skip-mev/petri/core/v2/types"
	"github.com/skip-mev/petri/cosmos/v2/chain"
	"github.com/skip-mev/petri/cosmos/v2/node"
	"go.uber.org/zap"
)

const (
	denom            = "slinkypp"
	prefix           = "slpp"
	homeDir          = "/petri-test"
)

func GetChainConfig() (petritypes.ChainConfig, error) {
	return petritypes.ChainConfig{
		Denom:         denom,
		Decimals:      6,
		NumValidators: 4,
		NumNodes:      2,
		BinaryName:    "wasmd",
		Image: provider.ImageDefinition{
			Image: "dydxprotocol-base",
			UID:   "1000",
			GID:   "1000",
		},
		GasPrices:      fmt.Sprintf("0%s", denom),
		GasAdjustment:  1.5,
		Bech32Prefix:   prefix,
		EncodingConfig: testutil.MakeTestEncodingConfig(wasm.AppModuleBasic{}, slpp.AppModuleBasic{}),
		HomeDir:        homeDir,
		SidecarHomeDir: "/etc",
		CoinType:       "118",
		ChainId:        "slinkypp-1",
		ModifyGenesis:  GetGenesisModifier(),
		WalletConfig: petritypes.WalletConfig{
			DerivationFn:     hd.Secp256k1.Derive(),
			GenerationFn:     hd.Secp256k1.Generate(),
			Bech32Prefix:     prefix,
			HDPath:           hd.CreateHDPath(0, 0, 0),
			SigningAlgorithm: "secp256k1",
		},
		NodeCreator:       node.CreateNode,
		GenesisDelegation: big.NewInt(10_000_000_000_000),
		GenesisBalance:    big.NewInt(100_000_000_000_000),
	}, nil
}

func GetGenesisModifier() petritypes.GenesisModifier {
	var genKVs = []chain.GenesisKV{
		{
			Key:   "app_state.gov.params.voting_period",
			Value: "10s",
		},
		{
			Key:   "app_state.gov.params.expedited_voting_period",
			Value: "5s",
		},
		{
			Key:   "app_state.gov.params.max_deposit_period",
			Value: "1s",
		},
		{
			Key:   "app_state.gov.params.min_deposit.0.denom",
			Value: denom,
		},
		{
			Key:   "app_state.gov.params.min_deposit.0.amount",
			Value: "1",
		},
		{
			Key:   "app_state.gov.params.threshold",
			Value: "0.1",
		},
		{
			Key:   "app_state.gov.params.quorum",
			Value: "0",
		},
		{
			Key:   "consensus.params.abci.vote_extensions_enable_height",
			Value: "2",
		},
		{
			Key:   "app_state.staking.params.bond_denom",
			Value: denom,
		},
	}

	return chain.ModifyGenesis(genKVs)
}


func GetChain(ctx context.Context, logger *zap.Logger, config petritypes.ChainConfig) (petritypes.ChainI, error) {
	prov, err := docker.NewDockerProvider(
		ctx,
		logger,
		"slinky-plus-plus-docker",
	)
	if err != nil {
		return nil, err
	}

	return chain.CreateChain(
		ctx,
		logger,
		prov,
		config,
	)
}
