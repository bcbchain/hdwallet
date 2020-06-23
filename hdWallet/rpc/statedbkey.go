package rpc

import (
	"github.com/bcbchain/bcbchain/abciapp_v1.0/smc"
	"strings"
)

func keyOfGenesisToken() string {
	return "/genesis/token"
}

func keyOfContract(contractAddr smc.Address) string {
	return "/contract/" + contractAddr
}

func keyOfToken(tokenAddress smc.Address) string {
	return "/token/" + tokenAddress
}

func keyOfBVMContract(contractAddr smc.Address) string {
	return "/bvm/contract/" + contractAddr
}

func keyOfTokenName(tokenName string) string {
	return "/token/name/" + strings.ToLower(tokenName)
}

func keyOfAccountToken(exAddress smc.Address, contractAddr smc.Address) string {
	return "/account/ex/" + exAddress + "/token/" + contractAddr
}

func keyOfAccount(exAddress smc.Address) string {
	return "/account/ex/" + exAddress
}

func keyOfAccountNonce(exAddress smc.Address) string {
	return "/account/ex/" + exAddress + "/account"
}
