# golang-web3

## install

```shell
go get -u github.com/mover-code/golang-web3.git@latest
```

> go get -u github.com/btcsuite/btcd@v0.22.0-beta 遇到ambiguous import

## desc

使用go操作区块链账户，发起交易，查询账户资产、合约交互、部署，事件监听等功能

## use

- 创建一个web3实例

```go
// It creates a new client for the given url.
//
// Args:
//   url (string): The address of the RPC service.
func NewCli(url string) *jsonrpc.Client {
    cli, err := jsonrpc.NewClient(url)
    if err != nil {
        panic(fmt.Sprintf("error:%s", url))
    }
    return cli
}
```

- 数据签名示例

```go
    privateKey, _ := crypto.HexToECDSA(key)  // 账户私钥
    u256, _ := abi.NewType("uint256")
    account, _ := abi.NewType("address")
    type Ord struct {
        Id      *big.Int
        Account web3.Address
        Amount  *big.Int
        Fee     *big.Int
        Solt    *big.Int
        End     *big.Int
        Type    *big.Int
        State   *big.Int
    }

     argumentInfo := []*abi.TupleElem{
        &abi.TupleElem{
            Name: "id",
            Elem: u256,
        },
        &abi.TupleElem{
            Name: "account",
            Elem: account,
        },
        &abi.TupleElem{
            Name: "amount",
            Elem: u256,
        },
        &abi.TupleElem{
            Name: "solt",
            Elem: u256,
        },
        &abi.TupleElem{
            Name: "end",
            Elem: u256,
        },
        &abi.TupleElem{
            Name: "type",
            Elem: u256,
        },
    }

    param := abi.NewTupleType(argumentInfo)
    
    // 数据hash
    hash, _ := param.Encode(&Ord{ 
        Id:      id,
        Account: addr,
        Amount:  withDrawAmount,
        Fee:     big.NewInt(0),
        Solt:    solt,
        End:     endTime,
        Type:    wType,
        State:   big.NewInt(0),
    })

    hashs := crypto.Keccak256Hash(signString(""), crypto.Keccak256(hash)) // 按照合约中的方式组装数据
    signature, _ := crypto.Sign(hashs.Bytes(), privateKey) // 使用私钥签名获取签名后的数据
    r,s,v := SignRSV(signature) // 签名后的rsv数据
```

- 验证签名信息示例

```go
func VerifySig(msg, from, sigHex string) bool {

    sig := hexutil.MustDecode(sigHex)
    if sig[64] != 27 && sig[64] != 28 {
        return false
    }
    sig[64] -= 27

    addr, _ := wallet.Ecrecover(web3.SignHash([]byte(msg)), sig)
    return addr == web3.HexToAddress(from)
}
```

- 代币转账示例 主币

```go

import(
    web3 "github.com/mover-code/golang-web3"
    "github.com/mover-code/golang-web3/jsonrpc"
    "github.com/mover-code/golang-web3/wallet"

)
     // rpc地址
    cli, err := jsonrpc.NewClient(url)
    if err != nil {
        panic(fmt.Sprintf("error:%s", url))
    }

    block, _ := cli.Eth().BlockNumber() // 获取当前区块
    balance, _ := cli.Eth().GetBalance(web3.HexToAddress(addr), web3.BlockNumber(block)) // 查询账户余额

    nonce, err := cli.Eth().GetNonce(sender, web3.BlockNumber(block)) // 发起交易账户的序号
    gas, _ := cli.Eth().GasPrice() // 当前gas
    gasLimit, err := cli.Eth().EstimateGas(&web3.CallMsg{
        From:     sender,    // 发起账户
        To:       &receiver, // 接收账户
        GasPrice: gas,
        Value:    big.NewInt(amount), // 转账金额
    })
    // 组装交易信息
    t := &web3.Transaction{
        From:        sender,
        To:          &receiver,
        Value:       big.NewInt(100000000000000000),
        Nonce:       nonce,
        BlockNumber: uint64(block),
        GasPrice:    gas,
        Gas:         gasLimit,
    }

    chainId, _ := cli.Eth().ChainID()
    signer := wallet.NewEIP155Signer(chainId.Uint64()) 
    byteKey, _ := hex.DecodeString(private_hash) // 发起交易账户 私钥
    key, _ := wallet.NewWalletFromPrivKey(byteKey)
    data, err := signer.SignTx(t, key) // 签名数据
    cli.Eth().SenSignTransaction(data) // 发起交易

```

- [关于事件监听](./event/event_test.go)
