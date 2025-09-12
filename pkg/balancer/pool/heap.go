package pool

import (
	"container/heap"
	"fmt"
	"load-balancer/pkg/balancer/node"
)

func (h *NodeHeap) RemoveNode(target *node.Node) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for i, n := range h.heap {
		if n == target {
			last := len(h.heap) - 1
			h.heap[i] = h.heap[last]
			h.heap = h.heap[:last]
			if i < len(h.heap) {
				heap.Fix(h, i)
			}
			return
		}
	}
}

func (h *NodeHeap) Len() int {
	return len(h.heap)
}

func (h *NodeHeap) Less(i, j int) bool {
	return (h.heap[i].Metrics.Connections + uint32(h.heap[i].Queue.Len())) < (h.heap[j].Metrics.Connections + uint32(h.heap[j].Queue.Len()))
}

func (h *NodeHeap) Swap(i, j int) {
	h.heap[i], h.heap[j] = h.heap[j], h.heap[i]
}

func (h *NodeHeap) Push(x any) {
	h.heap = append(h.heap, x.(*node.Node))
}

func (h *NodeHeap) Pop() any {
	old := h.heap
	n := len(old)
	x := old[n-1]
	h.heap = old[0 : n-1]
	return x
}

func (h *NodeHeap) Add(n *node.Node) {
	h.mu.Lock()
	defer h.mu.Unlock()
	heap.Push(h, n)
}

func (h *NodeHeap) RemoveMin() (*node.Node, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.Len() == 0 {
		return nil, fmt.Errorf("cannot pop from empty heap")
	}
	return heap.Pop(h).(*node.Node), nil
}
