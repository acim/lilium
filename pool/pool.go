package pool

import (
	"fmt"
	"net/url"
	"sync"
	"time"
)

// Node is a type alias for url.URL.
type Node = url.URL

// Pool represents pool data structure.
type Pool struct {
	Req     chan chan Node
	status  map[Node]bool
	nodes   []Node
	once    sync.Once
	done    chan struct{}
	current int
}

// New returns a pool with requested selection policy.
func New(nodes []Node) *Pool {
	p := &Pool{
		Req:     make(chan chan Node),
		status:  map[Node]bool{},
		nodes:   nodes,
		done:    make(chan struct{}),
		current: 0,
	}
	for _, node := range nodes {
		p.status[node] = true
	}
	return p
}

// Start takes care of nodes status updates and nodes delivery.
func (p *Pool) Start(interval time.Duration) {
	p.once.Do(func() {
		tick := time.NewTicker(interval)
		go func() {
			for {
				select {
				case res := <-p.Req:
					if p.current >= len(p.nodes) {
						p.current = 0
					}
					res <- p.nodes[p.current]
					p.current++
				case t := <-tick.C:
					fmt.Printf("UPDATE STATUS %s\n", t.String())
				case <-p.done:
					return
				}
			}
		}()
	})
}

// Get returns next available node.
func (p *Pool) Get() Node {
	res := make(chan Node)
	p.Req <- res
	return <-res
}

// Stop stops pool functionality.
func (p *Pool) Stop() {
	p.done <- struct{}{}
}

// func (p *basePool) UpdateStatus(ns *map[url.URL]bool) error {
// 	for n, newStatus := range *ns {
// 		p.RLock()
// 		oldStatus, ok := p.status[n]
// 		p.RUnlock()
// 		if !ok {
// 			return fmt.Errorf("cannot update status for unknown node %s", n.String())
// 		}

// 		if newStatus != oldStatus {
// 			switch newStatus {
// 			case true:
// 				p.addNode(n)
// 			case false:
// 				p.removeNode(n)
// 			}
// 		}
// 	}
// 	return nil
// }

// func (p *basePool) addNode(n url.URL) {
// 	p.Lock()
// 	defer p.Unlock()
// 	p.nodes = append(p.nodes, n)
// 	p.status[n] = true
// }

// func (p *basePool) removeNode(n url.URL) {
// 	p.RLock()
// 	for i, node := range p.nodes {
// 		if n == node {
// 			p.RUnlock()
// 			p.Lock()
// 			defer p.Unlock()
// 			p.nodes = append(p.nodes[:i], p.nodes[i+1:]...)
// 			p.status[n] = false
// 			return
// 		}
// 	}
// }
