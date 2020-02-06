package rpc

import (
	"common/bcdb"
	"common/hdwal"
	"hdWallet/common"
	"path/filepath"
)

type DB struct {
	*bcdb.GILevelDB
}
var (
	db     DB
	dbName = "account"
)
// Init DB
func InitDB() error {
	var err error
	dbPath := absolutePath(common.GetConfig().KeyStorePath)
	db.GILevelDB, err = bcdb.OpenDB(dbPath, "", "")
	hdwal.SetDB(db.GILevelDB)
	return err
}
func absolutePath(path string) string {
	if filepath.IsAbs(path) {
		path = filepath.Join(path, dbName)
	} else {
		dir, err := common.CurrentDirectory()
		if err != nil {
			panic(err)
		}
		path = filepath.Join(dir, path, dbName)
	}
	return path
}