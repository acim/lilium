// +build race

package lilium_test

import (
	"math"
	"sync"
	"testing"

	"github.com/acim/lilium"
)

func Test(t *testing.T) {
	ns := nodes()
	p := lilium.Pool{Nodes: ns}
	if err := p.Run(); err != nil {
		t.Error(err)
	}

	distribution := map[string]int{
		"127.0.0.1": 0,
		"127.0.0.2": 0,
		"127.0.0.3": 0,
	}

	mux := sync.RWMutex{}
	wg := sync.WaitGroup{}
	loop := 100

	for i := 0; i < loop; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			node := p.Node()
			mux.Lock()
			distribution[node.Host]++
			mux.Unlock()
		}(i)
	}
	wg.Wait()

	m := int(loop / len(ns))
	for j, d := range distribution {
		if math.Abs(float64(d-m)) > math.Round(float64(loop/len(ns))) {
			t.Errorf("deviation too large: host %s count %d mean %d", j, d, m)
		}
	}
}

func nodes() []lilium.Node {
	return []lilium.Node{
		lilium.Node{Scheme: "http", Host: "127.0.0.1"},
		lilium.Node{Scheme: "http", Host: "127.0.0.2"},
		lilium.Node{Scheme: "http", Host: "127.0.0.3"},
	}
}
