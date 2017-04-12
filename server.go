package codingchallenge

import (
	"sync"
	"net/http"
	"log"
	"encoding/json"
)

type Server struct {
	ListenAddr string
	Counter    int
	wg         *sync.WaitGroup
	mu         *sync.Mutex
}

func NewServer(listenAddr string) *Server {
	return &Server{
		ListenAddr: listenAddr,
		wg: &sync.WaitGroup{},
		mu: &sync.Mutex{},
	}
}

func (s *Server) Run() error {
	http.HandleFunc("/counter", serverHandler(s))
	err := http.ListenAndServe(s.ListenAddr, nil)
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
		s.mu.Lock()
		s.Counter += 1
		counter := map[string]int{"Counter": s.Counter}
		s.mu.Unlock()
		json.NewEncoder(w).Encode(counter)

	}
}

//func(s *Server) Persist(filepath string) error{
//	w := bufio.NewWriter(filepath)
//    	n4, err := w.WriteString("buffered\n")
//    	fmt.Printf("wrote %d bytes\n", n4)
//    	w.Flush()
//}