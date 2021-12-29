package main

import (
	"database/sql"
	"fmt"
	"git.fputs.com/fputs/kccloot/pkg/raiders"
	"git.fputs.com/fputs/kccloot/pkg/util"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db               *sql.DB
	connectionString string
)

func init() {
	connectionString = os.Getenv("KCC_CONNSTR")
	if connectionString == "" {
		panic("KCC_CONNSTR is unset")
	}

	var err error
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

func main() {
	defer db.Close()

	statement, err := db.Prepare("select * from raiders")
	util.CheckErr(err)

	rows, err := statement.Query()
	util.CheckErr(err)
	defer rows.Close()

	var rs []raiders.Raider
	for rows.Next() {
		var r raiders.Raider
		err = rows.Scan(&r.Id, &r.Name, &r.Points)
		util.CheckErr(err)
		rs = append(rs, r)
	}

	for _, r := range rs {
		fmt.Printf("%s: %d\n", r.Name, r.Points)
	}
}
