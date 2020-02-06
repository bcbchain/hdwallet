package client

import (
	"common/rpc/lib/client"
	"encoding/json"
	"fmt"
	rpc3 "hdWallet/rpc"
	"strconv"
)

func WalletCreate(password, path, url string) (err error) {
	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)
	result := new(rpc3.WalletCreateResult)
	_, err = rpc.Call("bcb_walletCreate", map[string]interface{}{"password": password, "path": path}, result)
	if err != nil {
		fmt.Printf("Cannot create wallet, password=%s, path=%s,\n error=%s \n", password, path, err.Error())
		return nil
	}
	jsIndent, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsIndent))
	return
}
func Transfer(password, path, smcAddress, gasLimit, note, to, value, url string) (err error) {
	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)
	transferParam := rpc3.TransferParam{SmcAddress: smcAddress, GasLimit: gasLimit, Note: note, To: to, Value: value}
	result := new(rpc3.TransferResult)
	_, err = rpc.Call("bcb_transfer", map[string]interface{}{"password": password, "path": path, "walletParams": transferParam}, result)
	if err != nil {
		fmt.Printf("Cannot transfer, password=%s, path=%s, walletParam=%v,\n error=%s \n", password, path, transferParam, err.Error())
		return nil
	}
	jsIndent, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsIndent))
	return
}
func TransferOffline(password, path, smcAddress, gasLimit, note, to, value, nonce, url string) (err error) {
	rpc := rpcclient.NewJSONRPCClientEx(url, "", true)
	uNonce, err := strconv.ParseUint(nonce, 10, 64)
	if err != nil {
		return
	}
	transferParam := rpc3.TransferOfflineParam{SmcAddress: smcAddress, GasLimit: gasLimit, Note: note, Nonce: uNonce, To: to, Value: value}
	result := new(rpc3.TransferOfflineResult)
	_, err = rpc.Call("bcb_transferOffline", map[string]interface{}{"password": password, "path": path, "walletParams": transferParam}, result)
	if err != nil {
		fmt.Printf("Cannot transferOffline, password=%s, path=%s, walletParam=%v,\n error=%s \n", password, path, transferParam, err.Error())
		return nil
	}
	jsIndent, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsIndent))
	return
}