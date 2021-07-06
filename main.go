package main

import (
	"github.com/onlineGo/route"
	"github.com/onlineGo/conf"

	"net/http"
	"time"
	"fmt"
)


func main() {
	// Echo instance
	conf.Init()

	e := route.Register()

	addr := fmt.Sprintf("%s:%s", "localhost", "8000")
	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	println(addr)
	// Start server
	e.Logger.Fatal(e.StartServer(server))
}
