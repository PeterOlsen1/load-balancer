package batch

import (
	"load-balancer/pkg/types"
	"sync"
)

type Batch struct {
	batch []*types.Connection
	mu    sync.Mutex
	cap   uint32
}
