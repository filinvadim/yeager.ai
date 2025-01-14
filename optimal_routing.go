package main

import (
	"container/heap"
	"fmt"
	"math"
)

type Router struct {
	liter   string
	latency float64
}

type State struct {
	liter   string
	latency float64
	path    string // just for visualisation and comfort
}

type PriorityQueue []State

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].latency < pq[j].latency
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(State))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func findMinimumLatencyPath(
	graph map[string][]Router,
	compressionNodes []string,
	source, destination string,
) float64 {
	var (
		compressedSet = make(map[string]struct{})
		latencyMap    = make(map[string]float64)
	)
	for _, node := range compressionNodes {
		compressedSet[node] = struct{}{}
	}

	for node := range graph {
		latencyMap[node] = math.Inf(1)
	}

	latencyMap[source] = 0

	queue := &PriorityQueue{}
	heap.Init(queue)
	heap.Push(queue, State{liter: source, latency: 0, path: source})

	var (
		minLatency = math.Inf(1)
		bestPath   string
	)

	for queue.Len() > 0 {
		current := heap.Pop(queue).(State)

		// find a best latency for every possible way
		if current.liter == destination {
			if current.latency < minLatency {
				minLatency = current.latency
				bestPath = current.path
			}
		}

		for _, router := range graph[current.liter] {
			newPath := current.path + router.liter

			newLatency := current.latency + router.latency

			if _, ok := compressedSet[current.liter]; ok {
				newLatency = current.latency + router.latency/2 // compress new hop
			}

			latencyMap[router.liter] = newLatency

			heap.Push(
				queue,
				State{
					liter:   router.liter,
					latency: newLatency,
					path:    newPath,
				},
			)

		}
	}

	if minLatency == math.Inf(1) {
		return math.Inf(1)
	}

	fmt.Println("BEST PATH", bestPath)
	return minLatency
}
