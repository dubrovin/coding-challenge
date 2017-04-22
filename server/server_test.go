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
	testFile := "testFile1"
	newServer := NewServer(addr, testFile , time.Second * 60)
	require.Equal(t, addr, newServer.ListenAddr)
}

func TestServerRun(t *testing.T) {
	addr := ":8080"
	testFile := "servertest2"
	newServer := NewServer(addr, testFile, time.Second * 60)
	go newServer.Run()

	var cnt Counter
	time.Sleep(time.Second)
	for i := 0; i < 100; i++ {
		resp, err := http.Get("http://127.0.0.1:8080/counter")
		require.Nil(t, err)
		require.NotNil(t, resp)

		body, err := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &cnt)
		require.Nil(t, err)
		require.Equal(t, i+1, cnt.Counter)
	}
	newServer.Storage.Stop()
}
