package storage

import (
	"testing"
	"github.com/stretchr/testify/require"
	"time"
	"os"
	"bufio"
)

func TestNewStorage(t *testing.T) {
	storage := NewStorage("test")
	require.NotNil(t, storage)
	require.Equal(t, "test", storage.filePath)
}

func TestStoragePersist(t *testing.T) {
	storage := NewStorage("test")
	require.NotNil(t, storage)
	require.Equal(t, "test", storage.filePath)
	for i:=0; i < 10; i++ {
		storage.Add(time.Now())
		time.Sleep(time.Millisecond * 88)
	}
	storage.Persist()
	f, err := os.Open("test")
	require.Nil(t, err)
	require.NotNil(t, f)
	r := bufio.NewReader(f)
	for i:=0; i < 10; i++ {
		_, _, err  := r.ReadLine()
		require.Nil(t, err)
	}

	//_, _, err  = r.ReadLine()
	//require.Equal(t, io.EOF, err)
	//os.Remove("test")
}