package main

import (
	rpcserver "common/rpc/lib/server"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	cmn "github.com/tendermint/tmlibs/common"
	"hdWallet/common"
	"hdWallet/rpc"
	"net/http"
	"os"
)

func main() {
	err := common.InitAll()
	if err != nil {
		panic(err)
	}

	err = rpc.InitDB()
	if err != nil {
		common.GetLogger().Error("open db failed", "error", err.Error())
		panic(err)
	}

	if len(os.Args) == 1 {

		rpcLogger := common.GetLogger()

		coreCodec := amino.NewCodec()

		mux := http.NewServeMux()
		rpcserver.NoLog = true
		rpcserver.RegisterRPCFuncs(mux, rpc.Routes, coreCodec, rpcLogger)

		if common.GetConfig().UseHttps {
			crtPath, keyPath := common.OutCertFileIsExist()
			_, err = rpcserver.StartHTTPAndTLSServer(common.GetConfig().ServerAddr, mux, crtPath, keyPath, rpcLogger)
			if err != nil {
				cmn.Exit(err.Error())
			}
		} else {
			_, err = rpcserver.StartHTTPServer(common.GetConfig().ServerAddr, mux, rpcLogger)
			if err != nil {
				cmn.Exit(err.Error())
			}
		}
		// Wait forever
		cmn.TrapSignal(func(signal os.Signal) {
		})
	} else {
		err = Execute()
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	}

}

// flags
var (
	// global flags
	flagPassword string

	//HD wallet flag
	flagMnemonic string
)

var RootCmd = &cobra.Command{
	Use:   "hdWallet_rpc",
	Short: "bcb exchange wallet console",
	Long:  "hdWallet_rpc client that it can perform the wallet operation, query chain information and so on.",
}

func Execute() error {
	addFlags()
	addCommands()
	return RootCmd.Execute()
}

func addFlags() {
	addExportMnemonicFlag()
	addImportMnemonicFlag()
	addChangePasswordFlag()
}

func addCommands() {
	RootCmd.AddCommand(createMnemonicCmd)
	RootCmd.AddCommand(exportMnemonicCmd)
	RootCmd.AddCommand(importMnemonicCmd)
	RootCmd.AddCommand(changePasswordCmd)
}

var createMnemonicCmd = &cobra.Command{
	Use:   "createMnemonic",
	Short: "Create mnemonic",
	Long:  "Create a new mnemonic",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return rpc.CreateMnemonic()
	},
}

var exportMnemonicCmd = &cobra.Command{
	Use:   "exportMnemonic",
	Short: "Export mnemonic",
	Long:  "Export a Mmnemonic",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return rpc.ExportMnemonic(flagPassword)
	},
}

func addExportMnemonicFlag() {
	exportMnemonicCmd.PersistentFlags().StringVarP(&flagPassword, "password", "p", "", "mnemonic password")
}

var importMnemonicCmd = &cobra.Command{
	Use:   "importMnemonic",
	Short: "Import mnemonic",
	Long:  "Import a original mnemonic",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return rpc.ImportMnemonic(flagMnemonic)
	},
}

func addImportMnemonicFlag() {
	importMnemonicCmd.PersistentFlags().StringVarP(&flagMnemonic, "mnemonic", "m", "", "12 mnemonics")
}

var changePasswordCmd = &cobra.Command{
	Use:   "changePassword",
	Short: "Change password",
	Long:  "Change wallet password",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return rpc.ChangePassword(flagPassword)
	},
}

func addChangePasswordFlag() {
	changePasswordCmd.PersistentFlags().StringVarP(&flagPassword, "password", "p", "", "old account password")
}
