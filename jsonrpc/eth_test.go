package jsonrpc

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	web3 "github.com/mover-code/golang-web3"
	"github.com/mover-code/golang-web3/testutil"
	"github.com/stretchr/testify/assert"
)

var (
	addr0 = web3.Address{0x1}
	addr1 = web3.Address{0x2}
)

func TestEthAccounts(t *testing.T) {
	testutil.MultiAddr(t, nil, func(s *testutil.TestServer, addr string) {
		c, _ := NewClient(addr)
		defer c.Close()

		_, err := c.Eth().Accounts()
		assert.NoError(t, err)
	})
}

func TestEthBlockNumber(t *testing.T) {
	i := uint64(0)
	testutil.MultiAddr(t, nil, func(s *testutil.TestServer, addr string) {
		c, _ := NewClient(addr)
		defer c.Close()

		for count := 0; count < 10; count, i = count+1, i+1 {
			num, err := c.Eth().BlockNumber()
			assert.NoError(t, err)
			assert.Equal(t, num, i)
			assert.NoError(t, s.ProcessBlock())
			count++
		}
	})
}

func TestEthGetCode(t *testing.T) {
	s := testutil.NewTestServer(t, nil)
	defer s.Close()

	c, _ := NewClient(s.HTTPAddr())

	cc := &testutil.Contract{}
	cc.AddEvent(testutil.NewEvent("A").
		Add("address", true).
		Add("address", true))

	cc.EmitEvent("setA1", "A", addr0.String(), addr1.String())
	cc.EmitEvent("setA2", "A", addr1.String(), addr0.String())

	_, addr := s.DeployContract(cc)

	code, err := c.Eth().GetCode(addr, web3.Latest)
	assert.NoError(t, err)
	assert.NotEqual(t, code, "0x")

	code2, err := c.Eth().GetCode(addr, web3.BlockNumber(0))
	assert.NoError(t, err)
	assert.Equal(t, code2, "0x")
}

func TestEthGetBalance(t *testing.T) {
	s := testutil.NewTestServer(t, nil)
	defer s.Close()

	c, _ := NewClient(s.HTTPAddr())

	before, err := c.Eth().GetBalance(s.Account(0), web3.Latest)
	assert.NoError(t, err)

	amount := big.NewInt(10)
	txn := &web3.Transaction{
		From:  s.Account(0),
		To:    &testutil.DummyAddr,
		Value: amount,
	}
	receipt, err := s.SendTxn(txn)
	assert.NoError(t, err)

	after, err := c.Eth().GetBalance(s.Account(0), web3.Latest)
	assert.NoError(t, err)

	// the balance in 'after' must be 'before' - 'amount'
	assert.Equal(t, new(big.Int).Add(after, amount).Cmp(before), 0)

	// get balance at block 0
	before2, err := c.Eth().GetBalance(s.Account(0), web3.BlockNumber(0))
	assert.NoError(t, err)
	assert.Equal(t, before, before2)

	{
		// query the balance with different options
		cases := []web3.BlockNumberOrHash{
			web3.Latest,
			receipt.BlockHash,
			web3.BlockNumber(receipt.BlockNumber),
		}
		for _, ca := range cases {
			res, err := c.Eth().GetBalance(s.Account(0), ca)
			assert.NoError(t, err)
			assert.Equal(t, res, after)
		}
	}
}

func TestEthGetBlockByNumber(t *testing.T) {
	s := testutil.NewTestServer(t, nil)
	defer s.Close()

	c, _ := NewClient(s.HTTPAddr())

	block, err := c.Eth().GetBlockByNumber(0, true)
	assert.NoError(t, err)
	assert.Equal(t, block.Number, uint64(0))

	// block 1 has not been processed yet, do not fail but returns nil
	block, err = c.Eth().GetBlockByNumber(1, true)
	assert.NoError(t, err)
	assert.Nil(t, block)

	// process a new block
	assert.NoError(t, s.ProcessBlock())

	// there exists a block 1 now
	block, err = c.Eth().GetBlockByNumber(1, true)
	assert.NoError(t, err)
	assert.Equal(t, block.Number, uint64(1))
}

func TestEthGetBlockByHash(t *testing.T) {
	testutil.MultiAddr(t, nil, func(s *testutil.TestServer, addr string) {
		c, _ := NewClient(addr)
		defer c.Close()

		// get block 0 first by number
		block, err := c.Eth().GetBlockByNumber(0, true)
		assert.NoError(t, err)
		assert.Equal(t, block.Number, uint64(0))

		// get block 0 by hash
		block2, err := c.Eth().GetBlockByHash(block.Hash, true)
		assert.NoError(t, err)
		assert.Equal(t, block, block2)
	})
}

func TestEthGasPrice(t *testing.T) {
	testutil.MultiAddr(t, nil, func(s *testutil.TestServer, addr string) {
		c, _ := NewClient(addr)
		defer c.Close()

		_, err := c.Eth().GasPrice()
		assert.NoError(t, err)
	})
}

func TestEthSendTransaction(t *testing.T) {
	s := testutil.NewTestServer(t, nil)
	defer s.Close()

	c, _ := NewClient(s.HTTPAddr())

	txn := &web3.Transaction{
		From:     s.Account(0),
		GasPrice: testutil.DefaultGasPrice,
		Gas:      testutil.DefaultGasLimit,
		To:       &testutil.DummyAddr,
		Value:    big.NewInt(10),
	}
	hash, err := c.Eth().SendTransaction(txn)
	assert.NoError(t, err)

	var receipt *web3.Receipt
	for {
		receipt, err = c.Eth().GetTransactionReceipt(hash)
		if err != nil {
			t.Fatal(err)
		}
		if receipt != nil {
			break
		}
	}
}

func TestEthEstimateGas(t *testing.T) {
	s := testutil.NewTestServer(t, nil)
	defer s.Close()

	c, _ := NewClient(s.HTTPAddr())

	cc := &testutil.Contract{}
	cc.AddEvent(testutil.NewEvent("A").Add("address", true))
	cc.EmitEvent("setA", "A", addr0.String())

	// estimate gas to deploy the contract
	solcContract, err := cc.Compile()
	assert.NoError(t, err)

	input, err := hex.DecodeString(solcContract.Bin)
	assert.NoError(t, err)

	msg := &web3.CallMsg{
		From: s.Account(0),
		To:   nil,
		Data: input,
	}
	gas, err := c.Eth().EstimateGas(msg)
	assert.NoError(t, err)
	assert.NotEqual(t, gas, 0)

	_, addr := s.DeployContract(cc)

	msg = &web3.CallMsg{
		From: s.Account(0),
		To:   &addr,
		Data: testutil.MethodSig("setA"),
	}

	gas, err = c.Eth().EstimateGas(msg)
	assert.NoError(t, err)
	assert.NotEqual(t, gas, 0)
}

func TestEthGetLogs(t *testing.T) {
	s := testutil.NewTestServer(t, nil)
	defer s.Close()

	c, _ := NewClient(s.HTTPAddr())

	cc := &testutil.Contract{}
	cc.AddEvent(testutil.NewEvent("A").
		Add("address", true).
		Add("address", true))

	cc.EmitEvent("setA1", "A", addr0.String(), addr1.String())
	cc.EmitEvent("setA2", "A", addr1.String(), addr0.String())

	_, addr := s.DeployContract(cc)

	r := s.TxnTo(addr, "setA2")

	filter := &web3.LogFilter{
		BlockHash: &r.BlockHash,
	}
	logs, err := c.Eth().GetLogs(filter)
	assert.NoError(t, err)
	assert.Len(t, logs, 1)

	log := logs[0]
	assert.Len(t, log.Topics, 3)
	assert.Equal(t, log.Address, addr)

	// first topic is the signature of the event
	assert.Equal(t, log.Topics[0].String(), cc.GetEvent("A").Sig())

	// topics have 32 bytes and the addr are 20 bytes, then, assert.Equal wont work.
	// this is a workaround until we build some helper function to test this better
	assert.True(t, bytes.HasSuffix(log.Topics[1][:], addr1[:]))
	assert.True(t, bytes.HasSuffix(log.Topics[2][:], addr0[:]))
}

func TestEthChainID(t *testing.T) {
	testutil.MultiAddr(t, nil, func(s *testutil.TestServer, addr string) {
		c, _ := NewClient(addr)
		defer c.Close()

		num, err := c.Eth().ChainID()
		assert.NoError(t, err)
		assert.Equal(t, num.Uint64(), uint64(1337)) // chainid of geth-dev
	})
}

func TestEthGetNonce(t *testing.T) {
	s := testutil.NewTestServer(t, nil)
	defer s.Close()

	c, _ := NewClient(s.HTTPAddr())

	num, err := c.Eth().GetNonce(s.Account(0), web3.Latest)
	assert.NoError(t, err)
	assert.Equal(t, num, uint64(0))

	receipt, err := s.ProcessBlockWithReceipt()
	assert.NoError(t, err)

	// query the balance with different options
	cases := []web3.BlockNumberOrHash{
		web3.Latest,
		receipt.BlockHash,
		web3.BlockNumber(receipt.BlockNumber),
	}
	for _, ca := range cases {
		num, err = c.Eth().GetNonce(s.Account(0), ca)
		assert.NoError(t, err)
		assert.Equal(t, num, uint64(1))
	}
}

func TestEthTransactionsInBlock(t *testing.T) {
	s := testutil.NewTestServer(t, nil)
	defer s.Close()

	c, _ := NewClient(s.HTTPAddr())

	_, err := c.Eth().GetBlockByNumber(0, false)
	assert.NoError(t, err)

	// Process a block with a transaction
	assert.NoError(t, s.ProcessBlock())

	// get a non-full block
	block0, err := c.Eth().GetBlockByNumber(1, false)
	assert.NoError(t, err)

	assert.Len(t, block0.TransactionsHashes, 1)
	assert.Len(t, block0.Transactions, 0)

	// get a full block
	block1, err := c.Eth().GetBlockByNumber(1, true)
	assert.NoError(t, err)

	assert.Len(t, block1.TransactionsHashes, 0)
	assert.Len(t, block1.Transactions, 1)

	assert.Equal(t, block0.TransactionsHashes[0], block1.Transactions[0].Hash)
}

func TestEthGetStorageAt(t *testing.T) {
	s := testutil.NewTestServer(t, nil)
	defer s.Close()

	c, _ := NewClient(s.HTTPAddr())

	cc := &testutil.Contract{}

	// add global variables
	cc.AddCallback(func() string {
		return "uint256 val;"
	})

	// add setter method
	cc.AddCallback(func() string {
		return `function setValue() public payable {
			val = 10;
		}`
	})

	_, addr := s.DeployContract(cc)
	receipt := s.TxnTo(addr, "setValue")

	cases := []web3.BlockNumberOrHash{
		web3.Latest,
		receipt.BlockHash,
		web3.BlockNumber(receipt.BlockNumber),
	}
	for _, ca := range cases {
		res, err := c.Eth().GetStorageAt(addr, web3.Hash{}, ca)
		assert.NoError(t, err)
		assert.True(t, strings.HasSuffix(res.String(), "a"))
	}
}
