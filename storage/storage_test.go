package storage

import (
	"bufio"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
	"time"
)

func TestNewStorage(t *testing.T) {
	storage := NewStorage("test1", time.Second*60)
	require.NotNil(t, storage)
	require.Equal(t, "test1", storage.filePath)
}

func TestStoragePersist(t *testing.T) {
	testFile := "test1"
	nodesNum := 10
	storage := NewStorage(testFile, time.Second*2)
	go storage.Worker()
	require.NotNil(t, storage)
	require.Equal(t, testFile, storage.filePath)
	for i := 0; i < nodesNum; i++ {
		storage.Inc(NewNode(time.Now()))
	}
	storage.persistCh <- struct {}{}
	time.Sleep(time.Second)
	f, err := os.Open(testFile)
	require.Nil(t, err)
	require.NotNil(t, f)
	r := bufio.NewReader(f)
	for i := 0; i < nodesNum; i++ {
		_, _, err := r.ReadLine()
		require.Nil(t, err)
	}

	_, _, err = r.ReadLine()
	require.Equal(t, io.EOF, err)
	require.Equal(t, nodesNum, storage.GetCount())
	time.Sleep(time.Second * 2)
	require.Equal(t, 0, storage.GetCount())
	time.Sleep(time.Second)
	os.Remove(testFile)

	//storage.cleanCh <- struct {}{}
	//time.Sleep(time.Second)
	//require.Empty(t, storage.nodes)
	//for i := 0; i < nodesNum; i++ {
	//	storage.Inc(NewNode(time.Now()))
	//	time.Sleep(time.Millisecond * 88)
	//}
	//require.Len(t, storage.nodes, nodesNum)
}

func TestStorageLoad(t *testing.T) {
	testFile := "test2"
	nodesNum := 10
	storage := NewStorage(testFile, time.Second*60)
	go storage.Worker()
	require.NotNil(t, storage)
	require.Equal(t, testFile, storage.filePath)
	f, _ := os.Create(testFile)
	for i := 0; i < nodesNum; i++ {
		f.WriteString(time.Now().Format(time.RFC3339Nano)+"\n")
	}
	f.Close()
	storage.Load()
	require.Equal(t,nodesNum, storage.GetCount())
	os.Remove(testFile)
}
