package erc20

import (
	"testing"

	web3 "github.com/mover-code/golang-web3"
	"github.com/mover-code/golang-web3/jsonrpc"
	"github.com/mover-code/golang-web3/testutil"
	"github.com/stretchr/testify/assert"
)

var (
	url   = "https://mainnet.infura.io"
	zeroX = web3.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498")
)

func TestERC20Decimals(t *testing.T) {
	c, _ := jsonrpc.NewClient(testutil.TestInfuraEndpoint(t))
	erc20 := NewERC20(zeroX, c)

	decimals, err := erc20.Decimals()
	assert.NoError(t, err)
	if decimals != 18 {
		t.Fatal("bad")
	}
}

func TestERC20Name(t *testing.T) {
	c, _ := jsonrpc.NewClient(testutil.TestInfuraEndpoint(t))
	erc20 := NewERC20(zeroX, c)

	name, err := erc20.Name()
	assert.NoError(t, err)
	assert.Equal(t, name, "0x Protocol Token")
}

func TestERC20Symbol(t *testing.T) {
	c, _ := jsonrpc.NewClient(testutil.TestInfuraEndpoint(t))
	erc20 := NewERC20(zeroX, c)

	symbol, err := erc20.Symbol()
	assert.NoError(t, err)
	assert.Equal(t, symbol, "ZRX")
}

func TestTotalSupply(t *testing.T) {
	c, _ := jsonrpc.NewClient(testutil.TestInfuraEndpoint(t))
	erc20 := NewERC20(zeroX, c)

	supply, err := erc20.TotalSupply()
	assert.NoError(t, err)
	assert.Equal(t, supply.String(), "1000000000000000000000000000")
}
