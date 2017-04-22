package server

import (
	"encoding/json"
	"github.com/dubrovin/coding-challenge/storage"
	"log"
	"net/http"
	"time"
	"fmt"
)

type Server struct {
	ListenAddr string
	Storage    *storage.Storage
}

func NewServer(listenAddr, filePath string, countTime time.Duration) *Server {
	return &Server{
		ListenAddr: listenAddr,
		Storage:    storage.NewStorage(filePath, countTime),
	}
}

func (s *Server) Run() error {
	http.HandleFunc("/counter", serverHandler(s))
	go s.Storage.Worker()
	err := s.Storage.Load()
	if err != nil {
		fmt.Println(err)
	}
	go s.Storage.Persister("1s")
	go s.Storage.Cleaner("60s")
	err = http.ListenAndServe(s.ListenAddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		return err
	}
	return nil
}

func serverHandler(s *Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		s.Storage.Inc(storage.NewNode(time.Now()))
		counter := map[string]int{"Counter": s.Storage.GetCount()}
		json.NewEncoder(w).Encode(counter)
	}
}
