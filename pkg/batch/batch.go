package batch

import (
	"fmt"
	"load-balancer/pkg/types"
)

func (b *Batch) Add(c *types.Connection) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.batch) == cap(b.batch) {
		return fmt.Errorf("batch is at capacity")
	}

	b.batch = append(b.batch, c)
	return nil
}

func (b *Batch) Flush() []*types.Connection {
    b.mu.Lock()
    defer b.mu.Unlock()

    flushed := b.batch
    b.batch = make([]*types.Connection, 0, b.cap)
    return flushed
}
