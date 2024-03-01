package main

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"encoding/json"

	slpptypes "github.com/CosmWasm/wasmd/x/slpp/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	cmthttp "github.com/cometbft/cometbft/rpc/client/http"
	clientTx "github.com/cosmos/cosmos-sdk/client/tx"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
	privKey = "35c5eda475033fc7dd17fe65efb3b05d14335ee988127d68aa8a1dbb63d3ec5e"
	accNum = 144030
	seq = 4
	contractAddress = "neutron160htmszmgr5wc7sx364e7z22r2ze5j70nx7nreww6qjkqqtqcx0s7vaat0"
)

func main() {
	sdk.GetConfig().SetBech32PrefixForAccount("neutron", "neutronpubkey")
	
	pk, err := hex.DecodeString(privKey)
	if err != nil {
		panic(err)
	}
	privKey := secp256k1.PrivKey{Key: pk}
	
	type SubmitRootBody struct { 
		ChainID string `json:"chain_id"`
		Root    []byte `json:"root"`
	}
	type SubmitRoot struct {
		SubmitRootBody `json:"submit_root"`
	}

	rootBz, err := hex.DecodeString("dce2799628ffd02598141491ae5fa9c5186b75a3c05b66d626402a5013eefdae")
	if err != nil {
		panic(err)
	}

	bz, err := json.Marshal(SubmitRoot{
		SubmitRootBody: SubmitRootBody{
			ChainID: "1",
			Root: rootBz,
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("bz", string(bz))

	// generate the message
	msg := &wasmtypes.MsgExecuteContract{
		Sender: sdk.AccAddress(privKey.PubKey().Address()).String(),
		Contract: contractAddress,
		Msg: bz,
	}

	client, err := cmthttp.New("https://neutron-rpc.lavenderfive.com:443", "/websocket")
	if err != nil {
		panic(err)
	}

	// generate the transaction
	txBytes, err := genTx(&privKey, msg)
	if err != nil {
		panic(err)
	}

	// broadcast the transaction
	res, err := client.BroadcastTxSync(context.Background(), txBytes)
	if err != nil {
		panic(err)
	}

	// print the result
	fmt.Println("response", res)
}


func genTx(pk cryptotypes.PrivKey, msg sdk.Msg) ([]byte, error) {
	ir := codectypes.NewInterfaceRegistry()
	slpptypes.RegisterInterfaces(ir)
	authtypes.RegisterInterfaces(ir)
	cryptocodec.RegisterInterfaces(ir)

	txc := tx.NewTxConfig(codec.NewProtoCodec(ir), tx.DefaultSignModes)

	txb := txc.NewTxBuilder()

	txb.SetMsgs(msg)
	txb.SetGasLimit(4636701)
	txb.SetFeeAmount(sdk.NewCoins(sdk.NewInt64Coin("untrn", 24575)))

	var sigsV2 []signing.SignatureV2
	sigV2 := signing.SignatureV2{
		PubKey: pk.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  signing.SignMode(txc.SignModeHandler().DefaultMode()),
			Signature: nil,
		},
		Sequence: seq,
	}

	sigsV2 = append(sigsV2, sigV2)

	err := txb.SetSignatures(sigsV2...)
	if err != nil {
		panic(err)
	}

	// Second round: all signer infos are set, so each signer can sign.
	sigsV2 = []signing.SignatureV2{}
	signerData := xauthsigning.SignerData{
		ChainID:       "neutron-1",
		AccountNumber: accNum,
		Sequence:      seq,
		PubKey: 	  pk.PubKey(),
	}
	sigV2, err = clientTx.SignWithPrivKey(
		context.Background(),
		signing.SignMode(txc.SignModeHandler().DefaultMode()), signerData,
		txb, pk, txc, uint64(seq))
	if err != nil {
		panic(err)
	}

	sigsV2 = append(sigsV2, sigV2)

	err = txb.SetSignatures(sigsV2...)
	if err != nil {
		panic(err)
	}

	return txc.TxEncoder()(txb.GetTx())
}
