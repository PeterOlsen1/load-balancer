package batch

import "time"

func InitBatch[T any](cap uint32, flushInterval time.Duration, onFlush func([]T)) *Batch[T] {
	b := &Batch[T]{
		batch:     make([]T, 0, cap),
		cap:       cap,
		onFlush:   onFlush,
		fullChan:  make(chan struct{}),
		closeChan: make(chan struct{}),
	}

	batchTicker := time.NewTicker(flushInterval)
	go func() {
		for {
			select {
			// batch has closed, flush all items and close channels
			case <-b.closeChan:
				b.Flush()
				batchTicker.Stop()
				close(b.closeChan)
				close(b.fullChan)
				return

			// ticker went off or batch is full. flush those items!
			case <-b.fullChan:
			case <-batchTicker.C:
				b.Flush()
			}
		}
	}()

	return b
}
