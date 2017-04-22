package main

import (
	"flag"
	"fmt"
	"github.com/dubrovin/coding-challenge/server"
	"time"
)

var (
	addr      = flag.String("addr", ":8080", "http service address")
	filepath  = flag.String("file", "storage.txt", "storage file path")
	countTime = flag.String("count", "60s", "time for counting")
)

func main() {
	flag.Parse()
	duration, err := time.ParseDuration(*countTime)
	if err != nil {
		fmt.Print("Can't parse duration, set to default 60s")
		duration = time.Second * 60

	}
	newServer := server.NewServer(*addr, *filepath, duration)
	newServer.Run()
}
