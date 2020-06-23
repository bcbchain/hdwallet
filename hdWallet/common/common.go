package common

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bcbchain/bcbchain/abciapp_v1.0/types"
	"github.com/bcbchain/bclib/tendermint/go-crypto"
	"github.com/bcbchain/bclib/tendermint/tmlibs/log"
	tx2 "github.com/bcbchain/bclib/tx/v2"
	"github.com/bcbchain/sdk/sdk/std"
	core_types "github.com/bcbchain/tendermint/rpc/core/types"
	"github.com/pkg/errors"
	"hdwallet/hdWallet/common/config"
)

var (
	hdWalletConfig  config.Config
	logger          log.Logger
	savePathFile    = "./configPath.conf"
	FaultCounterMap = make(map[string]int)
	CorrectUrls     = make([]string, 0)
	ContractMap     = make(map[string]*types.Contract)
	BVMContractMap  = make(map[string]*std.BvmContract)
	RWLock          = new(sync.RWMutex)
)

func InitAll() error {
	configFile := filepath.Join(getConfigPath(), "hdWallet.yaml")
	//configFile := "./.config/hdWallet.yaml"
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

func SetConfigPath(newPath string) error {
	fi, err := os.Create(savePathFile)
	if err != nil {
		return err
	}
	defer fi.Close()

	_, err = fi.Write([]byte(newPath))
	if err != nil {
		return err
	}

	return nil
}

func getConfigPath() string {
	_, err := os.Stat(savePathFile)
	if err != nil && !os.IsExist(err) {
		return "./.config"
	}

	content, err := ioutil.ReadFile(savePathFile)
	if err != nil {
		return "./.config"
	}

	return string(content)
}

func initLog(moduleName string) {
	l := log.NewTMLogger("./log", moduleName)
	l.SetOutputToFile(hdWalletConfig.LoggerFile)
	l.SetOutputToScreen(hdWalletConfig.LoggerScreen)
	l.AllowLevel(hdWalletConfig.LoggerLevel)
	logger = l
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
