package client

import (
	"github.com/CosmWasm/wasmd/x/slpp/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func RegisterAVSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "register-avs [id]",
		RunE: func(cmd *cobra.Command, args []string) error {
			msg := &types.MsgRegisterAVS{}
			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
