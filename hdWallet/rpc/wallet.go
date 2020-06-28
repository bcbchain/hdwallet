package rpc

import (
	"encoding/hex"
	"errors"
	"strings"

	types2 "github.com/bcbchain/bcbchain/abciapp_v1.0/types"
	"github.com/bcbchain/bclib/bn"
	"github.com/bcbchain/bclib/hdwal"
	"github.com/bcbchain/bclib/tendermint/go-crypto"
	tx2 "github.com/bcbchain/bclib/tx/v2"
	"github.com/bcbchain/bclib/types"
	"github.com/btcsuite/btcutil/base58"
	"hdwallet/hdWallet/common"
)

const (
	pattern = "^[a-zA-Z0-9_@.-]{1,40}$"
)

func transfer(password, path string, gasLimit uint64, walletParams TransferParam) (result *TransferResult, err error) {

	config := common.GetConfig()
	result = new(TransferResult)

	var txStr string
	var nonceResult *NonceResult
	var accPrikeyBytes crypto.PrivKeyEd25519

	if !strings.HasPrefix(path, "m/44'/60'/0'/0/") {
		err = errors.New("path must start with m/44'/60'/0'/0/ ")
		return
	}

	cfg := common.GetConfig()
	var accPrikey crypto.PrivKey
	accPrikey, err = hdwal.NewPrivKey(path, password)
	if err != nil {
		return nil, errors.New("NewPrivKey wrong")
	}
	accPrikeyBytes = accPrikey.(crypto.PrivKeyEd25519)
	accountAddr := accPrikey.PubKey().Address(cfg.ChainID)

	// 获取nonce
	nonceResult, err = nonce(accountAddr)
	if err != nil {
		return
	}

	value := bn.NewNumberStringBase(walletParams.Value, 10)

	if config.ChainVersion == "1" {
		txStr, err = PackAndSignTx(nonceResult.Nonce, gasLimit, walletParams.Note, walletParams.SmcAddress, walletParams.To, value.Bytes(), accPrikeyBytes[:])
		if err != nil {
			return nil, err
		}
	} else if config.ChainVersion == "2" {
		var method uint32 = 0x44D8CA60
		v := bn.NewNumberStringBase(walletParams.Value, 10)
		V2Paramss := []interface{}{walletParams.To, v}
		prikey := "0x" + hex.EncodeToString(accPrikeyBytes[:])

		txStr = GenerateTx(walletParams.SmcAddress, method, V2Paramss, nonceResult.Nonce, int64(gasLimit), walletParams.Note, prikey)
	} else {
		return nil, errors.New("ChainVersion wrong, please check!")
	}

	commitResult := new(types2.ResultBroadcastTxCommit)
	param := map[string]interface{}{"tx": []byte(txStr)}
	err = common.DoHttpRequestAndParseExBlock(config.NodeAddrSlice, "broadcast_tx_commit", param, commitResult)
	if err != nil {
		return
	}

	if commitResult.CheckTx.Code != 200 {
		result.Log = commitResult.CheckTx.Log
		result.Code = commitResult.CheckTx.Code
	} else {
		result.Log = commitResult.DeliverTx.Log
		result.Code = commitResult.DeliverTx.Code
	}
	result.Fee = commitResult.DeliverTx.Fee
	result.Height = commitResult.Height
	result.TxHash = "0x" + hex.EncodeToString(commitResult.Hash)

	return
}

func walletTransferOffline(password, path string, gasLimit uint64, walletParams TransferOfflineParam) (result *TransferOfflineResult, err error) {

	config := common.GetConfig()
	value := bn.NewNumberStringBase(walletParams.Value, 10)

	var txStr string
	var accPrikeyBytes crypto.PrivKeyEd25519

	if !strings.HasPrefix(path, "m/44'/60'/0'/0/") {
		err = errors.New("path must start with m/44'/60'/0'/0/ ")
		return
	}

	var accPrikey crypto.PrivKey
	accPrikey, err = hdwal.NewPrivKey(path, password)
	if err != nil {
		return nil, errors.New("NewPrivKey wrong")
	}
	accPrikeyBytes = accPrikey.(crypto.PrivKeyEd25519)

	if config.ChainVersion == "1" {
		txStr, err = PackAndSignTx(walletParams.Nonce, gasLimit, walletParams.Note, walletParams.SmcAddress, walletParams.To, value.Bytes(), accPrikeyBytes[:])
		if err != nil {
			return nil, err
		}
	} else if config.ChainVersion == "2" {
		var method uint32 = 0x44D8CA60
		v := bn.NewNumberStringBase(walletParams.Value, 10)
		V2Paramss := []interface{}{walletParams.To, v}
		prikey := "0x" + hex.EncodeToString(accPrikeyBytes[:])

		txStr = GenerateTx(walletParams.SmcAddress, method, V2Paramss, walletParams.Nonce, int64(gasLimit), walletParams.Note, prikey)
	} else {
		return nil, errors.New("ChainVersion wrong, please check!")
	}

	result = new(TransferOfflineResult)
	result.Tx = txStr

	return
}

//GenerateTx generate tx with one contract method request
func GenerateTx(contract types.Address, method uint32, V2Paramss []interface{}, nonce uint64, gaslimit int64, note string, privKey string) string {
	items := tx2.WrapInvokeParams(V2Paramss...)
	message := types.Message{
		Contract: contract,
		MethodID: method,
		Items:    items,
	}
	payload := tx2.WrapPayload(nonce, gaslimit, note, message)
	return tx2.WrapTx(payload, privKey)
}

// HD wallet
func createMnemonic() (result *CreateMnemonicResult, err error) {
	logger := common.GetLogger()
	logger.Info("new Mnemonic")

	password := base58.Encode(crypto.CRandBytes(32))

	mnemonic, err := hdwal.Mnemonic(password)
	if err != nil {
		logger.Error("Create mnemonic failed")
		return
	}

	result = new(CreateMnemonicResult)
	result.Mnemonic = mnemonic
	result.Password = password

	return
}

func exportMnemonic(password string) (result *ExportMnemonicResult, err error) {
	logger := common.GetLogger()
	logger.Info("exportMnemonic", "password", PASSWORDSTAR)

	mnemonic, err := hdwal.Export(password)
	if err != nil {
		logger.Error("Export wallet failed")
		return
	}

	result = new(ExportMnemonicResult)
	result.Mnemonic = mnemonic

	return
}

func importMnemonic(mnemonic string) (result *ImportMnemonicResult, err error) {
	logger := common.GetLogger()
	logger.Info("importMnemonic", "mnemonic", mnemonic)

	password := base58.Encode(crypto.CRandBytes(32))

	err = hdwal.Import(mnemonic, password)
	if err != nil {
		logger.Error("Import mnemonic failed")
		return
	}

	result = new(ImportMnemonicResult)
	result.Password = password

	return
}

func changePassword(password string) (result *ChangePasswordResult, err error) {
	logger := common.GetLogger()
	logger.Info("changePassword", "password", PASSWORDSTAR)

	mnemonic, err := hdwal.Export(password)
	if err != nil {
		logger.Error("Export wallet failed")
		return
	}

	logger.Info("Save mnemonic", "newPassword", PASSWORDSTAR)
	newPassword := base58.Encode(crypto.CRandBytes(32))
	err = hdwal.Import(mnemonic, newPassword)
	if err != nil {
		logger.Error("Import mnemonic failed")
		return
	}

	result = new(ChangePasswordResult)
	result.Password = newPassword

	return
}

func createAccount(password, path string) (result *WalletCreateResult, err error) {
	cfg := common.GetConfig()
	logger := common.GetLogger()
	logger.Info("createAccount", "password", PASSWORDSTAR, "path", path)

	prikey, err := hdwal.NewPrivKey(path, password)
	if err != nil {
		logger.Error("Create Account failed")
		return
	}

	accountAddr := prikey.PubKey().Address(cfg.ChainID)

	result = new(WalletCreateResult)
	result.WalletAddress = accountAddr

	return
}
