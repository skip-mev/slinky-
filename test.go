package main

import (
	"fmt"
	"context"
	"encoding/hex"
	"encoding/base64"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	slpptypes "github.com/CosmWasm/wasmd/x/slpp/types"
	clientTx "github.com/cosmos/cosmos-sdk/client/tx"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cmthttp "github.com/cometbft/cometbft/rpc/client/http"
)

const (
	privKey = "6c5b738da068835ebc9d625424307108ccdb69d96ce23912c9cfb0153a0a333d"
	binFile = "./bin"
	accNum = 0
	seq = 1
)

func init() {
	sdk.GetConfig().SetBech32PrefixForAccount("wasm", "wasmpubkey")
}

func main() {
	pk, err := hex.DecodeString(privKey)
	if err != nil {
		panic(err)
	}
	privKey := secp256k1.PrivKey{Key: pk}

	bin, err := os.ReadFile(binFile)
	if err != nil {
		panic(err)
	}

	// base64 decode bytecode
	fmt.Println(string(bin))
	binCode, err := base64.RawStdEncoding.DecodeString(string(bin))
	if err != nil {
		panic(err)
	}

	sender := sdk.AccAddress(privKey.PubKey().Address())
	msg := slpptypes.MsgRegisterAVS{
		SidecarDockerImage: "abc",
		Sender: sender.String(),
		ContractBin: binCode,
		InstantiateMsg: []byte("{}"),
	}

	// build the tx
	client, err := cmthttp.New("http://localhost:26657", "/websocket")
	if err != nil {
		panic(err)
	}

	txBytes, err := genTx(&privKey, &msg)
	if err != nil {
		panic(err)
	}

	// broadcast the tx
	res, err := client.BroadcastTxSync(context.Background(), txBytes)
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Code, res.Log, res.Data, res.Hash)
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
	txb.SetFeeAmount(sdk.NewCoins(sdk.NewInt64Coin("stake", 1000000)))

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
		ChainID:       "skip-1",
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
