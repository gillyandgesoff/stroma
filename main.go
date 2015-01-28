package main

import (
	"flag"
	"fmt"
	r "github.com/dancannon/gorethink"
	"log"
	"net/http"
)

var session *r.Session

type User struct {
	Id		string	`gorethink:"id,omitempty"`
	Name	string	`gorethink:"name"`
	UName	string	`gorethink:"username"`
	Pin		string	`gorethink:"pin"`
	Form	string	`gorethink:"form"`
}

func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Hello Chappo!")
}

func listUsers(res http.ResponseWriter, req *http.Request) {
	rows, err := r.Table("user").Run(session)
	if err != nil {
		log.Fatalln(err.Error())
	}
	var users []User
	err = rows.All(&users)
	if err != nil {
		log.Fatalln(err.Error())
	}
	var resstring string
	for _, user := range users {
		resstring = resstring + user.Name + " " + user.Id + " " + user.Form
	}
	fmt.Fprintf(res, resstring)
}

func main() {
	dbAddress := flag.String("dbaddr", "localhost:28015", "The address and port of the RethinkDB cluster in host:port format.")
	dbName := flag.String("dbname", "stroma", "The name of the db in RethinkDB")
	maxConn := flag.Int("maxconn", 1, "The maximum number of active connections in the RethinkDB connection pool.")
	flag.Parse()
	var err error
	session, err = r.Connect(r.ConnectOpts{
		Address:   *dbAddress,
		Database:  *dbName,
		MaxActive: *maxConn,
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/users", listUsers)
	http.ListenAndServe(":3000", nil)
}
