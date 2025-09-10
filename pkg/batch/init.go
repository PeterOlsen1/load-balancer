package batch

import "load-balancer/pkg/types"

func InitBatch(cap int) *Batch {
	return &Batch{
		batch: make([]*types.Connection, 0, cap),
		cap:   cap,
	}
}
