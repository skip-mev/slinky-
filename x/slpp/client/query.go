package client

import (
	"strconv"

	"github.com/CosmWasm/wasmd/x/slpp/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the slpp module",
		RunE:  client.ValidateCmd,
	}

	cmd.AddCommand(
		GetAVSCmd(),
	)

	return cmd
}

func GetAVSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "avs [id]",
		Short: "Query for the price of a specified currency-pair",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// get context
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// parse id
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			// create client
			qc := types.NewQueryClient(clientCtx)

			// query for prices
			res, err := qc.GetAVS(cmd.Context(), &types.GetAVSRequest{
				Id: uint64(id),
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
