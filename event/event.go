package event

import (
    "fmt"
    "reflect"
    "time"

    web3 "github.com/mover-code/golang-web3"
    "github.com/mover-code/golang-web3/abi"
    "github.com/mover-code/golang-web3/contract"
    "github.com/mover-code/golang-web3/jsonrpc"
)

type (
    LoadWrapper func(d interface{}, v map[string]interface{}, l *web3.Log)

    MyContract struct {
        Addr         string
        Contract     *contract.Contract
        Cli          *jsonrpc.Client
        TimeDuration int64 // after some time do it
    }
)

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

// If the ABI is invalid, panic.
//
// Args:
//   s (string): The ABI string of the contract
func NewAbi(s string) *abi.ABI {
    a, err := abi.NewABI(s)
    if err != nil {
        panic("error")
    }
    return a
}

// `NewContract` takes an address, an ABI string, and a data struct, and returns a `MyContract` struct
//
// Args:
//   addr (string): The address of the contract
//   abiStr (string): The ABI of the contract.
func NewContract(addr, abiStr, rpc string) *MyContract {
    cli := NewCli(rpc)
    return &MyContract{
        Addr:     addr,
        Contract: contract.NewContract(web3.HexToAddress(addr), NewAbi(abiStr), cli),
        Cli:      cli,
    }
}

// A function that is used to parse the logs of a contract.
func (d *MyContract) ParseLogWithWrapper(l *web3.Log, wrapper LoadWrapper, name ...interface{}) {
    for _, n := range name {
        if e, b := d.Contract.Event(reflect.TypeOf(n).Name()); b {
            data, err := e.ParseLog(l)
            if err == nil {
                wrapper(n, data, l)
            }
        }
    }
}

// Creating a filter for the contract events.
func (d *MyContract) NewFilter(name ...interface{}) *web3.LogFilter {
    topics := []*web3.Hash{}
    for _, n := range name {
        e, b := d.Contract.Event(reflect.TypeOf(n).Name())
        if b {
            topic := e.Encode()
            topics = append(topics, &topic)
        }
    }
    return &web3.LogFilter{
        Address: []web3.Address{web3.HexToAddress(d.Addr)},
        Topics:  topics,
    }
}

// It returns the current block number
func (d *MyContract) NowBlock() web3.BlockNumber {
    blockNumber, _ := d.Cli.Eth().BlockNumber()
    return web3.BlockNumber(blockNumber)
}

// The above code is a Go function that is used to get the logs of a contract.
func (d *MyContract) GetLogs(wrapper LoadWrapper, name ...interface{}) {
    f := d.NewFilter(name...)
    logsInfo := make(chan *web3.Log)
    go func() {
        n := d.NowBlock()
        // init from now-5 block start catch logs
        f.SetToUint64(uint64(n - 5))
        old := *f.To
        f.SetFromUint64(uint64(old))
        for {
            time.Sleep(time.Second * time.Duration(d.TimeDuration))
            now := d.NowBlock()
            if now > old {
                logs, err := d.Cli.Eth().GetLogs(f)
                if err == nil && len(logs) > 0 {
                    for _, l := range logs {
                        logsInfo <- l
                    }
                }
                f.SetFromUint64(uint64(now))
                f.SetToUint64(uint64(now))
                old = now
            }
        }
    }()

    for {
        select {
        case l := <-logsInfo:
            d.ParseLogWithWrapper(l, wrapper, name...)
        }
    }
}

// The above code is a function that is used to get the history logs of a contract.
func (d *MyContract) GetHistoryLogs(start, step int64, wrapper LoadWrapper, name ...interface{}) {
    f := d.NewFilter(name...)
    stop := false
    block := d.NowBlock()
    for {
        time.Sleep(time.Millisecond * time.Duration(d.TimeDuration))
        if start < int64(block) {
            newBlock := start + step
            if newBlock > int64(block) {
                newBlock = int64(block)
                stop = true
            }
            f.SetFromUint64(uint64(start))
            f.SetToUint64(uint64(newBlock))
            logs, err := d.Cli.Eth().GetLogs(f)
            if err == nil && len(logs) > 0 {
                for _, l := range logs {
                    d.ParseLogWithWrapper(l, wrapper, name...)
                }
            }
            start = newBlock
        }
        if stop {
            break
        }
    }
}

// A function that is used to call the contract method.
func (d *MyContract) Call(method string, param ...interface{}) (interface{}, error) {
    time.Sleep(time.Millisecond * time.Duration(d.TimeDuration))
    return d.Contract.Call(method, d.NowBlock(), param...)
}
