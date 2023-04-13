package dao

import (
	"gateway-micro/common/lib"
	"gateway-micro/dto"
	"gateway-micro/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http/httptest"
	"strings"
	"sync"
)

type ServiceDetail struct {
	Info          *ServiceInfo   `json:"info" description:"基本信息"`
	HTTPRule      *HttpRule      `json:"http_rule" description:"http_rule"`
	TCPRule       *TcpRule       `json:"tcp_rule" description:"tcp_rule"`
	GRPCRule      *GrpcRule      `json:"grpc_rule" description:"grpc_rule"`
	LoadBalance   *LoadBalance   `json:"load_balance" description:"load_balance"`
	AccessControl *AccessControl `json:"access_control" description:"access_control"`
}

var ServiceManagerHandler *ServiceManager

func init() {
	ServiceManagerHandler = NewServiceManager()
}

type ServiceManager struct {
	ServiceMap   map[string]*ServiceDetail
	ServiceSlice []*ServiceDetail
	Lock         sync.Mutex
	init         sync.Once
	err          error
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		ServiceMap:   map[string]*ServiceDetail{},
		ServiceSlice: []*ServiceDetail{},
		Lock:         sync.Mutex{},
		init:         sync.Once{},
	}
}

func (s *ServiceManager) LoadOnce() error {
	s.init.Do(func() {
		serviceInfo := &ServiceInfo{}
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		db, err := lib.GetGormPool("default")
		if err != nil {
			s.err = err
			return
		}
		params := &dto.ServiceListInput{PageNo: 1, PageSize: 99999}
		list, _, err := serviceInfo.PageList(c, db, params)
		if err != nil {
			s.err = err
			return
		}
		s.Lock.Lock()
		defer s.Lock.Unlock()
		for _, listItem := range list {
			tmpItem := listItem
			serviceDetail, err := tmpItem.ServiceDetail(c, db, &tmpItem)
			if err != nil {
				s.err = err
				return
			}
			s.ServiceMap[listItem.ServiceName] = serviceDetail
			s.ServiceSlice = append(s.ServiceSlice, serviceDetail)
		}
	})
	return s.err
}

func (s *ServiceManager) HTTPAccessMode(c *gin.Context) (*ServiceDetail, error) {
	//1、前缀匹配 /abc ==> serviceSlice.rule
	//2、域名匹配 www.test.com ==> serviceSlice.rule
	//host c.Request.Host
	//path c.Request.URL.Path
	host := c.Request.Host
	host = host[0:strings.Index(host, ":")]
	//fmt.Println("host:", host)
	path := c.Request.URL.Path
	for _, serviceItem := range s.ServiceSlice {
		if serviceItem.Info.LoadType != public.LoadTypeHTTP {
			continue
		}
		if serviceItem.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
			if serviceItem.HTTPRule.Rule == host {
				return serviceItem, nil
			}
		}
		if serviceItem.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL {
			if strings.HasPrefix(path, serviceItem.HTTPRule.Rule) {
				return serviceItem, nil
			}
		}
	}
	return nil, errors.New("not matched service")
}
