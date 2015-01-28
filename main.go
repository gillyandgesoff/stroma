package main

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
)

func namesHandler(statement *sql.Stmt) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var names string
		rows, err := statement.Query()
		defer rows.Close()
		if err != nil {
			panic(err.Error())
		}
		var (
			name string
			pin int
		)
		for rows.Next() {
			err := rows.Scan(&name, &pin)
			if err != nil {
				panic(err.Error())
			}
			names += (name + " " + strconv.Itoa(pin) + "\n")
		}
		fmt.Fprintf(w, names)
	}
}
func serverH(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}
func main() {
	db, err := sql.Open("mysql", "go_db_server:stroma@/stroma?parseTime=true")
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	statement, err := db.Prepare("SELECT * FROM test")
	if err != nil {
		panic(err.Error())
	}
	defer statement.Close()
	hf := namesHandler(statement)
	http.HandleFunc("/", hf)
	http.ListenAndServe(":8080", nil)
}
