package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"mail-provider/config"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("mysql", config.Config().Database)
	if err != nil {
		log.Fatalln("open db fail:", err)
	}

	DB.SetMaxOpenConns(config.Config().MaxConns)
	DB.SetMaxIdleConns(config.Config().MaxIdle)

	err = DB.Ping()
	if err != nil {
		log.Fatalln("ping db fail:", err)
	}
}
