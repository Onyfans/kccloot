package main

import (
	"database/sql"
	"git.fputs.com/fputs/kccloot/pkg/raiders"
	"git.fputs.com/fputs/kccloot/pkg/util"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
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
	</br>
	<h2>System Rules</h2>
	<ul>
	<li>Every raider who shows up to the raid on time and is ready (Raid Team Rules #1 and #2) receives 1 point</li>
	<li>Every raider that participates in a kill receives 1 point</li>
	<li>If a raider participates in wipes on a boss but has to leave before the kill, they will still receive a point if a kill occurs</li>
	<li>For each drop, the raid leader will ask who wants the loot. The raider with the most points will receive the item and have their points reset to zero</li>
	<li>If there is a point tie, then a /roll will determine the winner. Only the winner spends their points</li>
	<li>There is a maximum cap on points: number_of_bosses * 3 + 1</li>
	<li>Points are reset when we move to a new raid.</li>
	<li>This system is only in effect for the current raid. Any old raid runs will be simple MS>OS rolls</li>
	</ul>
	`
	t := template.Must(template.New("").Parse(tmpl))
	if err := t.Execute(w, rs); err != nil {
		util.CheckErr(err)
	}
}
