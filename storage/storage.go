package storage

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)
// Node - structure for each request
type Node struct {
	timestamp time.Time
	persisted bool
}

// NewNode - creates new node, by default persisted = false
func NewNode(timestamp time.Time) *Node {
	return &Node{timestamp: timestamp, persisted: false}
}

// Storage -
type Storage struct {
	// filePath - path where we persists requests data
	filePath  string

	// nodes - in memory cache of requests
	nodes     []*Node

	// countTime - time for counting total number of requests
	countTime time.Duration

	// countCh - channel for counting
	countCh   chan *Node

	// persistCh - channel for persisting
	persistCh chan struct{}

	// cleanCh - channel for clean in memory cache of nodes
	cleanCh   chan struct{}

	// lenCh - channel for receiving current len of nodes
	lenCh     chan int

	// stop - channel for stopping worker
	stop      chan struct{}
}

func NewStorage(filePath string, countTime time.Duration) *Storage {
	return &Storage{
		filePath:  filePath,
		nodes:     make([]*Node, 0),
		countTime: countTime,
		countCh:   make(chan *Node),
		persistCh: make(chan struct{}),
		cleanCh:   make(chan struct{}),
		lenCh:     make(chan int),
		stop:      make(chan struct{}),
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
		if !node.persisted {
			_, err = w.WriteString(fmt.Sprint(node.timestamp.Format(time.RFC3339Nano), "\n"))
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

func (storage *Storage) add(node *Node) {
	storage.nodes = append(storage.nodes, node)
}

// GetCount - returns current count of requests
func (storage *Storage) GetCount() int {
	storage.lenCh <- 0
	return <-storage.lenCh
}

// filter - filter all nodes in cache by the given storage.countTime
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

// Persister - runs persist by tiker with given duration
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

// clean - clean cache of nodes
func (storage *Storage) clean() {
	var p []*Node
	tEnd := time.Now()
	tBegin := tEnd.Add(-storage.countTime)
	// remove all nodes which time greater than storage.countTime
	for _, node := range storage.nodes {
		if node.timestamp.After(tBegin) && node.timestamp.Before(tEnd) && !node.persisted {
			p = append(p, node)
		}
	}

	storage.nodes = p
}

// Cleaner - runs clean by tiker with given duration
func (storage *Storage) Cleaner() error {
	tiker := time.NewTicker(storage.countTime)
	defer tiker.Stop()
	for t := range tiker.C {
		fmt.Println("Cleaned at ", t.String())
		storage.cleanCh <- struct{}{}
	}

	return nil
}

// Inc - increase cache by node
func (storage *Storage) Inc(node *Node) {
	storage.countCh <- node
}

// Worker - runs worker which listen channel
func (storage *Storage) Worker() {

	for {
		select {
		case node := <-storage.countCh:
			storage.add(node)
		case <-storage.persistCh:
			err := storage.persist()
			if err != nil {
				fmt.Println("persist error ", err)
			}
		case <-storage.cleanCh:
			storage.clean()
		case <-storage.lenCh:
			storage.lenCh <- len(storage.filter())
		case <-storage.stop:
			break
		}

	}
}

// Load - loads nodes to cache from given storage.filePath
func (storage *Storage) Load() error {
	f, err := os.Open(storage.filePath)
	if err != nil {
		fmt.Println("open error", err)
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		} else {
			t, err := time.Parse(time.RFC3339Nano, string(line))
			if err != nil {
				fmt.Println("PARSE ERROR", err)
				return err
			}
			persistedNode := NewNode(t)
			persistedNode.persisted = true
			storage.Inc(persistedNode)
		}

	}

	return nil
}

// Stop - stops worker
func (storage *Storage) Stop() {
	storage.stop <- struct{}{}
}
