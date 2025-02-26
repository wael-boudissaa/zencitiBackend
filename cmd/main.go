package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/wael-boudissaa/marquinoBackend/cmd/api"
	"github.com/wael-boudissaa/marquinoBackend/configs"
	"github.com/wael-boudissaa/marquinoBackend/db"
	"log"
)

func main() {
	cfg := mysql.Config{
		User:                 configs.Env.DBUser,
		Passwd:               configs.Env.DBPassword,
		Addr:                 configs.Env.DBAdress,
		DBName:               configs.Env.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	db, err := db.NewMysqlStorage(cfg)
	if err != nil {
		fmt.Println(err)
	}
	initStorage(db)

	server := api.NewApiServer(fmt.Sprintf(":%s", configs.Env.Port), db)

	if err := server.Run(); err != nil {
		fmt.Println(err)
	} else {
		log.Fatal("server running")

	}

}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("database connected")
}
