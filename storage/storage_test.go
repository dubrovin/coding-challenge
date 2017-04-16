package storage

import (
	"testing"
	"github.com/stretchr/testify/require"
	"time"
	"os"
	"bufio"
	"io"
)

func TestNewStorage(t *testing.T) {
	storage := NewStorage("test")
	require.NotNil(t, storage)
	require.Equal(t, "test", storage.filePath)
}

func TestStoragePersist(t *testing.T) {
	testFile := "test"
	nodesNum := 10
	storage := NewStorage(testFile)
	require.NotNil(t, storage)
	require.Equal(t, testFile, storage.filePath)
	for i:=0; i < nodesNum; i++ {
		storage.Add(time.Now())
		time.Sleep(time.Millisecond * 88)
	}
	storage.Persist()
	f, err := os.Open(testFile)
	require.Nil(t, err)
	require.NotNil(t, f)
	r := bufio.NewReader(f)
	for i:=0; i < nodesNum; i++ {
		_, _, err  := r.ReadLine()
		require.Nil(t, err)
	}

	_, _, err  = r.ReadLine()
	require.Equal(t, io.EOF, err)
	os.Remove(testFile)
}