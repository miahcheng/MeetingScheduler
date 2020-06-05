package main

import (
	"database/sql"
	"fmt"
	"info441-finalproj/servers/gateway/handlers"
	"info441-finalproj/servers/gateway/models/users"
	"info441-finalproj/servers/gateway/sessions"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	addr := os.Getenv("ADDR")
	cert := os.Getenv("TLSCERT")
	key := os.Getenv("TLSKEY")
	sess := os.Getenv("SESSIONKEY")
	redisAddr := os.Getenv("REDISADDR")
	meetingAddr := os.Getenv("MEETINGADDR")
	dsn := os.Getenv("dsn")

	if len(addr) == 0 {
		addr = ":443"
	}

	if len(cert) == 0 || len(key) == 0 {
		fmt.Fprintln(os.Stderr, "Either the key or certificate was not found")
		os.Exit(1)
	}

	rclient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	dur, err2 := time.ParseDuration("24h")
	if err2 != nil {
		log.Fatal(err)
	}

	handler := handlers.HandlerContext{
		SigningKey:   sess,
		SessionStore: sessions.NewRedisStore(rclient, dur),
		UserStore:    users.GetNewStore(db),
	}

	mux := http.NewServeMux()

	meetingDirector := func(r *http.Request) {
		addresses := strings.Split(meetingAddr, ", ")
		serv := addresses[0]
		if len(addresses) > 1 {
			rand.Seed(time.Now().UnixNano())
			serv = addresses[rand.Intn(len(addresses))]
		}
		r.Host = serv
		r.URL.Host = serv
		r.URL.Scheme = "http"
	}

	meetingProxy := &httputil.ReverseProxy{Director: meetingDirector}

	mux.HandleFunc("/users", handler.UsersHandler)
	mux.HandleFunc("/sessions", handler.SessionsHandler)
	mux.Handle("/user/", meetingProxy)
	mux.Handle("/meeting/{id}", meetingProxy)
	mux.Handle("/meeting", meetingProxy)

	newMux := handlers.NewCorsMW(mux)

	log.Fatal(http.ListenAndServeTLS(addr, cert, key, newMux))
}
