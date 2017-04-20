package storage

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type Node struct {
	timestamp time.Time
	persisted bool
}

func NewNode(timestamp time.Time) *Node {
	return &Node{timestamp: timestamp, persisted: false}
}

type Storage struct {
	filePath  string
	nodes     []Node
	countTime time.Duration
}

func NewStorage(filePath string, countTime time.Duration) *Storage {
	return &Storage{
		filePath:  filePath,
		nodes:     make([]Node, 0),
		countTime: countTime,
	}
}

func (storage *Storage) persist() error {
	f, err := storage.open()

	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, node := range storage.nodes {
		_, err = w.WriteString(fmt.Sprint(node.timestamp.String(), "\n"))
		if err != nil {
			fmt.Println(err)
		}
		node.persisted = true
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

func (storage *Storage) Add(node Node) {
	storage.nodes = append(storage.nodes, node)
}

func (storage *Storage) GetCount() int {
	return len(storage.filter())
}

func (storage *Storage) filter() []Node {

	var p []Node

	tEnd := time.Now()
	tBegin := tEnd.Add(-storage.countTime)

	for _, node := range storage.nodes {
		if node.timestamp.After(tBegin) && node.timestamp.Before(tEnd) {
			p = append(p, node)
		}
	}
	return p
}

func (storage *Storage) Persister(duration string) error{
	dur, err := time.ParseDuration(duration)
	if err != nil {
		return err
	}

	tiker :=time.NewTicker(dur)
	defer tiker.Stop()

	for t := range tiker.C {
		fmt.Println("Persisted at ", t.String())
		err = storage.persist()
		if err != nil {
			fmt.Println("Persister error ", err)
		}
	}
	return nil
}

func (storage *Storage) clean() {
	storage.nodes = storage.filter()
}

func (storage *Storage) Cleaner(duration string) error {
	dur, err := time.ParseDuration(duration)
	if err != nil {
		return err
	}

	tiker :=time.NewTicker(dur)
	defer tiker.Stop()
	for t := range tiker.C {
		fmt.Println("Cleaned at ", t.String())
		storage.clean()
	}


	return nil
}
