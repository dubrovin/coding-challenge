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
	storage := NewStorage("test", time.Second*60)
	require.NotNil(t, storage)
	require.Equal(t, "test", storage.filePath)
}

func TestStoragePersist(t *testing.T) {
	testFile := "test"
	nodesNum := 10
	storage := NewStorage(testFile, time.Second*1)
	require.NotNil(t, storage)
	require.Equal(t, testFile, storage.filePath)
	for i := 0; i < nodesNum; i++ {
		storage.Add(time.Now())
		time.Sleep(time.Millisecond * 88)
	}
	storage.Persist()
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
	os.Remove(testFile)
}
