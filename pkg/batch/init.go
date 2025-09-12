package batch

import "load-balancer/pkg/types"

func InitBatch(cap uint32) *Batch {
	return &Batch{
		batch: make([]*types.Connection, 0, cap),
		cap:   cap,
	}
}
