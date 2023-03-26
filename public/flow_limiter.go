package public

import (
	"golang.org/x/time/rate"
	"sync"
)

var FlowLimiterHandler *FlowLimiter

func init() {
	FlowLimiterHandler = NewFlowLimiter()
}

type FlowLimiter struct {
	FlowLimiterMap   map[string]*FlowLimiterItem
	FlowLimiterSlice []*FlowLimiterItem
	Lock             sync.Mutex
}

type FlowLimiterItem struct {
	ServiceName string
	Limiter     *rate.Limiter
}

func NewFlowLimiter() *FlowLimiter {
	return &FlowLimiter{
		FlowLimiterMap:   map[string]*FlowLimiterItem{},
		FlowLimiterSlice: []*FlowLimiterItem{},
		Lock:             sync.Mutex{},
	}
}

func (counter *FlowLimiter) GetLimiter(serverName string, qps float64) (*rate.Limiter, error) {
	for _, item := range counter.FlowLimiterSlice {
		if item.ServiceName == serverName {
			return item.Limiter, nil
		}
	}

	newLimiter := rate.NewLimiter(rate.Limit(qps), int(qps*3))
	item := &FlowLimiterItem{
		ServiceName: serverName,
		Limiter:     newLimiter,
	}
	counter.Lock.Lock()
	defer counter.Lock.Unlock()
	counter.FlowLimiterMap[serverName] = item
	counter.FlowLimiterSlice = append(counter.FlowLimiterSlice, item)

	return newLimiter, nil
}
