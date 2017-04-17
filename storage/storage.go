package storage

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type Storage struct {
	filePath  string
	nodes     []time.Time
	countTime time.Duration
}

func NewStorage(filePath string, countTime time.Duration) *Storage {
	return &Storage{
		filePath:  filePath,
		nodes:     make([]time.Time, 0),
		countTime: countTime,
	}
}

func (storage *Storage) Persist() error {
	f, err := storage.open()

	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, node := range storage.nodes {
		_, err = w.WriteString(fmt.Sprint(node.String(), "\n"))
		if err != nil {
			fmt.Println(err)
		}
	}

	err = w.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (storage *Storage) create() (*os.File, error) {
	f, err := os.Create(storage.filePath)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (storage *Storage) open() (*os.File, error) {
	f, err := os.OpenFile(storage.filePath, os.O_RDWR|os.O_APPEND, 0660)

	if os.IsNotExist(err) {
		fmt.Println(err, "create new file with name ", storage.filePath)
		return storage.create()
	}
	return f, err
}

func (storage *Storage) Add(node time.Time) {
	storage.nodes = append(storage.nodes, node)
}

func (storage *Storage) GetCount() int {
	return len(storage.filter())
}

func (storage *Storage) filter() []time.Time {

	var p []time.Time

	tEnd := time.Now()
	tBegin := tEnd.Add(-storage.countTime)

	for _, node := range storage.nodes {
		if node.After(tBegin) && node.Before(tEnd) {
			p = append(p, node)
		}
	}
	return p
}
