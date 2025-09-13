package workerpool

func (p *WorkerPool[T]) Event(e T) {
	p.eventChan <- e
}

func (p *WorkerPool[T]) Close() {
	close(p.eventChan)
}

func (p *WorkerPool[T]) UpdateEventHandler(f func(T)) {
	p.eventHandler = f
}
