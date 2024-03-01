package slpp

import (
	"context"
	"encoding/json"

	"cosmossdk.io/core/appmodule"

	cmtabci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/CosmWasm/wasmd/x/slpp/types"
	"github.com/CosmWasm/wasmd/x/slpp/keeper"
	slppclient "github.com/CosmWasm/wasmd/x/slpp/client"
	"github.com/spf13/cobra"
)

// ConsensusVersion is the x/alerts module's current version, as modules integrate and
// updates are made, this value determines what version of the module is being run by the chain.
const ConsensusVersion = 1

var (
	_ module.AppModuleBasic = AppModule{}
	_ module.HasServices    = AppModule{}

	_ appmodule.AppModule = AppModule{}
)

// AppModuleBasic defines the base interface that the x/alerts module exposes to the
// application.
type AppModuleBasic struct {}

// Name returns the name of this module.
func (AppModuleBasic) Name() string { return types.ModuleName }

// RegisterLegacyAminoCodec registers the necessary types from the x/alerts module
// for amino serialization.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// register alerts legacy amino codec
	types.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the necessary implementations / interfaces in the
// x/alerts module w/ the interface-registry.
func (AppModuleBasic) RegisterInterfaces(ir codectypes.InterfaceRegistry) {
	// register the msgs / interfaces for the alerts module
	types.RegisterInterfaces(ir)
}

// RegisterGRPCGatewayRoutes registers the necessary REST routes for the GRPC-gateway to
// the x/alerts module QueryService on mux. This method panics on failure.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(cliCtx client.Context, mux *runtime.ServeMux) {
	// Register the gate-way routes w/ the provided mux.
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(cliCtx)); err != nil {
		panic(err)
	}
}

// GetTxCmd is a no-op, as no txs are registered for submission.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return slppclient.GetTxCmd()
}

// GetQueryCmd returns the x/alerts module base query cli-command.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return slppclient.GetQueryCmd()
}

// DefaultGenesis returns default genesis state as raw bytes for the alerts
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return nil
}

// ValidateGenesis performs genesis state validation for the alerts module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	return nil
}

// AppModule represents an application module for the x/alerts module.
type AppModule struct {
	AppModuleBasic

	k *keeper.Keeper
}

// NewAppModule returns an application module for the x/alerts module.
func NewAppModule(k *keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		k: k,
	}
}

// IsAppModule implements the appmodule.AppModule interface.
func (AppModule) IsAppModule() {}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModule) IsOnePerModuleType() {}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }

// RegisterServices registers the module's services with the app's module configurator.
func (am AppModule) RegisterServices(cfc module.Configurator) {
	// Register the query service.
	types.RegisterQueryServer(cfc.QueryServer(), keeper.NewQueryServer(am.k))

	// Register the message service.
	types.RegisterMsgServer(cfc.MsgServer(), keeper.NewMsgServer(am.k))
}

// RegisterInvariants registers the invariants of the alerts module. If an invariant
// deviates from its predicted value, the InvariantRegistry triggers appropriate
// logic (most often the chain will be halted). No invariants exist for the alerts module.
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis performs the genesis initialization for the x/alerts module. It determines the
// genesis state to initialize from via a json-encoded genesis-state. This method returns no validator set updates.
// This method panics on any errors.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, bz json.RawMessage) []cmtabci.ValidatorUpdate {
	// return no validator-set updates
	return []cmtabci.ValidatorUpdate{}
}

// ExportGenesis returns the alerts module's exported genesis state as raw
// JSON bytes. This method panics on any error.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return json.RawMessage{}
}
