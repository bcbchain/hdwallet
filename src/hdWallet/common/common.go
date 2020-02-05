package common

import (
	"blockchain/abciapp_v1.0/types"
	"blockchain/tx2"
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tendermint/go-crypto"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tmlibs/log"
	"hdWallet/common/config"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	hdWalletConfig  config.Config
	logger          log.Logger
	FaultCounterMap = make(map[string]int)
	CorrectUrls     = make([]string, 0)
	ContractMap     = make(map[string]*types.Contract)
	RWLock          = new(sync.RWMutex)
)

func InitAll() error {
	configFile := "./.config/hdWallet.yaml"
	moduleName := "hdWallet"

	err := hdWalletConfig.InitConfig(configFile)
	if err != nil {
		return errors.New("Init config fail err info : " + err.Error())
	}
	initLog(moduleName)

	if hdWalletConfig.ChainID == "" {
		return errors.New(" chainId cannot be empty")
	}
	crypto.SetChainId(hdWalletConfig.ChainID)
	tx2.Init(hdWalletConfig.ChainID)

	CheckChainVersion()

	for _, v := range GetConfig().NodeAddrSlice {
		FaultCounterMap[v] = 0
		CorrectUrls = append(CorrectUrls, v)
	}

	rand.Seed(time.Now().Unix())

	return nil
}

func initLog(moduleName string) {
	l := log.NewTMLogger("./log", moduleName)
	l.SetOutputToFile(hdWalletConfig.LoggerFile)
	l.SetOutputToScreen(hdWalletConfig.LoggerScreen)
	l.AllowLevel(hdWalletConfig.LoggerLevel)
	logger = l
}

func GetConfig() config.Config {
	return hdWalletConfig
}

func GetLogger() log.Logger {
	return logger
}

func FuncRecover(l log.Logger, errPtr *error) {
	if err := recover(); err != nil {
		msg := ""
		if errInfo, ok := err.(error); ok {
			msg = errInfo.Error()
		}

		if errInfo, ok := err.(string); ok {
			msg = errInfo
		}

		l.Error("FuncRecover", "error", msg)
		*errPtr = errors.New(msg)
	}
}

func OutCertFileIsExist() (string, string) {
	crtPath := "./.config/server.crt"
	keyPath := "./.config/server.key"

	_, err := os.Stat(hdWalletConfig.OutCertPath + ".crt")
	if err != nil {
		return crtPath, keyPath
	}

	_, err = os.Stat(hdWalletConfig.OutCertPath + ".key")
	if err != nil {
		return crtPath, keyPath
	}

	return hdWalletConfig.OutCertPath + ".crt", hdWalletConfig.OutCertPath + ".key"
}

func CheckChainVersion() {
	cfg := hdWalletConfig

	if cfg.ChainVersion != "1" && cfg.ChainVersion != "2" && cfg.ChainVersion != "" {
		fmt.Println("Config file error, please check chainVersion!")
		return
	}

	if cfg.ChainVersion == "2" {
		return
	}

	ChainVersion, err := queryChainVersion()
	if err != nil {
		fmt.Println("Query ChainVersion failed, please check!")
		return
	}

	if ChainVersion == "0" {
		ChainVersion = "1"
	}

	if cfg.ChainVersion != ChainVersion {
		changeChainVersion(ChainVersion)
		hdWalletConfig.ChainVersion = ChainVersion
	}
}

func queryChainVersion() (chainVersion string, err error) {
	result := new(core_types.ResultHealth)
	params := map[string]interface{}{}
	err = DoHttpRequestAndParseExBlock(GetConfig().NodeAddrSlice, "health", params, result)
	if err != nil {
		return "", err
	}

	chainVersion = strconv.FormatInt(result.ChainVersion, 10)
	return
}

func changeChainVersion(chainversion string) {
	configFile := "./.config/hdWallet.yaml"

	f, err := os.Open(configFile)
	if err != nil {
		fmt.Println("OpenFile failed, please check!")
		return
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	var Str string
	for {
		line, err := buf.ReadString('\n')
		if strings.HasPrefix(line, "chainVersion:") {
			line = "chainVersion: " + chainversion + "\n"
		}
		Str = Str + line

		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
	}

	file2, err := os.Create(configFile)
	if err != nil {
		fmt.Println("CreateFile failed, please check!")
		return
	}
	defer file2.Close()

	_, err = file2.WriteString(Str)
	if err != nil {
		fmt.Println("WriteFile failed, please check!")
		return
	}
}
