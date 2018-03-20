package main

import (
	"fmt"
	"net/url"
	"sync"
	"time"
)

// New returns a pool with requested selection policy.
func new(nodes []url.URL) *pool {
	p := &pool{
		in:      make(chan map[url.URL]bool),
		out:     make(chan url.URL),
		status:  map[url.URL]bool{},
		nodes:   nodes,
		done:    make(chan struct{}),
		current: 0,
	}
	for _, node := range nodes {
		p.status[node] = true
	}
	return p
}

type pool struct {
	in      chan map[url.URL]bool
	out     chan url.URL
	status  map[url.URL]bool
	nodes   []url.URL
	done    chan struct{}
	once    sync.Once
	current int
}

func (p *pool) run() {
	p.once.Do(func() {
		go func() {
			for {
				select {
				case newStatus := <-p.in:
					p.nodes = make([]url.URL, 0, len(newStatus))
					i := 0
					for node, status := range newStatus {
						p.status[node] = status
						if status {
							p.nodes = append(p.nodes, node)
							i++
						}
					}
				case <-p.done:
					return
				default:
					if p.current == len(p.nodes) {
						p.current = 0
					}
					p.out <- p.nodes[p.current]
					p.current++
				}
			}
		}()
	})
}

func main() {
	pool := new(nodes())
	pool.run()
	for i := 0; i < 10; i++ {
		node := <-pool.out
		fmt.Printf("%d: %s\n", i, node.String())
	}
	pool.in <- map[url.URL]bool{
		url.URL{Scheme: "http", Host: "127.0.0.1"}: true,
		url.URL{Scheme: "http", Host: "127.0.0.2"}: true,
		url.URL{Scheme: "http", Host: "127.0.0.3"}: false,
	}
	for i := 0; i < 10; i++ {
		node := <-pool.out
		fmt.Printf("%d: %s\n", i, node.String())
	}
	time.Sleep(5 * time.Second)
}

func nodes() []url.URL {
	return []url.URL{
		url.URL{Scheme: "http", Host: "127.0.0.1"},
		url.URL{Scheme: "http", Host: "127.0.0.2"},
		url.URL{Scheme: "http", Host: "127.0.0.3"},
	}
}

// // GetNode returns node from the pool using round-robin algorithm.
// func (p *RoundRobinPool) GetNode() (url.URL, error) {
// 	p.Lock()
// 	defer p.Unlock()

// 	if len(p.nodes) == 0 {
// 		return url.URL{}, fmt.Errorf("no nodes in the pool")
// 	}

// 	if p.current == len(p.nodes) {
// 		p.current = 0
// 	}

// 	result := p.nodes[p.current]
// 	p.current++
// 	return result, nil
// }

// // RunStatusUpdate starts goroutine to update nodes status
// func (p *basePool) RunStatusUpdate(interval time.Duration) {
// 	p.once = sync.Once{}
// 	p.once.Do(func() {
// 		p.logger.Debug("status update", "interval", interval.String())
// 		go func() {
// 			t := time.Tick(interval)

// 			for now := range t {
// 				p.logger.Debug("status update", "time", now.String())

// 				node, err := p.GetNode()
// 				if err != nil {
// 					p.logger.Warn("status update", "GetNode error", err)
// 				}

// 				ns, err := f(node)
// 				if err != nil {
// 					p.logger.Warn("status update", "get status update error", err)
// 				}

// 				err = p.UpdateStatus(ns)
// 				if err != nil {
// 					p.logger.Warn("status update", "error updating status", err)
// 				}
// 			}
// 		}()
// 	})
// }

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
