package rpc

import (
	"bytes"
	"errors"
	"strconv"
	"strings"

	"github.com/bcbchain/bcbchain/abciapp_v1.0/smc"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

// nolint
func checkAddress(chainID string, addr smc.Address) error {
	if !strings.HasPrefix(addr, chainID) {
		return errors.New("Address chainID is error! ")
	}
	base58Addr := strings.Replace(addr, chainID, "", 1)
	addrData := base58.Decode(base58Addr)
	dataLen := len(addrData)
	if dataLen < 4 {
		return errors.New("Base58Addr parse error! ")
	}

	hasher := ripemd160.New()
	hasher.Write(addrData[:dataLen-4])
	md := hasher.Sum(nil)

	if bytes.Compare(md[:4], addrData[dataLen-4:]) != 0 {
		return errors.New("Address checksum is error! ")
	}

	return nil
}

func requireUint64(valueStr string) (uint64, error) {
	value, err := strconv.ParseUint(valueStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}
