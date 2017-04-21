package server

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type Counter struct {
	Counter int `json:"Counter"`
}

func TestNewServer(t *testing.T) {
	addr := ":8081"
	newServer := NewServer(addr, "test1", time.Second * 60)
	require.Equal(t, addr, newServer.ListenAddr)
}

func TestServerRun(t *testing.T) {
	addr := ":8080"
	newServer := NewServer(addr, "test", time.Second * 60)
	go newServer.Run()

	var cnt Counter

	for i := 0; i < 100; i++ {
		resp, err := http.Get("http://127.0.0.1:8080/counter")
		require.Nil(t, err)
		require.NotNil(t, resp)

		body, err := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &cnt)
		require.Nil(t, err)
		require.Equal(t, i+1, cnt.Counter)
	}

}
