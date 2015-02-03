package main

import (
	"flag"
	"fmt"
	r "github.com/dancannon/gorethink"
	"log"
	"net/http"
	"github.com/julienschmidt/httprouter"
)
// global session object for RethinkDB connection pool
var session *r.Session

// struct representing user document in db
type User struct {
	Id		string	`gorethink:"id,omitempty"`
	Name	string	`gorethink:"name"`
	UName	string	`gorethink:"username"`
	Pin		string	`gorethink:"pin"`
	Form	string	`gorethink:"form"`
}

// struct representing the event in the db
type StatusChange struct {
	Id		string	`gorethink:"id,omitempty"`
	Current	string	`gorethink:"currentstatus"`
	Time	string	`gorethink:"lastupdated"`
}

// test route
func hello(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	fmt.Fprintf(res, authUser("r", "2948"))
}

func checkUser(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// check user is authenticated, returns id
	id := authUser(ps.ByName("username"), ps.ByName("pin"))
	if id == "failed" {
		fmt.Fprintf(res, "authentication failed")
	} else {
		// finds the status of the authenticated user
		rows, err := r.Table("status").Get(id).Run(session)
		if err != nil {
			// db error, not user's fault
			log.Fatalln(err.Error() + ps.ByName("username"))
			fmt.Fprintf(res, "failed, contact support")
		} else {
			// return current status
			var status StatusChange
			err = rows.One(&status)
			fmt.Fprintf(res, status.Current)
		}
	}
}

//route that lists the users in the db
func listUsers(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// db query, returns cursor
	rows, err := r.Table("user").Run(session)
	if err != nil {
		log.Fatalln(err.Error())
	}
	var users []User
	// scans cursor to struct
	err = rows.All(&users)
	if err != nil {
		log.Fatalln(err.Error())
	}
	// compiles a string of all the users
	var resstring string
	for _, user := range users {
		resstring = resstring + user.Name + " " + user.Id + " " + user.Form + "\n"
	}
	fmt.Fprintf(res, resstring)
}

// authenticates user
func authUser(username string, pin string) string {
	rows, err := r.Table("user").Filter(r.Row.Field("username").Eq(username)).Run(session)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if rows.IsNil() {
		return "failed"
	}
	var user User
	err = rows.One(&user)
	if user.Pin == pin {
		return user.Id
	} else {
		return "failed"
	}
}

func main() {
	// parse commandline flags
	// dbAddress := flag.String("dbaddr", "localhost:28015", "The address and port of the RethinkDB cluster in host:port format.")
	// dbName := flag.String("dbname", "stroma", "The name of the db in RethinkDB")
	maxConn := flag.Int("maxconn", 1, "The maximum number of active connections in the RethinkDB connection pool.")
	flag.Parse()
	var err error
	// connect to db and initialise session with options
	session, err = r.Connect(r.ConnectOpts{
		Address:   "db:28015",
		Database:  "stroma",
		MaxActive: *maxConn,
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	router := httprouter.New()
	router.GET("/hello", hello)
	router.GET("/users", listUsers)
	router.GET("/user/:username/:pin/status", checkUser)
	//TODO: router.POST("/user/:username/:pin/status", setUser)
	// need to know how to update documents, set time
	// need to restructure status table to allow history
	http.ListenAndServe(":3000", router)
}
