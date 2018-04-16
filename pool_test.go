// +build race

package lilium_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/acim/lilium"
)

func Test(t *testing.T) {
	p := lilium.Pool{Nodes: nodes()}
	p.Run()

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			node := p.Node()
			fmt.Printf("%d: %s\n", i, node.String())

		}(i)
	}
	wg.Wait()
}

func nodes() []lilium.Node {
	return []lilium.Node{
		lilium.Node{Scheme: "http", Host: "127.0.0.1"},
		lilium.Node{Scheme: "http", Host: "127.0.0.2"},
		lilium.Node{Scheme: "http", Host: "127.0.0.3"},
	}
}
