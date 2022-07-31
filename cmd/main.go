package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/appak21/forum/pkg/config"
	"github.com/appak21/forum/pkg/controllers"
	"github.com/appak21/forum/pkg/routes"
)

func main() {
	addr := flag.String("addr", ":9010", "http network address")
	flag.Parse()
	logInf, logErr := config.Logger()
	srv := &http.Server{
		Addr:         *addr,
		Handler:      routes.RegisterForumRoutes(),
		ErrorLog:     logErr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go controllers.DeleteExpiredSessions()
	logInf.Printf("The server is running on port %s", *addr)
	logErr.Fatal(srv.ListenAndServe())
}
