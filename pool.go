package lilium

import (
	"errors"
	"net/url"
	"sync"
)

// Node is a type alias for url.URL.
type Node = url.URL

// Pool represents pool data structure.
type Pool struct {
	Nodes   []Node
	req     chan chan Node
	current int
	once    sync.Once
}

// Run starts a goroutine to accept requests to get next node.
func (p *Pool) Run() error {
	if len(p.Nodes) == 0 {
		return errors.New("No nodes defined")
	}

	p.once.Do(func() {
		p.req = make(chan chan Node)
		go func() {
			for req := range p.req {

				if p.current >= len(p.Nodes) {
					p.current = 0
				}

				req <- p.Nodes[p.current]
				p.current++
			}
		}()
	})

	return nil
}

// Node returns next available node.
func (p *Pool) Node() Node {
	r := make(chan Node)
	p.req <- r
	defer close(r)
	return <-r
}
