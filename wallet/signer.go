package wallet

import (
	"math/big"

	web3 "github.com/mover-code/golang-web3"

	"github.com/mover-code/golang-web3/fastrlp"
)

type Signer interface {
	// RecoverSender returns the sender to the transaction
	RecoverSender(tx *web3.Transaction) (web3.Address, error)

	// SignTx signs a transaction
	SignTx(tx *web3.Transaction, key *Key) (*web3.Transaction, error)
}

type EIP1155Signer struct {
	chainID uint64
}

func NewEIP155Signer(chainID uint64) *EIP1155Signer {
	return &EIP1155Signer{chainID: chainID}
}

func (e *EIP1155Signer) RecoverSender(tx *web3.Transaction) (web3.Address, error) {
	v := new(big.Int).SetBytes(tx.V).Uint64()
	v -= e.chainID * 2
	v -= 8
	v -= 27

	sig, err := encodeSignature(tx.R, tx.S, byte(v))
	if err != nil {
		return web3.Address{}, err
	}
	addr, err := Ecrecover(signHash(tx, e.chainID), sig)
	if err != nil {
		return web3.Address{}, err
	}
	return addr, nil
}

func (e *EIP1155Signer) SignTx(tx *web3.Transaction, key *Key) (*web3.Transaction, error) {
	hash := signHash(tx, e.chainID)

	sig, err := key.Sign(hash)
	if err != nil {
		return nil, err
	}

	vv := uint64(sig[64]) + 35 + e.chainID*2

	tx.R = sig[:32]
	tx.S = sig[32:64]
	tx.V = new(big.Int).SetUint64(vv).Bytes()
	return tx, nil
}

func signHash(tx *web3.Transaction, chainID uint64) []byte {
	a := fastrlp.DefaultArenaPool.Get()

	v := a.NewArray()
	v.Set(a.NewUint(tx.Nonce))
	v.Set(a.NewUint(tx.GasPrice))
	v.Set(a.NewUint(tx.Gas))
	if tx.To == nil {
		v.Set(a.NewNull())
	} else {
		v.Set(a.NewCopyBytes((*tx.To)[:]))
	}
	v.Set(a.NewBigInt(tx.Value))
	v.Set(a.NewCopyBytes(tx.Input))

	// EIP155
	if chainID != 0 {
		v.Set(a.NewUint(chainID))
		v.Set(a.NewUint(0))
		v.Set(a.NewUint(0))
	}

	hash := keccak256(v.MarshalTo(nil))
	fastrlp.DefaultArenaPool.Put(a)
	return hash
}

func encodeSignature(R, S []byte, V byte) ([]byte, error) {
	sig := make([]byte, 65)
	copy(sig[32-len(R):32], R)
	copy(sig[64-len(S):64], S)
	sig[64] = V
	return sig, nil
}

// It encodes the signature of a transaction.
//
// Args:
//   R ([]byte): The X coordinate of a point R on the curve
//   S ([]byte): the signature
//   V (byte): The recovery ID.
func EncodeSignature(R, S []byte, V byte) ([]byte, error) {
	return encodeSignature(R, S, V)
}

// It takes a signature and a chainID and returns the R, S, and V values of the signature.
//
// Args:
//   sig ([]byte): the signature
//   chainID (uint64): The chain ID of the network you're signing for.
func DecodeSignature(sig []byte, chainID uint64) (R, S, V []byte) {
	vv := uint64(sig[64]) + 35 + chainID*2
	R = sig[:32]
	S = sig[32:64]
	V = new(big.Int).SetUint64(vv).Bytes()
	return
}
