package main

import (
	"database/sql"
	"info441-finalproj/servers/meeting/meetings"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	mux := mux.NewRouter()
	addr := os.Getenv("ADDR")
	dsn := os.Getenv("DSN")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	server := meetings.Context{CalendarStore: db}

	mux.HandleFunc("/user/", server.SpecificUserHandler)
	mux.HandleFunc("/meeting/{id}", server.SpecificMeetingHandler)
	mux.HandleFunc("/meeting", server.MeetingsHandler)

	http.ListenAndServe(addr, mux)
}
