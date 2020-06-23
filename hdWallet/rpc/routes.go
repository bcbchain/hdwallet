package rpc

import (
	rpcserver "github.com/bcbchain/bclib/rpc/lib/server"
)

var Routes = map[string]*rpcserver.RPCFunc{
	// hdWallet api
	"bcb_walletCreate":    rpcserver.NewRPCFunc(WalletCreate, "password,path"),
	"bcb_transfer":        rpcserver.NewRPCFunc(WalletTransfer, "password,path,walletParams"),
	"bcb_transferOffline": rpcserver.NewRPCFunc(WalletTransferOffline, "password,path,walletParams"),

	// block chain api
	"bcb_blockHeight":    rpcserver.NewRPCFunc(BlockHeight, ""),
	"bcb_block":          rpcserver.NewRPCFunc(Block, "height"),
	"bcb_transaction":    rpcserver.NewRPCFunc(Transaction, "txHash"),
	"bcb_balance":        rpcserver.NewRPCFunc(Balance, "address"),
	"bcb_balanceOfToken": rpcserver.NewRPCFunc(BalanceOfToken, "address,tokenAddress,tokenName"),
	"bcb_allBalance":     rpcserver.NewRPCFunc(AllBalance, "address"),
	"bcb_nonce":          rpcserver.NewRPCFunc(Nonce, "address"),
	"bcb_commitTx":       rpcserver.NewRPCFunc(CommitTx, "tx"),
	"bcb_version":        rpcserver.NewRPCFunc(Version, ""),
}
