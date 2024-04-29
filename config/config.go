package config

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func Dbconnect() (*sql.DB, error) {
	dburi := fmt.Sprintf("%v:%v@tcp(%v)/%v",
		GetConfigValues().DBConfig.Username,
		GetConfigValues().DBConfig.Password,
		GetConfigValues().DBConfig.Server,
		GetConfigValues().DBConfig.Schema)
	engine := GetConfigValues().DBConfig.Engine
	db, err := sql.Open(engine, dburi)
	if err != nil {
		logger.Error(err.Error())
	
	}
	return db, nil
}
