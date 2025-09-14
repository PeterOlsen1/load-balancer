package batch

// Adds an item to the batch
//
// If the batch is full, flush it
func (b *Batch[T]) Add(item T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.batch) == cap(b.batch)-1 {
		b.batch = append(b.batch, item)
		b.flushUnsafe()
		return
	}

	b.batch = append(b.batch, item)
}

// Same as flush method but with no locking.
//
// Done to avoid deadlock in the add method when we add + flush
func (b *Batch[T]) flushUnsafe() {
	//copy batch and apply fush function
	batchCopy := append([]T(nil), b.batch...)
	b.onFlush(batchCopy)

	b.batch = make([]T, 0, b.cap)
}

// The flush method applies the `onFlush` function to every item that was in the batch.
//
// Once this is done, reset the batch to have no members in it
func (b *Batch[T]) Flush() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.flushUnsafe()
}

// This method is basically the same as Flush, but the user can provide
// a custom method to be applied to eveyr item
func (b *Batch[T]) FlushCustom(onFlush func([]T)) {
	b.mu.Lock()
	defer b.mu.Unlock()

	//apply custom flush func
	batchCopy := append([]T(nil), b.batch...)
	onFlush(batchCopy)

	b.batch = make([]T, 0, b.cap)
}

// Close the batch. All tickers and channels will be closed
func (b *Batch[T]) Close() {
	b.closeChan <- struct{}{}
}
