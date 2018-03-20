package pool

import (
	"net/url"
	"sync"
)

// Node is a type alias for url.URL.
type Node = url.URL

// Pool represents pool data structure.
type Pool struct {
	In      chan map[Node]bool
	Out     chan Node
	status  map[Node]bool
	nodes   []Node
	once    sync.Once
	current int
}

// New returns a pool with requested selection policy.
func New(nodes []Node) *Pool {
	p := &Pool{
		In:      make(chan map[Node]bool, 1),
		Out:     make(chan Node),
		status:  map[Node]bool{},
		nodes:   nodes,
		current: 0,
	}
	for _, node := range nodes {
		p.status[node] = true
	}
	return p
}

// Run takes care of nodes status updates and nodes delivery.
func (p *Pool) Run() {
	p.once.Do(func() {
		go func() {
			for {
				select {
				case newStatus := <-p.In:
					p.nodes = make([]Node, 0, len(newStatus))
					for node, status := range newStatus {
						p.status[node] = status
						if status {
							p.nodes = append(p.nodes, node)
						}
					}
				default:
					if p.current >= len(p.nodes) {
						p.current = 0
					}
					p.Out <- p.nodes[p.current]
					p.current++
				}
			}
		}()
	})
}

// Get returns next available node.
// func (p *Pool) Get() <-chan Node {
// 	return p.out
// }

// Update updates status of the existing nodes.
// func (p *Pool) Update() chan<- map[Node]bool {
// 	return p.in
// }

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
