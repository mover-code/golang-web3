/*
 * @Author: small_ant xms.chnb@gmail.com
 * @Time: 2022-01-15 13:50:16
 * @LastAuthor: small_ant xms.chnb@gmail.com
 * @lastTime: 2022-12-08 14:14:42
 * @FileName: units
 * @Desc:
 *
 * Copyright (c) 2022 by small_ant xms.chnb@gmail.com, All Rights Reserved.
 */
package web3

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/shopspring/decimal"
)

func convert(val uint64, decimals int64) *big.Int {
	v := big.NewInt(int64(val))
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil)
	return v.Mul(v, exp)
}

func devert(v *big.Int, decimals int64) *big.Int {
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil)
	return v.Div(v, exp)
}

// Ether converts a value to the ether unit with 18 decimals
func Ether(i uint64) *big.Int {
	return convert(i, 18)
}

// Gwei converts a value to the gwei unit with 9 decimals

func Gwei(i uint64) *big.Int {
	return convert(i, 9)
}

// `FWei` converts a `big.Int` from wei to ether
//
// Args:
//   i: The amount of Wei you want to convert.
func FWei(i *big.Int) *big.Int {
	return devert(i, 18)
}

// It takes a signature in the form of a byte array or hex string, and returns the R, S, and V values
// as byte arrays and a uint8
//
// Args:
//   isig: signature
func SignRSV(isig interface{}) ([32]byte, [32]byte, uint8) {
	var sig []byte
	switch v := isig.(type) {
	case []byte:
		sig = v
	case string:
		sig, _ = hexutil.Decode(v)
	}

	signStr := common.Bytes2Hex(sig)
	rS := signStr[0:64]
	sS := signStr[64:128]
	R := [32]byte{}
	S := [32]byte{}
	copy(R[:], common.FromHex(rS))
	copy(S[:], common.FromHex(sS))
	vStr := signStr[128:130]
	vI, _ := strconv.Atoi(vStr)
	V := uint8(vI + 27)

	return R, S, V
}

// It takes an interface{} and converts it to a decimal.Decimal
//
// Args:
//   ivalue: The value to convert to a decimal.
//   decimals (int): The number of decimal places to round to.
func ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

// It takes an amount and a decimal count, and returns the amount in wei
//
// Args:
//   iamount: The amount of tokens you want to send.
//   decimals (int): The number of decimals the token uses.
func ToWei(iamount interface{}, decimals int) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iamount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}

// It takes a byte array, and returns a byte array
//
// Args:
//   buf ([]byte): The data to be hashed.
//
// Returns:
//   The hash of the input buffer.
func Keccak256(buf []byte) []byte {
	return crypto.Keccak256(buf)
}

// It takes a byte array, converts it to a string, prepends a string to it, and then hashes the result
//
// Args:
//   data ([]byte): The data to sign.
//
// Returns:
//   The hash of the message.
func SignHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return Keccak256([]byte(msg))
}

// It takes a string and returns a byte array
//
// Args:
//   data (string): The data to sign.
//
// Returns:
//   The return value is a byte array.
func SignString(data string) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", 32)
	return []byte(msg)
}

// It takes a variable number of byte slices and returns a single byte slice
// 
// Returns:
//   The hash of the data.
func Keccak256HashData(data ...[]byte) []byte {
	h := crypto.Keccak256Hash(data...)
	return h.Bytes()
}
