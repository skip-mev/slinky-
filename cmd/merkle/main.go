package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"github.com/cbergoon/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "net/http/pprof" //nolint: gosec
)

type TestContent struct {
	tx *types.Transaction
}

// CalculateHash hashes the values of a TestContent
func (t TestContent) CalculateHash() ([]byte, error) {
	return t.tx.Hash().Bytes(), nil
}

// Equals tests for equality of two Contents
func (t TestContent) Equals(other merkletree.Content) (bool, error) {
	h1, err := t.CalculateHash()
	if err != nil {
		return false, err
	}
	h2, err := t.CalculateHash()
	if err != nil {
		return false, err
	}
	return bytes.Compare(h1, h2) == 0, nil
}

// start the oracle-grpc server + oracle process, cancel on interrupt or terminate.
func main() {

	// create server
	conn, err := ethclient.Dial("https://ethereum-rpc.publicnode.com")
	if err != nil {
		panic(err.Error())
	}
	receipt, err := conn.TransactionReceipt(context.Background(), common.HexToHash("0x675d94edf2c4fe9b4de7b8633d46ab8c6cba932bfcdbb733abb1106bb85aa9e0"))
	if err != nil {
		panic(err.Error())
	}
	block, err := conn.BlockByNumber(context.Background(), receipt.BlockNumber)
	if err != nil {
		panic(err.Error())
	}
	txs := block.Transactions()
	var contents []merkletree.Content
	txContent := TestContent{txs[receipt.TransactionIndex]}
	for i := range txs {
		contents = append(contents, TestContent{txs[i]})
	}
	m, err := merkletree.NewTree(contents)
	if err != nil {
		panic(err.Error())
	}
	a, b, err := m.GetMerklePath(txContent)
	if err != nil {
		panic(err.Error())
	}
	println(a)
	println(b)
	println(hex.EncodeToString(m.MerkleRoot()))
}

//type TestContent struct {
//	i uint8
//}
//
//// CalculateHash hashes the values of a TestContent
//func (t TestContent) CalculateHash() ([]byte, error) {
//	return crypto.Sha256([]byte{t.i}), nil
//}
//
//// Equals tests for equality of two Contents
//func (t TestContent) Equals(other merkletree.Content) (bool, error) {
//	h1, err := t.CalculateHash()
//	if err != nil {
//		return false, err
//	}
//	h2, err := t.CalculateHash()
//	if err != nil {
//		return false, err
//	}
//	return bytes.Compare(h1, h2) == 0, nil
//}
//
//// start the oracle-grpc server + oracle process, cancel on interrupt or terminate.
//func main() {
//
//	var contents []merkletree.Content
//	leaves := []uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
//	for i := range leaves {
//		contents = append(contents, TestContent{leaves[i]})
//	}
//	m, err := merkletree.NewTree(contents)
//	if err != nil {
//		panic(err.Error())
//	}
//	println(hex.EncodeToString(m.MerkleRoot()))
//}
