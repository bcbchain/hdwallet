package rpc

import (
	"blockchain/abciapp_v1.0/keys"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tendermint/go-crypto"
	"hdWallet/common"
	"io/ioutil"
	"os"
	"strings"
)

var (
	PASSWORDSTAR = "********"
)

// WalletCreate - create wallet
func WalletCreate(password, path string) (result *WalletCreateResult, err error) {
	logger := common.GetLogger()
	logger.Info("WalletCreate begin", "password", PASSWORDSTAR, "path", path)
	defer common.FuncRecover(logger, &err)

	if len(password) == 0 {
		buf := bufio.NewReader(os.Stdin)
		password, err = getHDWalletPassword("Enter Password("+path+"):", buf)
		if err != nil {
			return
		}
	}

	result = new(WalletCreateResult)

	if !strings.HasPrefix(path, "m/44'/60'/0'/0/") {
		err = errors.New("path must start with m/44'/60'/0'/0/ ")
		return
	}

	result, err = createAccount(password, path)
	if err != nil {
		logger.Error("Cannot create account", "error", err)
		return
	}

	logger.Info("WalletCreate finish", "path", path, "password", PASSWORDSTAR)
	return
}

// WalletTransfer - transfer token
func WalletTransfer(password, path string, walletParams TransferParam) (result *TransferResult, err error) {
	logger := common.GetLogger()

	defer common.FuncRecover(logger, &err)
	logger.Trace("bcb_transfer", "path", path, "gasLimit", walletParams.GasLimit, "note", walletParams.Note, "to", walletParams.To, "Value", walletParams.Value)

	if len(password) == 0 {
		buf := bufio.NewReader(os.Stdin)
		password, err = getHDWalletPassword("Enter Password("+path+"):", buf)
		if err != nil {
			return
		}
	}

	//parse gasLimit
	gasLimit, err := requireUint64(walletParams.GasLimit)
	if err != nil {
		return
	}

	// check value
	if _, err = requireUint64(walletParams.Value); err != nil {
		return
	}

	// check smcAddress
	if err = checkAddress(crypto.GetChainId(), walletParams.SmcAddress); err != nil {
		return
	}

	// check to address
	if err = checkAddress(crypto.GetChainId(), walletParams.To); err != nil {
		return
	}

	result, err = transfer(password, path, gasLimit, walletParams)
	if err != nil {
		logger.Error("Cannot transfer", "error", err)
	}

	return
}

// WalletTransferOffline - pack transfer transaction offline
func WalletTransferOffline(password, path string, walletParams TransferOfflineParam) (result *TransferOfflineResult, err error) {
	logger := common.GetLogger()

	defer common.FuncRecover(logger, &err)
	logger.Trace("bcb_transferOffline", "path", path, "gasLimit", walletParams.GasLimit, "note", walletParams.Note, "to", walletParams.To, "Value", walletParams.Value)

	if len(password) == 0 {
		buf := bufio.NewReader(os.Stdin)
		password, err = getHDWalletPassword("Enter Password("+path+"):", buf)
		if err != nil {
			return
		}
	}

	//parse gasLimit
	gasLimit, err := requireUint64(walletParams.GasLimit)
	if err != nil {
		return
	}

	// check value
	if _, err = requireUint64(walletParams.Value); err != nil {
		return
	}

	// check smcAddress
	if err = checkAddress(crypto.GetChainId(), walletParams.SmcAddress); err != nil {
		return
	}

	// check to address
	if err = checkAddress(crypto.GetChainId(), walletParams.To); err != nil {
		return
	}

	result, err = walletTransferOffline(password, path, gasLimit, walletParams)
	if err != nil {
		logger.Error("Cannot pack transfer transaction", "error", err)
	}

	return
}

// BlockHeight - get current block height
func BlockHeight() (result *BlockHeightResult, err error) {
	defer common.FuncRecover(common.GetLogger(), &err)

	common.GetLogger().Trace("bcb_blockHeight")

	result, err = blockHeight()
	if err != nil {
		common.GetLogger().Error("Cannot get current block height", "error", err)
	}

	return
}

// Block - get block data with height
func Block(height int64) (result *BlockResult, err error) {
	defer common.FuncRecover(common.GetLogger(), &err)

	common.GetLogger().Trace("bcb_block", "height", height)

	// if height is 0, set it current height
	if height == 0 {
		var blkHeight *BlockHeightResult
		if blkHeight, err = blockHeight(); err != nil {
			common.GetLogger().Error("Cannot get current block height", "error", err)
			return
		}
		height = blkHeight.LastBlock
	} else if height < 0 {
		return nil, errors.New("Height cannot be negative ")
	}

	result, err = block(height)
	if err != nil {
		common.GetLogger().Error("Cannot get block data", "height", height, "error", err)
	}

	return
}

// Transaction - get transaction data with txHash
func Transaction(txHash string) (result *TxResult, err error) {
	defer common.FuncRecover(common.GetLogger(), &err)

	common.GetLogger().Trace("bcb_transaction", "txHash", txHash)

	if txHash == "" {
		return nil, errors.New("TxHash cannot be empty ")
	}
	if txHash[:2] == "0x" {
		txHash = txHash[2:]
	}

	result, err = transaction(txHash, nil)
	if err != nil {
		common.GetLogger().Error("Cannot get transaction data", "error", err)
	}

	return
}

// Balance - get balance of account address
func Balance(address keys.Address) (result *BalanceResult, err error) {
	defer common.FuncRecover(common.GetLogger(), &err)

	common.GetLogger().Trace("bcb_balance", "address", address)

	if address == "" {
		return nil, errors.New("Address cannot be empty ")
	}

	if err = checkAddress(crypto.GetChainId(), address); err != nil {
		return
	}

	result, err = balance(address)
	if err != nil {
		common.GetLogger().Error("Cannot get balance", "error", err)
	}

	return
}

// BalanceOfToken - get balance of account address and token address
func BalanceOfToken(address, tokenAddress keys.Address, tokenName string) (result *BalanceResult, err error) {
	defer common.FuncRecover(common.GetLogger(), &err)

	common.GetLogger().Trace("bcb_balanceOfToken", "address", address, "tokenAddress", tokenAddress, "tokenName", tokenName)

	if address == "" {
		return nil, errors.New("Address cannot be empty ")
	}

	if err = checkAddress(crypto.GetChainId(), address); err != nil {
		return
	}

	if tokenAddress != "" {
		if err = checkAddress(crypto.GetChainId(), tokenAddress); err != nil {
			return
		}
	} else if tokenName == "" {
		return nil, errors.New("TokenAddress and TokenName cannot empty with both ")
	}

	result, err = balanceOfToken(address, tokenAddress, tokenName)
	if err != nil {
		common.GetLogger().Error("Cannot get balance of token", "error", err)
	}

	return
}

// AllBalance - get all token balance of account address
func AllBalance(address keys.Address) (result *[]AllBalanceItemResult, err error) {
	defer common.FuncRecover(common.GetLogger(), &err)

	common.GetLogger().Trace("bcb_allBalance", "address", address)

	if err = checkAddress(crypto.GetChainId(), address); err != nil {
		return
	}

	result, err = allBalance(address)
	if err != nil {
		common.GetLogger().Error("Cannot get all balance", "error", err)
	}

	return
}

// Nonce - get nonce of account address
func Nonce(address keys.Address) (result *NonceResult, err error) {
	defer common.FuncRecover(common.GetLogger(), &err)

	common.GetLogger().Trace("bcb_nonce", "address", address)

	if err = checkAddress(crypto.GetChainId(), address); err != nil {
		return
	}

	result, err = nonce(address)
	if err != nil {
		common.GetLogger().Error("Cannot get nonce", "error", err)
	}

	return
}

// CommitTx - commit transaction
func CommitTx(tx string) (result *CommitTxResult, err error) {
	defer common.FuncRecover(common.GetLogger(), &err)

	common.GetLogger().Trace("bcb_commitTx", "tx", tx)

	if tx == "" {
		return nil, errors.New("Tx cannot be empty ")
	}

	result, err = commitTx(tx)
	if err != nil {
		common.GetLogger().Error("Cannot commit tx", "error", err)
	}

	return
}

// Version - return current app version
func Version() (result *VersionResult, err error) {
	defer common.FuncRecover(common.GetLogger(), &err)

	common.GetLogger().Trace("bcb_version")

	var version []byte
	version, err = ioutil.ReadFile("./.config/version")
	if err != nil {
		common.GetLogger().Error("Read version file error", "error", err)
		return
	}
	result = new(VersionResult)
	NewVersion := string(version)
	NewVersion = strings.Replace(NewVersion, "\r\n", "", -1)
	NewVersion = strings.Replace(NewVersion, "\n", "", -1)
	result.Version = NewVersion

	return
}

// CreateMnemonic - create mnemonic
func CreateMnemonic() (err error) {
	logger := common.GetLogger()
	logger.Info("CreateMnemonic begin")
	defer common.FuncRecover(logger, &err)

	result, err := createMnemonic()
	if err != nil {
		logger.Error("Cannot create mnemonic", "error", err)
		return
	}

	logger.Info("CreateMnemonic finish")
	jsIndent, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsIndent))

	return
}

// ExportMnemonic - export mnemonic
func ExportMnemonic(password string) (err error) {
	logger := common.GetLogger()
	logger.Info("ExportWallet begin", "password", PASSWORDSTAR)
	defer common.FuncRecover(logger, &err)

	if len(password) == 0 {
		buf := bufio.NewReader(os.Stdin)
		password, err = getHDWalletPassword("Enter Password:", buf)
		if err != nil {
			return
		}
	}

	result, err := exportMnemonic(password)
	if err != nil {
		logger.Error("Cannot export mnemonic", "error", err)
		return
	}

	logger.Info("ExportMnemonic finish", "password", PASSWORDSTAR)
	jsIndent, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsIndent))

	return
}

// ImportMnemonic - import mnemonic
func ImportMnemonic(mnemonic string) (err error) {
	logger := common.GetLogger()
	logger.Info("ImportMnemonic begin", "mnemonic", mnemonic)
	defer common.FuncRecover(logger, &err)

	mnemonic = strings.TrimSpace(mnemonic)
	if len(strings.Split(mnemonic, " ")) != 12 {
		logger.Error("Incorrect mnemonic format, the length of mnemonics must be 12 and only a space between each mnemonic")
		err = errors.New("Incorrect mnemonic format, the length of mnemonics must be 12 and only a space between each mnemonic")
		return
	}

	result, err := importMnemonic(mnemonic)
	if err != nil {
		logger.Error("Cannot import mnemonic", "error", err)
		return
	}

	logger.Info("ImportMnemonic finish", "mnemonic", mnemonic)
	jsIndent, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsIndent))

	return
}

// ChangePassword - change password
func ChangePassword(password string) (err error) {
	logger := common.GetLogger()
	logger.Info("ChangePassword begin", "password", PASSWORDSTAR)
	defer common.FuncRecover(logger, &err)

	if len(password) == 0 {
		buf := bufio.NewReader(os.Stdin)
		password, err = getHDWalletPassword("Enter password:", buf)
		if err != nil {
			return
		}
	}

	result, err := changePassword(password)
	if err != nil {
		logger.Error("Cannot change password", "error", err)
		return
	}

	logger.Info("ChangePassword finish", "password", PASSWORDSTAR)
	jsIndent, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsIndent))
	return
}
