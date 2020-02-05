package rpc

import (
	"blockchain/abciapp_v1.0/keys"
	tx1 "blockchain/abciapp_v1.0/tx/tx"
	"blockchain/abciapp_v1.0/types"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/rlp"
	"blockchain/smcsdk/sdk/std"
	"blockchain/tx2"
	types3 "blockchain/types"
	"common/jsoniter"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	types2 "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/tendermint/rpc/core/types"
	"hdWallet/common"
	"math/big"
	"strings"
)

var genesisTokenAddr = ""

func blockHeight() (blkHeight *BlockHeightResult, err error) {

	result := new(core_types.ResultABCIInfo)
	params := map[string]interface{}{}
	err = common.DoHttpRequestAndParseExBlock(common.GetConfig().NodeAddrSlice, "abci_info", params, result)
	if err != nil {
		return
	}

	blkHeight = new(BlockHeightResult)
	blkHeight.LastBlock = result.Response.LastBlockHeight

	return
}

func block(height int64) (blk *BlockResult, err error) {

	result := new(core_types.ResultBlock)
	params := map[string]interface{}{"height": height}
	err = common.DoHttpRequestAndParseExBlock(common.GetConfig().NodeAddrSlice, "block", params, result)
	if err != nil {
		return
	}

	blk = new(BlockResult)
	blk.BlockHeight = result.BlockMeta.Header.Height
	blk.BlockHash = "0x" + hex.EncodeToString(result.BlockMeta.BlockID.Hash)
	blk.ParentHash = "0x" + hex.EncodeToString(result.BlockMeta.Header.LastBlockID.Hash)
	blk.ChainID = result.BlockMeta.Header.ChainID
	blk.ValidatorHash = "0x" + hex.EncodeToString(result.BlockMeta.Header.ValidatorsHash)
	blk.ConsensusHash = "0x" + hex.EncodeToString(result.BlockMeta.Header.ConsensusHash)
	blk.BlockTime = result.BlockMeta.Header.Time.String()
	blk.BlockSize = result.BlockSize
	blk.ProposerAddress = result.BlockMeta.Header.ProposerAddress

	var blkResults *core_types.ResultBlockResults
	if blkResults, err = blockResults(height); err != nil {
		return
	}

	blk.Txs = make([]TxResult, 0)
	for k, ResDeliver := range blkResults.Results.DeliverTx {
		var tx *TxResult
		if tx, err = transactionBlock(k, ResDeliver, result); err != nil {
			return
		}
		blk.Txs = append(blk.Txs, *tx)
	}

	return
}

func transactionBlock(k int, ResDeliver *types2.ResponseDeliverTx, resultBlock *core_types.ResultBlock) (tx *TxResult, err error) {

	//ParseTX
	var (
		transaction tx1.Transaction
		fromAddr    string
		msg         Message
		GasLimit    uint64
		Nonce       uint64
		Note        string
	)

	messages := make([]Message, 0)
	txStr := string(resultBlock.Block.Txs[k])
	splitTx := strings.Split(txStr, ".")

	if splitTx[1] == "v1" {
		// parse transaction V1
		fromAddr, _, err = transaction.TxParse(crypto.GetChainId(), txStr)
		if err != nil {
			return
		}
		msg, err = messageV1Parse(transaction)
		if err != nil {
			return
		}
		messages = append(messages, msg)
		GasLimit = transaction.GasLimit
		Nonce = transaction.Nonce
		Note = transaction.Note

	} else if splitTx[1] == "v2" {
		// parse transaction V2
		var txv2 types3.Transaction
		var pubKey crypto.PubKeyEd25519
		txv2, pubKey, err := tx2.TxParse(txStr)
		if err != nil {
			return nil, err
		}

		fromAddr = pubKey.Address(crypto.GetChainId())

		var msg Message
		for i := 0; i < len(txv2.Messages); i++ {
			msg, err = messageV2Parse(txv2.Messages[i])
			if err != nil {
				return nil, err
			}
			messages = append(messages, msg)
		}
		GasLimit = uint64(txv2.GasLimit)
		Nonce = txv2.Nonce
		Note = txv2.Note
	} else {
		err = errors.New("unsupported tx=" + txStr)
		return
	}

	tx = new(TxResult)
	tx.TxHash = "0x" + strings.ToLower(hex.EncodeToString(ResDeliver.TxHash))
	tx.TxTime = resultBlock.BlockMeta.Header.Time.String()
	tx.Code = ResDeliver.Code
	tx.Log = ResDeliver.Log
	tx.BlockHash = "0x" + hex.EncodeToString(resultBlock.BlockMeta.BlockID.Hash)
	tx.BlockHeight = resultBlock.BlockMeta.Header.Height
	tx.From = fromAddr
	tx.Nonce = Nonce
	tx.GasLimit = GasLimit
	tx.Fee = ResDeliver.Fee
	tx.Note = Note

	tx.Messages = make([]Message, 0)
	tx.Messages = messages
	tx.TransferReceipts, err = transferReceipts(*ResDeliver)

	return
}

func transaction(txHash string, resultBlock *core_types.ResultBlock) (tx *TxResult, err error) {

	result := new(core_types.ResultTx)
	params := map[string]interface{}{"hash": txHash}
	err = common.DoHttpRequestAndParseExBlock(common.GetConfig().NodeAddrSlice, "tx", params, result)
	if err != nil {
		return
	}

	if resultBlock == nil {
		resultBlock = new(core_types.ResultBlock)
		params = map[string]interface{}{"height": result.Height}
		err = common.DoHttpRequestAndParseExBlock(common.GetConfig().NodeAddrSlice, "block", params, resultBlock)
		if err != nil {
			return
		}
	}

	//ParseTX
	var (
		transaction tx1.Transaction
		fromAddr    string
		msg         Message
		GasLimit    uint64
		Nonce       uint64
		Note        string
		txStr       string
	)

	var blkResults *core_types.ResultBlockResults
	if blkResults, err = blockResults(result.Height); err != nil {
		return
	}

	for k, v := range blkResults.Results.DeliverTx {
		hash := hex.EncodeToString(v.TxHash)
		if hash[:2] == "0x" {
			txHash = txHash[2:]
		}
		if strings.ToLower(txHash) == hash {
			txStr = string(resultBlock.Block.Txs[k])
		}
	}

	messages := make([]Message, 0)

	splitTx := strings.Split(txStr, ".")
	if splitTx[1] == "v1" {
		// parse transaction V1
		fromAddr, _, err = transaction.TxParse(crypto.GetChainId(), txStr)
		if err != nil {
			return
		}
		msg, err = messageV1Parse(transaction)
		if err != nil {
			return
		}
		messages = append(messages, msg)
		GasLimit = transaction.GasLimit
		Nonce = transaction.Nonce
		Note = transaction.Note

	} else if splitTx[1] == "v2" {
		// parse transaction V2
		var txv2 types3.Transaction
		var pubKey crypto.PubKeyEd25519
		txv2, pubKey, err = tx2.TxParse(txStr)
		if err != nil {
			return
		}

		fromAddr = pubKey.Address(crypto.GetChainId())

		var msg Message
		for i := 0; i < len(txv2.Messages); i++ {
			msg, err = messageV2Parse(txv2.Messages[i])
			if err != nil {
				return
			}
			messages = append(messages, msg)
		}
		GasLimit = uint64(txv2.GasLimit)
		Nonce = txv2.Nonce
		Note = txv2.Note
	} else {
		err = errors.New("unsupported tx=" + txStr)
		return
	}

	tx = new(TxResult)
	tx.TxHash = "0x" + txHash
	tx.TxTime = resultBlock.BlockMeta.Header.Time.String()
	tx.Code = result.DeliverResult.Code
	tx.Log = result.DeliverResult.Log
	tx.BlockHash = "0x" + hex.EncodeToString(resultBlock.BlockMeta.BlockID.Hash)
	tx.BlockHeight = resultBlock.BlockMeta.Header.Height
	tx.From = fromAddr
	tx.Nonce = Nonce
	tx.GasLimit = GasLimit
	tx.Fee = result.DeliverResult.Fee
	tx.Note = Note

	tx.Messages = make([]Message, 0)
	tx.Messages = messages
	tx.TransferReceipts, err = transferReceipts(result.DeliverResult)

	return
}

func messageV2Parse(message types3.Message) (msg Message, err error) {

	methodID := fmt.Sprintf("%x", message.MethodID)

	msg.SmcAddress = message.Contract
	if msg.SmcName, msg.Method, err = contractNameAndMethod(message.Contract, methodID); err != nil {
		return
	}

	if len(message.Items) != 2 {
		return msg, errors.New("items count error")
	}

	if methodID == transferMethodIDV2 {

		var to types3.Address
		if err = rlp.DecodeBytes(message.Items[0], &to); err != nil {
			return
		}

		var value bn.Number
		if err = rlp.DecodeBytes(message.Items[1], &value); err != nil {
			return
		}
		msg.To = to
		msg.Value = value.String()
		msg.Method = "Transfer(types.Address,bn.Number)"
	}

	return
}

func balance(address keys.Address) (result *BalanceResult, err error) {

	return balanceOfToken(address, genesisToken(), "")
}

func balanceOfToken(address, tokenAddress keys.Address, tokenName string) (result *BalanceResult, err error) {

	var value []byte
	if tokenName != "" {

		var tmpAddress keys.Address
		param := map[string]interface{}{"path": keyOfTokenName(tokenName)}
		result := new(types.ResultABCIQuery)
		if err = common.DoHttpRequestAndParseExBlock(common.GetConfig().NodeAddrSlice, "abci_query", param, result); err != nil {
			return nil, err
		}

		value = result.Response.Value
		if len(value) == 0 {
			return nil, errors.New("invalid tokenName")
		}

		if err = json.Unmarshal(value, &tmpAddress); err != nil {
			return nil, err
		}

		if tokenAddress != "" && tokenAddress != tmpAddress {
			return nil, errors.New("tokenAddress and tokenName not be same token")
		}
		tokenAddress = tmpAddress
	} else if tokenAddress == "" {
		return nil, errors.New("tokenAddress and tokenName cannot be empty with both")
	}

	NewResult := new(types.ResultABCIQuery)
	param := map[string]interface{}{"path": keyOfAccountToken(address, tokenAddress)}
	if err = common.DoHttpRequestAndParseExBlock(common.GetConfig().NodeAddrSlice, "abci_query", param, &NewResult); err != nil {
		return
	}

	value = NewResult.Response.Value
	result = new(BalanceResult)
	if len(value) == 0 {
		result.Balance = "0"
	} else {
		var tokenBalance types.TokenBalance
		if err = json.Unmarshal(value, &tokenBalance); err != nil {
			return
		}
		result.Balance = tokenBalance.Balance.String()
	}

	return
}

func allBalance(address keys.Address) (items *[]AllBalanceItemResult, err error) {

	tokens := make([]string, 0)
	result := new(types.ResultABCIQuery)
	param := map[string]interface{}{"path": keyOfAccount(address)}
	if err = common.DoHttpRequestAndParseExBlock(common.GetConfig().NodeAddrSlice, "abci_query", param, result); err != nil {
		return
	}

	err = json.Unmarshal(result.Response.Value, &tokens)
	if err != nil {
		return
	}

	balanceItems := make([]AllBalanceItemResult, 0)
	for _, token := range tokens {
		splitToken := strings.Split(token, "/")
		if splitToken[4] != "token" || len(splitToken) != 6 {
			continue
		}

		tokenBalance := new(types.TokenBalance)
		result := new(types.ResultABCIQuery)
		param := map[string]interface{}{"path": token}
		if err = common.DoHttpRequestAndParseExBlock(common.GetConfig().NodeAddrSlice, "abci_query", param, &result); err != nil {
			return
		}

		err = json.Unmarshal(result.Response.Value, tokenBalance)
		if err != nil {
			return
		}

		var name string
		if name, err = tokenName(tokenBalance.Address); err != nil {
			return
		}

		balanceItems = append(balanceItems,
			AllBalanceItemResult{
				TokenAddress: tokenBalance.Address,
				TokenName:    name,
				Balance:      tokenBalance.Balance.String()})
	}

	return &balanceItems, err
}

func nonce(acctAddress keys.Address) (result *NonceResult, err error) {

	type account struct {
		Nonce uint64 `json:"nonce"`
	}

	a := new(account)

	NewResult := new(types.ResultABCIQuery)
	param := map[string]interface{}{"path": keyOfAccountNonce(acctAddress)}
	if err = common.DoHttpRequestAndParseExBlock(common.GetConfig().NodeAddrSlice, "abci_query", param, &NewResult); err != nil {
		return
	}

	value := NewResult.Response.Value
	result = new(NonceResult)
	if len(value) == 0 {
		result.Nonce = 1
	} else {
		err = json.Unmarshal(value, a)
		if err != nil {
			return
		}

		result.Nonce = a.Nonce + 1
	}

	return
}

func commitTx(tx string) (commit *CommitTxResult, err error) {

	result := new(types.ResultBroadcastTxCommit)
	param := map[string]interface{}{"tx": []byte(tx)}
	err = common.DoHttpRequestAndParseExBlock(common.GetConfig().NodeAddrSlice, "broadcast_tx_commit", param, result)
	if err != nil {
		return
	}

	commit = new(CommitTxResult)
	if result.CheckTx.Code != types2.CodeTypeOK {
		commit.Code = result.CheckTx.Code
		commit.Log = result.CheckTx.Log
	} else {
		commit.Code = result.DeliverTx.Code
		commit.Log = result.DeliverTx.Log
	}
	commit.Fee = result.DeliverTx.Fee
	commit.TxHash = "0x" + hex.EncodeToString(result.Hash)
	commit.Height = result.Height

	return
}

func blockResults(height int64) (blkResults *core_types.ResultBlockResults, err error) {

	blkResults = new(core_types.ResultBlockResults)
	params := map[string]interface{}{"height": height}
	err = common.DoHttpRequestAndParseExBlock(common.GetConfig().NodeAddrSlice, "block_results", params, blkResults)
	if err != nil {
		return
	}

	return
}

func messageV1Parse(transation tx1.Transaction) (msg Message, err error) {

	var methodInfo tx1.MethodInfo
	if err = rlp.DecodeBytes(transation.Data, &methodInfo); err != nil {
		return
	}
	methodID := fmt.Sprintf("%x", methodInfo.MethodID)

	msg.SmcAddress = transation.To
	if msg.SmcName, msg.Method, err = contractNameAndMethod(transation.To, methodID); err != nil {
		return
	}

	if methodID == transferMethodIDV1 {
		var itemsBytes = make([][]byte, 0)
		if err = rlp.DecodeBytes(methodInfo.ParamData, &itemsBytes); err != nil {
			return
		}
		msg.To = string(itemsBytes[0])
		msg.Value = new(big.Int).SetBytes(itemsBytes[1][:]).String()
	}

	return
}

func transferReceipts(deliverTx types2.ResponseDeliverTx) ([]std.Transfer, error) {

	receipts := make([]std.Transfer, 0)
	for _, tag := range deliverTx.Tags {
		var receipt std.Receipt
		err := jsoniter.Unmarshal(tag.Value, &receipt)
		if err != nil {
			return nil, err
		}

		if receipt.Name == "std::transfer" || receipt.Name == "transfer" {
			var transferReceipt std.Transfer
			err = jsoniter.Unmarshal(receipt.Bytes, &transferReceipt)
			if err != nil {
				return nil, err
			}

			receipts = append(receipts, transferReceipt)
		}
	}

	return receipts, nil
}

func contractNameAndMethod(contractAddress keys.Address, methodID string) (contractName string, method string, err error) {

	contract := new(types.Contract)
	common.RWLock.RLock()
	v, ok := common.ContractMap[contractAddress]
	common.RWLock.RUnlock()
	if ok == true {
		contract = v
	} else {
		param := map[string]interface{}{"path": keyOfContract(contractAddress)}
		result := new(types.ResultABCIQuery)
		if err = common.DoHttpRequestAndParseExBlock(common.GetConfig().NodeAddrSlice, "abci_query", param, result); err != nil {
			return
		}
		err = json.Unmarshal(result.Response.Value, contract)
		if err != nil {
			return
		}
		common.RWLock.Lock()
		common.ContractMap[contractAddress] = contract
		common.RWLock.Unlock()
	}

	for _, methodItem := range contract.Methods {
		if methodItem.MethodId == methodID {
			method = methodItem.Prototype
			break
		}
	}

	return contract.Name, method, nil
}

func tokenName(tokenAddress keys.Address) (name string, err error) {

	token := new(types.IssueToken)
	param := map[string]interface{}{"path": keyOfToken(tokenAddress)}
	result := new(types.ResultABCIQuery)
	if err = common.DoHttpRequestAndParseExBlock(common.GetConfig().NodeAddrSlice, "abci_query", param, result); err != nil {
		return
	}
	err = json.Unmarshal(result.Response.Value, token)
	if err != nil {
		return
	}

	return token.Name, err
}

func genesisToken() string {
	if genesisTokenAddr == "" {
		token := new(types.IssueToken)
		param := map[string]interface{}{"path": keyOfGenesisToken()}
		result := new(types.ResultABCIQuery)
		if err := common.DoHttpRequestAndParseExBlock(common.GetConfig().NodeAddrSlice, "abci_query", param, result); err != nil {
			return ""
		} else {
			err = json.Unmarshal(result.Response.Value, token)
			genesisTokenAddr = token.Address
		}
	}

	return genesisTokenAddr
}
