package batch

import (
	"sync"
)

type Batch[T any] struct {
	// The underlying sice that holds batch data
	batch []T

	// Mutex for safe locking of slice in goroutines
	mu sync.Mutex

	// Capacity of the batch, need to keep track for when batch is reset
	cap uint32

	// Signal when the batch is full, will be emptied after
	fullChan chan struct{}

	// Signal when the batch is closed
	closeChan chan struct{}

	// Function passed in from the init method
	//
	// When flushing, this function will be applied to all
	// members of the batch
	onFlush func([]T)
}
