package public

import (
	"sync"
	"time"
)

var FlowCounterHandler *FlowCounter

func init() {
	FlowCounterHandler = NewFlowCounter()
}

type FlowCounter struct {
	RedisFlowCountMap   map[string]*RedisFlowCountService
	RedisFlowCountSlice []*RedisFlowCountService
	Lock                sync.Mutex
}

func NewFlowCounter() *FlowCounter {
	return &FlowCounter{
		RedisFlowCountMap:   map[string]*RedisFlowCountService{},
		RedisFlowCountSlice: []*RedisFlowCountService{},
		Lock:                sync.Mutex{},
	}
}

func (counter *FlowCounter) GetCounter(serverName string) (*RedisFlowCountService, error) {
	for _, item := range counter.RedisFlowCountSlice {
		if item.TenantID == serverName {
			return item, nil
		}
	}

	newCounter := NewRedisFlowCountService(serverName, 1*time.Second)
	counter.Lock.Lock()
	defer counter.Lock.Unlock()
	counter.RedisFlowCountMap[serverName] = newCounter
	counter.RedisFlowCountSlice = append(counter.RedisFlowCountSlice, newCounter)

	return newCounter, nil
}
