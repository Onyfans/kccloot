package main

import (
	"database/sql"
	"fmt"
	"git.fputs.com/fputs/kccloot/pkg/raiders"
	"git.fputs.com/fputs/kccloot/pkg/util"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
	"time"
)

var connectionString string

func main() {
	connectionString = os.Getenv("KCCWEB_CONNSTR")
	if connectionString == "" {
		panic("KCCWEB_CONNSTR is unset")
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8012", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

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
		fmt.Fprintf(w, "%s: %d\n", r.Name, r.Points)
	}
}
