package pool_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/acim/lilium/pool"
)

func Test(t *testing.T) {
	pool := pool.New(nodes())
	pool.Run()
	for i := 0; i < 10; i++ {
		node := <-pool.Out
		fmt.Printf("%d: %s\n", i, node.String())
	}
	time.Sleep(2 * time.Second)
	pool.In <- nodesStatus()
	time.Sleep(2 * time.Second)
	for i := 0; i < 10; i++ {
		node := <-pool.Out
		fmt.Printf("%d: %s\n", i, node.String())
	}
	time.Sleep(2 * time.Second)
}

func nodes() []pool.Node {
	return []pool.Node{
		pool.Node{Scheme: "http", Host: "127.0.0.1"},
		pool.Node{Scheme: "http", Host: "127.0.0.2"},
		pool.Node{Scheme: "http", Host: "127.0.0.3"},
	}
}

func nodesStatus() map[pool.Node]bool {
	status := map[pool.Node]bool{}
	s := false
	for _, node := range nodes() {
		rs := rand1()
		status[node] = rs
		s = s || rs
	}
	if s == false {
		status[pool.Node{Scheme: "http", Host: "127.0.0.1"}] = true
	}
	return status
}

func rand1() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Float32() < 0.5
}
