package main

import (
	"database/sql"
	_ "github.com/Go-SQL-Driver/MySQL"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	db, e := sql.Open("mysql", DB_USER+":"+DB_PASS+"@/"+DB_NAME+"?charset=utf8")
	if e != nil {
		panic(e)
	}
	defer db.Close()

	r := mux.NewRouter()

	runApiServer(db, r.PathPrefix("/api").Subrouter())

	http.Handle("/", r)
	http.ListenAndServe(":1802", nil)
}
