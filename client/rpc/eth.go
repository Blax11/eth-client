package rpc

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	client "github.com/ethereum/go-ethereum/rpc"
)

type Eth interface {
	PublicTransactionPool
}

type eth struct {
	PublicTransactionPool
}

func NewEth(client *client.Client) Eth {
	return &eth{
		PublicTransactionPool: NewPublicTransactionPool(client),
	}
}

// SendTxArgs represents the arguments to sumbit a new transaction into the transaction pool.
type SendTxArgs struct {
	From     common.Address `json:"from"`
	To       common.Address `json:"to"`
	Gas      hexutil.Big    `json:"gas"`
	GasPrice hexutil.Big    `json:"gasPrice"`
	Value    hexutil.Big    `json:"value"`
	Data     hexutil.Bytes  `json:"data"`
	Nonce    hexutil.Uint64 `json:"nonce"`
}

// SignTransactionResult represents a RLP encoded signed transaction.
type SignTransactionResult struct {
	Raw hexutil.Bytes      `json:"raw"`
	Tx  *types.Transaction `json:"tx"`
}

// RPCTransaction represents a transaction that will serialize to the RPC representation of a transaction
type RPCTransaction struct {
	BlockHash        common.Hash     `json:"blockHash"`
	BlockNumber      *hexutil.Big    `json:"blockNumber"`
	From             common.Address  `json:"from"`
	Gas              *hexutil.Big    `json:"gas"`
	GasPrice         *hexutil.Big    `json:"gasPrice"`
	Hash             common.Hash     `json:"hash"`
	Input            hexutil.Bytes   `json:"input"`
	Nonce            hexutil.Uint64  `json:"nonce"`
	To               *common.Address `json:"to"`
	TransactionIndex hexutil.Uint    `json:"transactionIndex"`
	Value            *hexutil.Big    `json:"value"`
	V                *hexutil.Big    `json:"v"`
	R                *hexutil.Big    `json:"r"`
	S                *hexutil.Big    `json:"s"`
}

type PublicTransactionPool interface {
	// GetBlockTransactionCountByNumber returns the number of transactions in the block with the given block number.
	GetBlockTransactionCountByNumber(ctx context.Context, blockNr string) (*hexutil.Uint, error)
	// GetBlockTransactionCountByHash returns the number of transactions in the block with the given hash.
	GetBlockTransactionCountByHash(ctx context.Context, blockHash common.Hash) (*hexutil.Uint, error)
	// GetTransactionByBlockNumberAndIndex returns the transaction for the given block number and index.
	GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNr string, index hexutil.Uint) (*RPCTransaction, error)
	// GetTransactionByBlockHashAndIndex returns the transaction for the given block hash and index.
	GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (*RPCTransaction, error)
	// GetRawTransactionByBlockNumberAndIndex returns the bytes of the transaction for the given block number and index.
	GetRawTransactionByBlockNumberAndIndex(ctx context.Context, blockNr string, index hexutil.Uint) (hexutil.Bytes, error)
	// GetRawTransactionByBlockHashAndIndex returns the bytes of the transaction for the given block hash and index.
	GetRawTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (hexutil.Bytes, error)
	// GetTransactionCount returns the number of transactions the given address has sent for the given block number
	GetTransactionCount(ctx context.Context, address common.Address, blockNr string) (*hexutil.Uint64, error)
	// GetTransactionByHash returns the transaction for the given hash
	GetTransactionByHash(ctx context.Context, hash common.Hash) (*RPCTransaction, error)
	// GetRawTransactionByHash returns the bytes of the transaction for the given hash.
	GetRawTransactionByHash(ctx context.Context, hash common.Hash) (hexutil.Bytes, error)
	// GetTransactionReceipt returns the transaction receipt for the given transaction hash.
	GetTransactionReceipt(ctx context.Context, hash common.Hash) (map[string]interface{}, error)
	// SendTransaction creates a transaction for the given argument, sign it and submit it to the
	// transaction pool.
	SendTransaction(ctx context.Context, args SendTxArgs) (common.Hash, error)
	// SendRawTransaction will add the signed transaction to the transaction pool.
	// The sender is responsible for signing the transaction and using the correct nonce.
	SendRawTransaction(ctx context.Context, encodedTx hexutil.Bytes) (common.Hash, error)
	// Sign calculates an ECDSA signature for:
	// keccack256("\x19Ethereum Signed Message:\n" + len(message) + message).
	//
	// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
	// where the V value will be 27 or 28 for legacy reasons.
	//
	// The account associated with addr must be unlocked.
	//
	// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign
	Sign(ctx context.Context, addr common.Address, data hexutil.Bytes) (hexutil.Bytes, error)
	// SignTransaction will sign the given transaction with the from account.
	// The node needs to have the private key of the account corresponding with
	// the given from address and it needs to be unlocked.
	SignTransaction(ctx context.Context, args SendTxArgs) (*SignTransactionResult, error)
	// PendingTransactions returns the transactions that are in the transaction pool and have a from address that is one of
	// the accounts this node manages.
	PendingTransactions(ctx context.Context) ([]*RPCTransaction, error)
	// Resend accepts an existing transaction and a new gas price and limit. It will remove
	// the given transaction from the pool and reinsert it with the new gas price and limit.
	Resend(ctx context.Context, sendArgs SendTxArgs, gasPrice, gasLimit hexutil.Big) (common.Hash, error)
}

type publicTransactionPool struct {
	client *client.Client
}

func NewPublicTransactionPool(client *client.Client) PublicTransactionPool {
	return &publicTransactionPool{
		client: client,
	}
}

// GetBlockTransactionCountByNumber returns the number of transactions in the block with the given block number.
func (pub *publicTransactionPool) GetBlockTransactionCountByNumber(ctx context.Context, blockNr string) (*hexutil.Uint, error) {
	var r *hexutil.Uint
	err := pub.client.CallContext(ctx, &r, "eth_getBlockTransactionCountByNumber", blockNr)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetBlockTransactionCountByHash returns the number of transactions in the block with the given hash.
func (pub *publicTransactionPool) GetBlockTransactionCountByHash(ctx context.Context, blockHash common.Hash) (*hexutil.Uint, error) {
	var r *hexutil.Uint
	err := pub.client.CallContext(ctx, &r, "eth_getBlockTransactionCountByHash", blockHash)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetTransactionByBlockNumberAndIndex returns the transaction for the given block number and index.
func (pub *publicTransactionPool) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNr string, index hexutil.Uint) (*RPCTransaction, error) {
	var r *RPCTransaction
	err := pub.client.CallContext(ctx, &r, "eth_getTransactionByBlockNumberAndIndex", blockNr, index)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetTransactionByBlockHashAndIndex returns the transaction for the given block hash and index.
func (pub *publicTransactionPool) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (*RPCTransaction, error) {
	var r *RPCTransaction
	err := pub.client.CallContext(ctx, &r, "eth_getTransactionByBlockNumberAndIndex", blockHash, index)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetRawTransactionByBlockNumberAndIndex returns the bytes of the transaction for the given block number and index.
func (pub *publicTransactionPool) GetRawTransactionByBlockNumberAndIndex(ctx context.Context, blockNr string, index hexutil.Uint) (hexutil.Bytes, error) {
	var r hexutil.Bytes
	err := pub.client.CallContext(ctx, &r, "eth_getRawTransactionByBlockNumberAndIndex", blockNr, index)
	if err != nil {
		return r, err
	}
	return r, nil
}

// GetRawTransactionByBlockHashAndIndex returns the bytes of the transaction for the given block hash and index.
func (pub *publicTransactionPool) GetRawTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (hexutil.Bytes, error) {
	var r hexutil.Bytes
	err := pub.client.CallContext(ctx, &r, "eth_getRawTransactionByBlockHashAndIndex", blockHash, index)
	if err != nil {
		return r, err
	}
	return r, nil
}

// GetTransactionCount returns the number of transactions the given address has sent for the given block number
func (pub *publicTransactionPool) GetTransactionCount(ctx context.Context, address common.Address, blockNr string) (*hexutil.Uint64, error) {
	var r *hexutil.Uint64
	err := pub.client.CallContext(ctx, &r, "eth_getTransactionCount", address, blockNr)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetTransactionByHash returns the transaction for the given hash
func (pub *publicTransactionPool) GetTransactionByHash(ctx context.Context, hash common.Hash) (*RPCTransaction, error) {
	var r *RPCTransaction
	err := pub.client.CallContext(ctx, &r, "eth_getTransactionByHash", hash)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetRawTransactionByHash returns the bytes of the transaction for the given hash.
func (pub *publicTransactionPool) GetRawTransactionByHash(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	var r hexutil.Bytes
	err := pub.client.CallContext(ctx, &r, "eth_getRawTransactionByHash", hash)
	if err != nil {
		return r, err
	}
	return r, nil
}

// GetTransactionReceipt returns the transaction receipt for the given transaction hash.
func (pub *publicTransactionPool) GetTransactionReceipt(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	var r map[string]interface{}
	err := pub.client.CallContext(ctx, &r, "eth_getTransactionReceipt", hash)
	if err != nil {
		return r, err
	}
	return r, nil
}

// SendTransaction creates a transaction for the given argument, sign it and submit it to the
// transaction pool.
func (pub *publicTransactionPool) SendTransaction(ctx context.Context, args SendTxArgs) (common.Hash, error) {
	var r common.Hash
	err := pub.client.CallContext(ctx, &r, "eth_sendTransaction", args)
	if err != nil {
		return r, err
	}
	return r, nil
}

// SendRawTransaction will add the signed transaction to the transaction pool.
// The sender is responsible for signing the transaction and using the correct nonce.
func (pub *publicTransactionPool) SendRawTransaction(ctx context.Context, encodedTx hexutil.Bytes) (common.Hash, error) {
	var r common.Hash
	err := pub.client.CallContext(ctx, &r, "eth_sendRawTransaction", encodedTx)
	if err != nil {
		return r, err
	}
	return r, nil
}

// Sign calculates an ECDSA signature for:
// keccack256("\x19Ethereum Signed Message:\n" + len(message) + message).
//
// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
// where the V value will be 27 or 28 for legacy reasons.
//
// The account associated with addr must be unlocked.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign
func (pub *publicTransactionPool) Sign(ctx context.Context, addr common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	var r hexutil.Bytes
	err := pub.client.CallContext(ctx, &r, "eth_sign", addr, data)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Resend accepts an existing transaction and a new gas price and limit. It will remove
// the given transaction from the pool and reinsert it with the new gas price and limit.
func (pub *publicTransactionPool) Resend(ctx context.Context, sendArgs SendTxArgs, gasPrice, gasLimit hexutil.Big) (common.Hash, error) {
	var r common.Hash
	err := pub.client.CallContext(ctx, &r, "eth_resend", sendArgs, gasPrice, gasLimit)
	if err != nil {
		return r, err
	}
	return r, nil
}

// SignTransaction will sign the given transaction with the from account.
// The node needs to have the private key of the account corresponding with
// the given from address and it needs to be unlocked.
func (pub *publicTransactionPool) SignTransaction(ctx context.Context, args SendTxArgs) (*SignTransactionResult, error) {
	var r *SignTransactionResult
	err := pub.client.CallContext(ctx, &r, "eth_signTransaction", args)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// TODO: how to implement submitTransaction?

// GetRawTransaction returns the bytes of the transaction for the given hash.
func (pub *publicTransactionPool) GetRawTransaction(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	var r hexutil.Bytes
	err := pub.client.CallContext(ctx, &r, "eth_getRawTransactionByHash", hash)
	if err != nil {
		return r, err
	}
	return r, nil
}

// PendingTransactions returns the transactions that are in the transaction pool and have a from address that is one of
// the accounts this node manages.
func (pub *publicTransactionPool) PendingTransactions(ctx context.Context) ([]*RPCTransaction, error) {
	var r []*RPCTransaction
	err := pub.client.CallContext(ctx, &r, "eth_pendingTransactions")
	if err != nil {
		return r, err
	}
	return r, nil
}
