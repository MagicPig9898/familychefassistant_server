package db_config

import (
	"github.com/MagicPig9898/easy_db/mysql"
	conf "github.com/MagicPig9898/familychefassistant_server/conf"
)

type dbconfig struct {
	mclient *mysql.Client
}

var dbmgr *dbconfig

func NewDbConfig() error {
	c := conf.Cfg.DB
	mclient, err := mysql.NewClient(c.Host, c.Port, c.User, c.Password, c.DBName)
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
