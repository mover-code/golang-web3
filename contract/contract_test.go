package contract

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/mover-code/golang-web3/jsonrpc"

	"github.com/mover-code/golang-web3/testutil"

	"github.com/mover-code/golang-web3/abi"

	web3 "github.com/mover-code/golang-web3"
	"github.com/stretchr/testify/assert"
)

var (
	addr0  = "0x0000000000000000000000000000000000000000"
	addr0B = web3.HexToAddress(addr0)
)

func TestContractNoInput(t *testing.T) {
	s := testutil.NewTestServer(t, nil)
	defer s.Close()

	cc := &testutil.Contract{}
	cc.AddOutputCaller("set")

	contract, addr := s.DeployContract(cc)

	abi0, err := abi.NewABI(contract.Abi)
	assert.NoError(t, err)

	p, _ := jsonrpc.NewClient(s.HTTPAddr())
	c := NewContract(addr, abi0, p)

	vals, err := c.Call("set", web3.Latest)
	assert.NoError(t, err)
	assert.Equal(t, vals["0"], big.NewInt(1))

	abi1, err := abi.NewABIFromList([]string{
		"function set () view returns (uint256)",
	})
	assert.NoError(t, err)

	c1 := NewContract(addr, abi1, p)
	vals, err = c1.Call("set", web3.Latest)
	assert.NoError(t, err)
	assert.Equal(t, vals["0"], big.NewInt(1))
}

func TestContractIO(t *testing.T) {
	s := testutil.NewTestServer(t, nil)
	defer s.Close()

	cc := &testutil.Contract{}
	cc.AddDualCaller("setA", "address", "uint256")

	contract, addr := s.DeployContract(cc)

	abi, err := abi.NewABI(contract.Abi)
	assert.NoError(t, err)

	p, _ := jsonrpc.NewClient(s.HTTPAddr())
	c := NewContract(addr, abi, p)
	c.SetFrom(s.Account(0))

	resp, err := c.Call("setA", web3.Latest, addr0B, 1000)
	assert.NoError(t, err)

	assert.Equal(t, resp["0"], addr0B)
	assert.Equal(t, resp["1"], big.NewInt(1000))
}

func TestDeployContract(t *testing.T) {
	s := testutil.NewTestServer(t, nil)
	defer s.Close()

	p, _ := jsonrpc.NewClient(s.HTTPAddr())

	cc := &testutil.Contract{}
	cc.AddConstructor("address", "uint256")

	artifact, err := cc.Compile()
	assert.NoError(t, err)

	abi, err := abi.NewABI(artifact.Abi)
	assert.NoError(t, err)

	bin, err := hex.DecodeString(artifact.Bin)
	assert.NoError(t, err)

	txn := DeployContract(p, s.Account(0), abi, bin, web3.Address{0x1}, 1000)

	if err := txn.Do(); err != nil {
		t.Fatal(err)
	}
	if err := txn.Wait(); err != nil {
		t.Fatal(err)
	}

	i := NewContract(txn.Receipt().ContractAddress, abi, p)
	resp, err := i.Call("val_0", web3.Latest)
	assert.NoError(t, err)
	assert.Equal(t, resp["0"], web3.Address{0x1})

	resp, err = i.Call("val_1", web3.Latest)
	assert.NoError(t, err)
	assert.Equal(t, resp["0"], big.NewInt(1000))
}
