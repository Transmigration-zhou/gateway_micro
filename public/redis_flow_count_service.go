package public

import (
	"fmt"
	"gateway-micro/common/lib"
	"github.com/gomodule/redigo/redis"
	"sync/atomic"
	"time"
)

type RedisFlowCountService struct {
	TenantID    string
	Interval    time.Duration
	QPS         int64
	Unix        int64
	TickerCount int64
	TotalCount  int64
}

func NewRedisFlowCountService(tenantID string, interval time.Duration) *RedisFlowCountService {
	service := &RedisFlowCountService{
		TenantID: tenantID,
		Interval: interval,
		QPS:      0,
		Unix:     0,
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		ticker := time.NewTicker(interval)
		for {
			<-ticker.C
			tickerCount := atomic.LoadInt64(&service.TickerCount) //获取数据
			atomic.StoreInt64(&service.TickerCount, 0)            //重置数据

			currentTime := time.Now()
			dayKey := service.GetDayKey(currentTime)
			hourKey := service.GetHourKey(currentTime)
			if err := RedisConfPipeline(func(c redis.Conn) {
				c.Send("INCRBY", dayKey, tickerCount)
				c.Send("EXPIRE", dayKey, 86400*2)
				c.Send("INCRBY", hourKey, tickerCount)
				c.Send("EXPIRE", hourKey, 86400*2)
			}); err != nil {
				fmt.Println("RedisConfPipeline err", err)
				continue
			}

			totalCount, err := service.GetDayData(currentTime)
			if err != nil {
				fmt.Println("service.GetDayData err", err)
				continue
			}

			nowUnix := time.Now().Unix()
			if service.Unix == 0 {
				service.Unix = time.Now().Unix()
				continue
			}

			tickerCount = totalCount - service.TotalCount
			if nowUnix > service.Unix {
				service.TotalCount = totalCount
				service.QPS = tickerCount / (nowUnix - service.Unix)
				service.Unix = time.Now().Unix()
			}
		}
	}()
	return service
}

func (service *RedisFlowCountService) GetDayKey(t time.Time) string {
	dayStr := t.In(lib.TimeLocation).Format("20060102")
	return fmt.Sprintf("%s_%s_%service", RedisFlowDayKey, dayStr, service.TenantID)
}

func (service *RedisFlowCountService) GetHourKey(t time.Time) string {
	hourStr := t.In(lib.TimeLocation).Format("2006010215")
	return fmt.Sprintf("%s_%s_%service", RedisFlowHourKey, hourStr, service.TenantID)
}

func (service *RedisFlowCountService) GetHourData(t time.Time) (int64, error) {
	return redis.Int64(RedisConfDo("GET", service.GetHourKey(t)))
}

func (service *RedisFlowCountService) GetDayData(t time.Time) (int64, error) {
	return redis.Int64(RedisConfDo("GET", service.GetDayKey(t)))
}

// Increase 原子增加
func (service *RedisFlowCountService) Increase() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		atomic.AddInt64(&service.TickerCount, 1)
	}()
}
