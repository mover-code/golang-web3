# golang-web3

## install

```shell
go get -u https://github.com/mover-code/golang-web3.git
```

> 导入包时go.mod 下增加 replace github.com/btcsuite/btcutil => github.com/btcsuite/btcutil v1.0.2

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

- [关于事件监听](./event/event_test.go)
