package client

import (
	"github.com/CosmWasm/wasmd/x/slpp/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	slppTxCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "SLPP transactions subcommands",
	}

	slppTxCmd.AddCommand(RegisterAVSCmd())

	return slppTxCmd
}

func RegisterAVSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "register-avs [id]",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			dockerImage := args[0]
			contractBin := []byte(args[1])

			msg := &types.MsgRegisterAVS{
				SidecarDockerImage: dockerImage,
				ContractBin:        contractBin,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
