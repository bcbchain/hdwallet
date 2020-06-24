package rpc

import (
	"bufio"
	"github.com/bgentry/speakeasy"
)

func getHDWalletPassword(prompt string, buf *bufio.Reader) (pass string, err error) {
	pass, err = speakeasy.Ask(prompt)
	if err != nil {
		return "", err
	}

	return pass, nil
}
