package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"os"
	"time"

	"git.fputs.com/fputs/kccloot/pkg/raiders"
	"git.fputs.com/fputs/kccloot/pkg/util"
	_ "github.com/go-sql-driver/mysql"
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
		err = rows.Scan(&r.Id, &r.Name, &r.Points, &r.Class, &r.Spec)
		util.CheckErr(err)
		rs = append(rs, r)
	}
	raiders.SortSlice(rs)

	const tmpl = `
	{{ $length := len . }}
	<head>
		<link rel="stylesheet" href="https://fputs.com/webdisplay.css">
		<title>KCC Loot Tracker</title>
	</head>
	<h1>KCC Loot Tracker</h1>
	<table>
	<tr>
		<th>Name</th>
		<th>Class</th>
		<th>Spec</th>
		<th>Points</th>
	</tr>
		{{range .}}
			<tr><td>{{.Name}}</td><td>{{.Class}}</td><td>{{.Spec}}</td><td>{{.Points}}</td></tr>
		{{end}}
	</table>
	<p>Total Raiders: {{ $length }}</p>
	`
	t := template.Must(template.New("").Parse(tmpl))
	if err := t.Execute(w, rs); err != nil {
		util.CheckErr(err)
	}
}
