package dto

import (
	"gateway-micro/public"
	"github.com/gin-gonic/gin"
)

type ServiceListInput struct {
	Info     string `json:"info" form:"info" comment:"关键词" example:"" validate:""`                                    //关键词
	PageNo   int    `json:"page_no" form:"page_no" comment:"页数" example:"1" validate:"required,min=1,max=999"`        //页数
	PageSize int    `json:"page_size" form:"page_size" comment:"每页条数" example:"20" validate:"required,min=1,max=999"` //每页条数
}

func (param *ServiceListInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceListOutput struct {
	Total int64                   `json:"total" form:"total" comment:"总数"` //总数
	List  []ServiceListItemOutput `json:"list" form:"list" comment:"列表"`   //列表
}

type ServiceListItemOutput struct {
	ID          int64  `json:"id" form:"id"`                     //id
	ServiceName string `json:"service_name" form:"service_name"` //服务名称
	ServiceDesc string `json:"service_desc" form:"service_desc"` //服务描述
	LoadType    int    `json:"load_type" form:"load_type"`       //类型
	ServiceAddr string `json:"service_addr" form:"service_addr"` //服务地址
	Qps         int64  `json:"qps" form:"qps"`                   //每秒请求量
	Qpd         int64  `json:"qpd" form:"qpd"`                   //日请求量
	TotalNode   int    `json:"total_node" form:"total_node"`     //节点数
}

type ServiceInput struct {
	ID int64 `json:"id" form:"id" comment:"服务ID" example:"56" validate:"required"` //服务ID
}

func (param *ServiceInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceAddHTTPInput struct {
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名" example:"" validate:"required,valid_service_name"` //服务名
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"" validate:"required,min=1,max=255"`     //服务描述

	RuleType       int    `json:"rule_type" form:"rule_type" comment:"接入类型" validate:"min=0,max=1"`                                      //接入类型
	Rule           string `json:"rule" form:"rule" comment:"接入路径：域名或者前缀" example:"" validate:"required,valid_rule"`                      //域名或者前缀
	NeedHttps      int    `json:"need_https" form:"need_https" comment:"支持https" validate:"min=0,max=1"`                                 //支持https
	NeedStripUri   int    `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri" validate:"min=0,max=1"`                     //启用strip_uri
	NeedWebsocket  int    `json:"need_websocket" form:"need_websocket" comment:"是否支持websocket" validate:"min=0,max=1"`                   //是否支持websocket
	UrlRewrite     string `json:"url_rewrite" form:"url_rewrite" comment:"url重写功能" example:"" validate:"valid_url_rewrite"`              //url重写功能
	HeaderTransfer string `json:"header_transfer" form:"header_transfer" comment:"header转换" example:"" validate:"valid_header_transfer"` //header转换

	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限" validate:"min=0,max=1"`                  //是否开启权限
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单ip" example:"" validate:"valid_list"`       //黑名单ip
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单ip" example:"" validate:"valid_list"`       //白名单ip
	ClientIpFlowLimit int    `json:"client_ip_flow_limit" form:"client_ip_flow_limit" comment:"客户端ip限流" validate:"min=0"` //客户端ip限流
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:"min=0"`       //服务端限流

	RoundType              int    `json:"round_type" form:"round_type" comment:"轮询方式" validate:"min=0,max=3"`                               //轮询方式
	IpList                 string `json:"ip_list" form:"ip_list" comment:"ip列表" example:"" validate:"required,valid_ip_list"`               //ip列表
	WeightList             string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"" validate:"required,valid_weight_list"`   //权重列表
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" form:"upstream_connect_timeout" comment:"建立连接超时，单位s" validate:"min=0"`   //建立连接超时，单位s
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout" form:"upstream_header_timeout" comment:"获取header超时，单位s" validate:"min=0"` //获取header超时，单位s
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout" form:"upstream_idle_timeout" comment:"链接最大空闲时间，单位s" validate:"min=0"`       //链接最大空闲时间，单位s
	UpstreamMaxIdle        int    `json:"upstream_max_idle" form:"upstream_max_idle" comment:"最大空闲链接数" validate:"min=0"`                    //最大空闲链接数
}

func (param *ServiceAddHTTPInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceUpdateHTTPInput struct {
	ID          int64  `json:"id" form:"id" comment:"服务ID" example:"63" validate:"required,min=1"`                                                     //服务ID
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名" example:"test_http_service_indb" validate:"required,valid_service_name"` //服务名
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"test_http_service_indb" validate:"required,min=1,max=255"`     //服务描述

	RuleType       int    `json:"rule_type" form:"rule_type" comment:"接入类型" validate:"min=0,max=1"`                                        //接入类型
	Rule           string `json:"rule" form:"rule" comment:"接入路径：域名或者前缀" example:"/test_http_service_indb" validate:"required,valid_rule"` //域名或者前缀
	NeedHttps      int    `json:"need_https" form:"need_https" comment:"支持https" validate:"min=0,max=1"`                                   //支持https
	NeedStripUri   int    `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri" validate:"min=0,max=1"`                       //启用strip_uri
	NeedWebsocket  int    `json:"need_websocket" form:"need_websocket" comment:"是否支持websocket" validate:"min=0,max=1"`                     //是否支持websocket
	UrlRewrite     string `json:"url_rewrite" form:"url_rewrite" comment:"url重写功能" example:"" validate:"valid_url_rewrite"`                //url重写功能
	HeaderTransfer string `json:"header_transfer" form:"header_transfer" comment:"header转换" example:"" validate:"valid_header_transfer"`   //header转换

	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限" validate:"min=0,max=1"`                  //是否开启权限
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单ip" example:"" validate:"valid_list"`       //黑名单ip
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单ip" example:"" validate:"valid_list"`       //白名单ip
	ClientIpFlowLimit int    `json:"client_ip_flow_limit" form:"client_ip_flow_limit" comment:"客户端ip限流" validate:"min=0"` //客户端ip限流
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:"min=0"`       //服务端限流

	RoundType              int    `json:"round_type" form:"round_type" comment:"轮询方式" validate:"min=0,max=3"`                               //轮询方式
	IpList                 string `json:"ip_list" form:"ip_list" comment:"ip列表" example:"127.0.0.1:80" validate:"required,valid_ip_list"`   //ip列表
	WeightList             string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"50" validate:"required,valid_weight_list"` //权重列表
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" form:"upstream_connect_timeout" comment:"建立连接超时，单位s" validate:"min=0"`   //建立连接超时，单位s
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout" form:"upstream_header_timeout" comment:"获取header超时，单位s" validate:"min=0"` //获取header超时，单位s
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout" form:"upstream_idle_timeout" comment:"链接最大空闲时间，单位s" validate:"min=0"`       //链接最大空闲时间，单位s
	UpstreamMaxIdle        int    `json:"upstream_max_idle" form:"upstream_max_idle" comment:"最大空闲链接数" validate:"min=0"`                    //最大空闲链接数
}

func (param *ServiceUpdateHTTPInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceStatisticsOutput struct {
	Today     []int64 `json:"today" form:"today" comment:"今日流量" validate:"required"`         //今日流量
	Yesterday []int64 `json:"yesterday" form:"yesterday" comment:"昨日流量" validate:"required"` //昨日流量
}

type ServiceAddTcpInput struct {
	ServiceName    string `json:"service_name" form:"service_name" comment:"服务名" example:""  validate:"required,valid_service_name"`       //服务名
	ServiceDesc    string `json:"service_desc" form:"service_desc" comment:"服务描述" example:""  validate:"required"`                         //服务描述
	Port           int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"`                    //端口，需要设置8001-8999范围内
	HeaderTransfer string `json:"header_transfer" form:"header_transfer" comment:"header头转换" example:""  validate:"valid_header_transfer"` //header头转换

	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:"min=0,max=1"`                            //是否开启权限验证
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" example:"" validate:"valid_list"` //黑名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" example:"" validate:"valid_list"` //白名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机名，以逗号间隔" example:"" validate:"valid_list"`  //白名单主机名，以逗号间隔
	ClientIpFlowLimit int    `json:"client_ip_flow_limit" form:"client_ip_flow_limit" comment:"客户端IP限流" validate:""`                  //客户端IP限流
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`                        //服务端限流
	RoundType         int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`                                         //轮询策略
	IpList            string `json:"ip_list" form:"ip_list" comment:"IP列表" example:"" validate:"required,valid_ip_list"`              //IP列表
	WeightList        string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"" validate:"required,valid_weight_list"`  //权重列表
	ForbidList        string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" example:"" validate:"valid_list"`                //禁用IP列表
}

func (param *ServiceAddTcpInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceUpdateTcpInput struct {
	ID             int64  `json:"id" form:"id" comment:"服务ID" validate:"required"`                                                        //服务ID
	ServiceName    string `json:"service_name" form:"service_name" comment:"服务名" example:"" validate:"required,valid_service_name"`       //服务名
	ServiceDesc    string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"" validate:"required"`                         //服务描述
	Port           int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"`                   //端口，需要设置8001-8999范围内
	HeaderTransfer string `json:"header_transfer" form:"header_transfer" comment:"header头转换" example:"" validate:"valid_header_transfer"` //header头转换

	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:"min=0,max=1"`                            //是否开启权限验证
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" example:"" validate:"valid_list"` //黑名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" example:"" validate:"valid_list"` //白名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机名，以逗号间隔" example:"" validate:"valid_list"`  //白名单主机名，以逗号间隔
	ClientIpFlowLimit int    `json:"client_ip_flow_limit" form:"client_ip_flow_limit" comment:"客户端IP限流" validate:""`                  //客户端IP限流
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`                        //服务端限流
	RoundType         int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`                                         //轮询策略
	IpList            string `json:"ip_list" form:"ip_list" comment:"IP列表" example:"" validate:"required,valid_ip_list"`              //IP列表
	WeightList        string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"" validate:"required,valid_weight_list"`  //权重列表
	ForbidList        string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" example:"" validate:"valid_list"`                //禁用IP列表
}

func (param *ServiceUpdateTcpInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceAddGrpcInput struct {
	ServiceName    string `json:"service_name" form:"service_name" comment:"服务名" example:"" validate:"required,valid_service_name"`        //服务名
	ServiceDesc    string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"" validate:"required"`                          //服务描述
	Port           int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"`                    //端口，需要设置8001-8999范围内
	HeaderTransfer string `json:"header_transfer" form:"header_transfer" comment:"header头转换" example:""  validate:"valid_header_transfer"` //header头转换

	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:"min=0,max=1"`                            //是否开启权限验证
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" example:"" validate:"valid_list"` //黑名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" example:"" validate:"valid_list"` //白名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机名，以逗号间隔" example:"" validate:"valid_list"`  //白名单主机名，以逗号间隔
	ClientIpFlowLimit int    `json:"client_ip_flow_limit" form:"client_ip_flow_limit" comment:"客户端IP限流" validate:""`                  //客户端IP限流
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`                        //服务端限流
	RoundType         int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`                                         //轮询策略
	IpList            string `json:"ip_list" form:"ip_list" comment:"IP列表" example:"" validate:"required,valid_ip_list"`              //IP列表
	WeightList        string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"" validate:"required,valid_weight_list"`  //权重列表
	ForbidList        string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" example:"" validate:"valid_list"`                //禁用IP列表
}

func (param *ServiceAddGrpcInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type ServiceUpdateGrpcInput struct {
	ID             int64  `json:"id" form:"id" comment:"服务ID" validate:"required"`                                                         //服务ID
	ServiceName    string `json:"service_name" form:"service_name" comment:"服务名" example:"" validate:"required,valid_service_name"`        //服务名
	ServiceDesc    string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"" validate:"required"`                          //服务描述
	Port           int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"`                    //端口，需要设置8001-8999范围内
	HeaderTransfer string `json:"header_transfer" form:"header_transfer" comment:"header头转换" example:""  validate:"valid_header_transfer"` //header头转换

	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:"min=0,max=1"`                            //是否开启权限验证
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" example:"" validate:"valid_list"` //黑名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" example:"" validate:"valid_list"` //白名单IP，以逗号间隔，白名单优先级高于黑名单
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机名，以逗号间隔" example:"" validate:"valid_list"`  //白名单主机名，以逗号间隔
	ClientIpFlowLimit int    `json:"client_ip_flow_limit" form:"client_ip_flow_limit" comment:"客户端IP限流" validate:""`                  //客户端IP限流
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`                        //服务端限流
	RoundType         int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`                                         //轮询策略
	IpList            string `json:"ip_list" form:"ip_list" comment:"IP列表" example:"" validate:"required,valid_ip_list"`              //IP列表
	WeightList        string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"" validate:"required,valid_weight_list"`  //权重列表
	ForbidList        string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" example:"" validate:"valid_list"`                //禁用IP列表
}

func (param *ServiceUpdateGrpcInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}
