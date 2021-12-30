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
		err = rows.Scan(&r.Id, &r.Name, &r.Points)
		util.CheckErr(err)
		rs = append(rs, r)
	}

	const tmpl = `
	<style>
	body {
		background-color: #282828;
		color: #ebdbb2;
	}
	h1 { color: #83a598; }
	table, th, td { 
    	border: 1px solid #504945;
		text-align: left; 
	}
	table { width: 50%; }	
	th { color: #d65d03; }
	</style>
	<h1>KCC Loot Tracker</h1>
	<table>
	<tr>
		<th>Name</th>
		<th>Points</th>
	</tr>
		{{range .}}
			<tr><td>{{.Name}}</td><td>{{.Points}}</td></tr>
		{{end}}
	</table>`
	t := template.Must(template.New("").Parse(tmpl))
	if err := t.Execute(w, rs); err != nil {
		util.CheckErr(err)
	}

	/*
		for _, r := range rs {
			fmt.Fprintf(w, "%s: %d\n", r.Name, r.Points)
		}
	*/
}
