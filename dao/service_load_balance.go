package dao

import (
	"encoding/json"
	"fmt"
	"gateway-micro/public"
	"gateway-micro/reverse_proxy/load_balance"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type LoadBalance struct {
	ID            int64  `json:"id" gorm:"primary_key"`
	ServiceID     int64  `json:"service_id" gorm:"column:service_id" description:"服务id	"`
	CheckMethod   int    `json:"check_method" gorm:"column:check_method" description:"检查方法 tcpchk=检测端口是否握手成功	"`
	CheckTimeout  int    `json:"check_timeout" gorm:"column:check_timeout" description:"check超时时间	"`
	CheckInterval int    `json:"check_interval" gorm:"column:check_interval" description:"检查间隔, 单位s		"`
	RoundType     int    `json:"round_type" gorm:"column:round_type" description:"轮询方式 round/weight_round/random/ip_hash"`
	IpList        string `json:"ip_list" gorm:"column:ip_list" description:"ip列表"`
	WeightList    string `json:"weight_list" gorm:"column:weight_list" description:"权重列表"`
	ForbidList    string `json:"forbid_list" gorm:"column:forbid_list" description:"禁用ip列表"`

	UpstreamConnectTimeout int `json:"upstream_connect_timeout" gorm:"column:upstream_connect_timeout" description:"下游建立连接超时, 单位s"`
	UpstreamHeaderTimeout  int `json:"upstream_header_timeout" gorm:"column:upstream_header_timeout" description:"下游获取header超时, 单位s	"`
	UpstreamIdleTimeout    int `json:"upstream_idle_timeout" gorm:"column:upstream_idle_timeout" description:"下游链接最大空闲时间, 单位s	"`
	UpstreamMaxIdle        int `json:"upstream_max_idle" gorm:"column:upstream_max_idle" description:"下游最大空闲链接数"`
}

func (t *LoadBalance) TableName() string {
	return "gateway_service_load_balance"
}

func (t *LoadBalance) First(c *gin.Context, db *gorm.DB, search *LoadBalance) (*LoadBalance, error) {
	model := &LoadBalance{}
	err := db.WithContext(c).Where(search).First(model).Error
	return model, err
}

func (t *LoadBalance) Save(c *gin.Context, db *gorm.DB) error {
	return db.WithContext(c).Save(t).Error
}

func (t *LoadBalance) Updates(c *gin.Context, db *gorm.DB) error {
	data, _ := json.Marshal(&t)
	m := make(map[string]interface{})
	json.Unmarshal(data, &m)
	return db.WithContext(c).Model(&LoadBalance{}).Where("id = ?", t.ID).Updates(&m).Error
}

func (t *LoadBalance) GetIpListByModel() []string {
	return strings.Split(t.IpList, ",")
}

func (t *LoadBalance) GetWeightListByModel() []string {
	return strings.Split(t.WeightList, ",")
}

var LoadBalancerHandler *LoadBalancer

func init() {
	LoadBalancerHandler = NewLoadBalancer()
}

type LoadBalancerItem struct {
	LoadBalance load_balance.LoadBalance
	ServiceName string
}

type LoadBalancer struct {
	LoadBalanceMap   map[string]*LoadBalancerItem
	LoadBalanceSlice []*LoadBalancerItem
	Lock             sync.Mutex
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		LoadBalanceMap:   map[string]*LoadBalancerItem{},
		LoadBalanceSlice: []*LoadBalancerItem{},
		Lock:             sync.Mutex{},
	}
}

func (l *LoadBalancer) GetLoadBalancer(service *ServiceDetail) (load_balance.LoadBalance, error) {
	for _, lbItem := range l.LoadBalanceSlice {
		if lbItem.ServiceName == service.Info.ServiceName {
			return lbItem.LoadBalance, nil
		}
	}

	schema := ""
	if service.Info.LoadType == public.LoadTypeHTTP {
		schema = "http://"
		if service.HTTPRule.NeedHttps == 1 {
			schema = "https://"
		}
	}

	ipList := service.LoadBalance.GetIpListByModel()
	weightList := service.LoadBalance.GetWeightListByModel()
	ipConf := map[string]string{}
	for ipIndex, ipItem := range ipList {
		ipConf[ipItem] = weightList[ipIndex]
	}

	mConf, err := load_balance.NewLoadBalanceCheckConf(fmt.Sprintf("%s%s", schema, "%s"), ipConf)
	if err != nil {
		return nil, err
	}

	lb := load_balance.LoadBalanceFactorWithConf(load_balance.LbType(service.LoadBalance.RoundType), mConf)
	lbItem := &LoadBalancerItem{
		LoadBalance: lb,
		ServiceName: service.Info.ServiceName,
	}
	l.Lock.Lock()
	defer l.Lock.Unlock()
	l.LoadBalanceMap[service.Info.ServiceName] = lbItem
	l.LoadBalanceSlice = append(l.LoadBalanceSlice, lbItem)

	return lb, nil
}

var TransporterHandler *Transporter

func init() {
	TransporterHandler = NewTransporter()
}

type TransportItem struct {
	Trans       *http.Transport
	ServiceName string
}

type Transporter struct {
	TransportMap   map[string]*TransportItem
	TransportSlice []*TransportItem
	Lock           sync.Mutex
}

func NewTransporter() *Transporter {
	return &Transporter{
		TransportMap:   map[string]*TransportItem{},
		TransportSlice: []*TransportItem{},
		Lock:           sync.Mutex{},
	}
}

func (t *Transporter) GetTransporter(service *ServiceDetail) (*http.Transport, error) {
	for _, transItem := range t.TransportSlice {
		if transItem.ServiceName == service.Info.ServiceName {
			return transItem.Trans, nil
		}
	}

	trans := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: time.Duration(service.LoadBalance.UpstreamConnectTimeout) * time.Second,
		}).DialContext,
		MaxIdleConns:          service.LoadBalance.UpstreamMaxIdle,
		IdleConnTimeout:       time.Duration(service.LoadBalance.UpstreamIdleTimeout) * time.Second,
		ResponseHeaderTimeout: time.Duration(service.LoadBalance.UpstreamHeaderTimeout) * time.Second,
	}

	transItem := &TransportItem{
		Trans:       trans,
		ServiceName: service.Info.ServiceName,
	}
	t.Lock.Lock()
	defer t.Lock.Unlock()
	t.TransportMap[service.Info.ServiceName] = transItem
	t.TransportSlice = append(t.TransportSlice, transItem)

	return trans, nil
}
