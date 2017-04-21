package main

import (
	"time"
	"github.com/dubrovin/coding-challenge/server"
)




func main() {
	addr := ":8080"
	newServer := server.NewServer(addr, "test", time.Second * 60)
	newServer.Run()
}
