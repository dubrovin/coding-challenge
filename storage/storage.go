package storage

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"io"
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
	nodes     []*Node
	countTime time.Duration
	countCh   chan *Node
	persistCh chan struct{}
	cleanCh   chan struct{}
	lenCh chan int
	stop chan struct{}
}

func NewStorage(filePath string, countTime time.Duration) *Storage {
	return &Storage{
		filePath:  filePath,
		nodes:     make([]*Node, 0),
		countTime: countTime,
		countCh:   make(chan *Node),
		persistCh: make(chan struct{}),
		cleanCh:   make(chan struct{}),
		lenCh: make(chan int),
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
		fmt.Println("NODE ", node)
		if !node.persisted{
			_, err = w.WriteString(fmt.Sprint(node.timestamp.String(), "\n"))
			if err != nil {
				fmt.Println(err)
			}
			node.persisted = true
		}


	}

	err = w.Flush()
	if err != nil {
		return err
	}
	fmt.Println("Persisted ")
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

func (storage *Storage) Add(node *Node) {
	storage.nodes = append(storage.nodes, node)
}

func (storage *Storage) GetCount() int {
	storage.lenCh <- 0
	return <-storage.lenCh
}

func (storage *Storage) filter() []*Node {

	var p []*Node

	tEnd := time.Now()
	tBegin := tEnd.Add(-storage.countTime)
	for _, node := range storage.nodes {
		if node.timestamp.After(tBegin) && node.timestamp.Before(tEnd) {
			p = append(p, node)
		}
	}
	return p
}

func (storage *Storage) Persister(duration string) error {
	dur, err := time.ParseDuration(duration)
	if err != nil {
		return err
	}

	tiker := time.NewTicker(dur)
	defer tiker.Stop()

	for t := range tiker.C {
		fmt.Println("Persisted at ", t.String())
		storage.persistCh <- struct{}{}
	}
	return nil
}

func (storage *Storage) clean() {

	var p []*Node

	tEnd := time.Now()
	tBegin := tEnd.Add(-storage.countTime)
	for _, node := range storage.nodes {
		if node.timestamp.After(tBegin) && node.timestamp.Before(tEnd) && !node.persisted {
			p = append(p, node)
		}
	}

	storage.nodes = p
}

func (storage *Storage) Cleaner(duration string) error {
	dur, err := time.ParseDuration(duration)
	if err != nil {
		return err
	}

	tiker := time.NewTicker(dur)
	defer tiker.Stop()
	for t := range tiker.C {
		fmt.Println("Cleaned at ", t.String())
		storage.cleanCh <- struct{}{}
	}

	return nil
}

func (storage *Storage) Inc(node *Node) {
	storage.countCh <- node
}

func (storage *Storage) Worker() {

	for {
		select {
		case node := <-storage.countCh:
			storage.Add(node)
		case <-storage.persistCh:
			err := storage.persist()
			fmt.Println(err)
		case <-storage.cleanCh:
			storage.clean()
		case <-storage.lenCh:
			storage.lenCh<-len(storage.filter())
		case <-storage.stop:
			break
		}

	}
}


func (storage *Storage) Load(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for err != io.EOF {
		line, _, err := r.ReadLine()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(line))
	}

	return nil

}