package db_config

import "github.com/MagicPig9898/easy_db/mysql"

type dbconfig struct {
	mclient *mysql.Client
}

var dbmgr *dbconfig

func NewDbConfig() error {
	mclient, err := mysql.NewClient("localhost", 3306, "root", "123456", "mysql")
	if err != nil {
		return err
	}
	dbmgr = &dbconfig{
		mclient: mclient,
	}
	return nil
}

func GetDb() *mysql.Client {
	if dbmgr == nil {
		return nil
	}
	return dbmgr.mclient
}

func Close() error {
	if dbmgr == nil {
		return nil
	}
	return dbmgr.mclient.Close()
}
