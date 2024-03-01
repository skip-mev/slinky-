package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/oracle interfaces (messages) on the
// cdc. These types are used for amino serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// register the MsgAddCurrencyPairs for amino serialization
	legacy.RegisterAminoMsg(cdc, &MsgRegisterAVS{}, "slinky/x/oracle/MsgAddCurrencyPairs")
}

// RegisterInterfaces registers the x/oracle messages + message service w/ the InterfaceRegistry (registry).
func RegisterInterfaces(registry types.InterfaceRegistry) {
	// register the MsgAddCurrencyPairs as an implementation of sdk.Msg
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterAVS{},
	)

	// register the x/oracle message-service
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
