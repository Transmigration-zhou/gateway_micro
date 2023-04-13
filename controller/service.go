package controller

import (
	"fmt"
	"gateway-micro/common/lib"
	"gateway-micro/dao"
	"gateway-micro/dto"
	"gateway-micro/middleware"
	"gateway-micro/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type ServiceController struct{}

func ServiceRegister(group *gin.RouterGroup) {
	service := &ServiceController{}
	group.GET("/service_list", service.ServiceList)
	group.GET("/service_delete", service.ServiceDelete)
	group.GET("/service_detail", service.ServiceDetail)
	group.GET("/service_statistics", service.ServiceStatistics)
	group.POST("/service_add_http", service.ServiceAddHTTP)
	group.POST("/service_update_http", service.ServiceUpdateHTTP)
	group.POST("/service_add_tcp", service.ServiceAddTcp)
	group.POST("/service_update_tcp", service.ServiceUpdateTcp)
	group.POST("/service_add_grpc", service.ServiceAddGrpc)
	group.POST("/service_update_grpc", service.ServiceUpdateGrpc)
}

// ServiceList godoc
// @Summary      服务列表
// @Description  服务列表
// @Tags         服务管理
// @Accept       json
// @Produce      json
// @Param        info 		query 	string  false	"关键词"
// @Param        page_no 	query 	int  	true	"当前页数"
// @Param        page_size	query	int  	true	"每页条数"
// @Success      200		{object}	middleware.Response{data=dto.ServiceListOutput}
// @Router       /service/service_list	[get]
func (service *ServiceController) ServiceList(c *gin.Context) {
	params := &dto.ServiceListInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	db, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	serviceInfo := &dao.ServiceInfo{}
	list, total, err := serviceInfo.PageList(c, db, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	var outList []dto.ServiceListItemOutput
	for _, listItem := range list {
		serviceDetail, err := listItem.ServiceDetail(c, db, &listItem)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			return
		}
		serviceAddr := "unknown"
		clusterIp := lib.GetStringConf("base.cluster.cluster_ip")
		clusterPort := lib.GetStringConf("base.cluster.cluster_port")
		clusterSSLPort := lib.GetStringConf("base.cluster.cluster_ssl_port")

		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 1 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIp, clusterSSLPort, serviceDetail.HTTPRule.Rule)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 0 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIp, clusterPort, serviceDetail.HTTPRule.Rule)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
			serviceAddr = serviceDetail.HTTPRule.Rule
		}
		if serviceDetail.Info.LoadType == public.LoadTypeTCP {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIp, serviceDetail.TCPRule.Port)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeGRPC {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIp, serviceDetail.GRPCRule.Port)
		}

		ipList := serviceDetail.LoadBalance.GetIpListByModel()
		counter, err := public.FlowCounterHandler.GetCounter(public.FlowServicePrefix + listItem.ServiceName)
		if err != nil {
			middleware.ResponseError(c, 2004, err)
			return
		}

		outItem := dto.ServiceListItemOutput{
			ID:          listItem.ID,
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServiceDesc,
			LoadType:    listItem.LoadType,
			ServiceAddr: serviceAddr,
			Qps:         counter.QPS,
			Qpd:         counter.TotalCount,
			TotalNode:   len(ipList),
		}
		outList = append(outList, outItem)
	}
	out := &dto.ServiceListOutput{
		Total: total,
		List:  outList,
	}
	middleware.ResponseSuccess(c, out)
}

// ServiceDelete godoc
// @Summary      服务删除
// @Description  服务删除
// @Tags         服务管理
// @Accept       json
// @Produce      json
// @Param        id		query	int64	true	"关键词"
// @Success      200	{object}	middleware.Response{data=string}
// @Router       /service/service_delete	[get]
func (service *ServiceController) ServiceDelete(c *gin.Context) {
	params := &dto.ServiceInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	db, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	serviceInfo := &dao.ServiceInfo{ID: params.ID}
	serviceInfo, err = serviceInfo.First(c, db, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	serviceInfo.IsDelete = 1
	if err := serviceInfo.Save(c, db); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	middleware.ResponseSuccess(c, "服务删除成功")
}

// ServiceDetail godoc
// @Summary      服务详情
// @Description  服务详情
// @Tags         服务管理
// @Accept       json
// @Produce      json
// @Param        id		query	int64	true	"关键词"
// @Success      200	{object}	middleware.Response{data=dao.ServiceDetail}
// @Router       /service/service_detail	[get]
func (service *ServiceController) ServiceDetail(c *gin.Context) {
	params := &dto.ServiceInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	db, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	serviceInfo := &dao.ServiceInfo{ID: params.ID}
	serviceInfo, err = serviceInfo.First(c, db, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	serviceDetail, err := serviceInfo.ServiceDetail(c, db, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	middleware.ResponseSuccess(c, serviceDetail)
}

// ServiceStatistics godoc
// @Summary      服务统计
// @Description  服务统计
// @Tags         服务管理
// @Accept       json
// @Produce      json
// @Param        id		query	int64	true	"关键词"
// @Success      200	{object}	middleware.Response{data=dto.ServiceStatisticsOutput}
// @Router       /service/service_statistics	[get]
func (service *ServiceController) ServiceStatistics(c *gin.Context) {
	params := &dto.ServiceInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	db, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	serviceInfo := &dao.ServiceInfo{ID: params.ID}
	serviceInfo, err = serviceInfo.First(c, db, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	serviceDetail, err := serviceInfo.ServiceDetail(c, db, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	counter, err := public.FlowCounterHandler.GetCounter(public.FlowServicePrefix + serviceDetail.Info.ServiceName)
	if err != nil {
		middleware.ResponseError(c, 2004, err)
		return
	}

	var todayList []int64
	todayTime := time.Now()
	for i := 0; i <= todayTime.Hour(); i++ {
		dateTime := time.Date(todayTime.Year(), todayTime.Month(), todayTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourData, _ := counter.GetHourData(dateTime)
		todayList = append(todayList, hourData)
	}

	var yesterdayList []int64
	yesterdayTime := todayTime.Add(-24 * time.Hour)
	for i := 0; i <= 23; i++ {
		dateTime := time.Date(yesterdayTime.Year(), yesterdayTime.Month(), yesterdayTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourData, _ := counter.GetHourData(dateTime)
		yesterdayList = append(yesterdayList, hourData)
	}

	middleware.ResponseSuccess(c, &dto.ServiceStatisticsOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	})
}

// ServiceAddHTTP godoc
// @Summary      添加HTTP服务
// @Description  添加HTTP服务
// @Tags         服务管理
// @Accept       json
// @Produce      json
// @Param        body	body		dto.ServiceAddHTTPInput	true	"body"
// @Success      200	{object}	middleware.Response{data=string}
// @Router       /service/service_add_http	[post]
func (service *ServiceController) ServiceAddHTTP(c *gin.Context) {
	params := &dto.ServiceAddHTTPInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2001, errors.New("IP列表与权重列表数量不一致"))
		return
	}

	db, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	tx := db.Begin()

	serviceInfo := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	if _, err := serviceInfo.First(c, tx, serviceInfo); err == nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, errors.New("服务名已经存在"))
		return
	}

	httpUrl := &dao.HttpRule{RuleType: params.RuleType, Rule: params.Rule}
	if _, err := httpUrl.First(c, tx, httpUrl); err == nil {
		tx.Rollback()
		middleware.ResponseError(c, 2004, errors.New("服务接入前缀或域名已存在"))
		return
	}

	serviceModel := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := serviceModel.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}

	httpRule := &dao.HttpRule{
		ServiceID:      serviceModel.ID,
		RuleType:       params.RuleType,
		Rule:           params.Rule,
		NeedHttps:      params.NeedHttps,
		NeedStripUri:   params.NeedStripUri,
		NeedWebsocket:  params.NeedWebsocket,
		UrlRewrite:     params.UrlRewrite,
		HeaderTransfer: params.HeaderTransfer,
	}
	if err := httpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}

	accessControl := &dao.AccessControl{
		ServiceID:         serviceModel.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		ClientIpFlowLimit: params.ClientIpFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	loadBalance := &dao.LoadBalance{
		ServiceID:              serviceModel.ID,
		RoundType:              params.RoundType,
		IpList:                 params.IpList,
		WeightList:             params.WeightList,
		UpstreamConnectTimeout: params.UpstreamConnectTimeout,
		UpstreamHeaderTimeout:  params.UpstreamHeaderTimeout,
		UpstreamIdleTimeout:    params.UpstreamIdleTimeout,
		UpstreamMaxIdle:        params.UpstreamMaxIdle,
	}
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}

	tx.Commit()
	middleware.ResponseSuccess(c, "添加HTTP服务成功")
}

// ServiceUpdateHTTP godoc
// @Summary      更新HTTP服务
// @Description  更新HTTP服务
// @Tags         服务管理
// @Accept       json
// @Produce      json
// @Param        body	body		dto.ServiceUpdateHTTPInput	true	"body"
// @Success      200	{object}	middleware.Response{data=string}
// @Router       /service/service_update_http	[post]
func (service *ServiceController) ServiceUpdateHTTP(c *gin.Context) {
	params := &dto.ServiceUpdateHTTPInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2001, errors.New("IP列表与权重列表数量不一致"))
		return
	}

	db, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	tx := db.Begin()

	serviceInfo := &dao.ServiceInfo{ID: params.ID}
	serviceDetail, err := serviceInfo.ServiceDetail(c, lib.GORMDefaultPool, serviceInfo)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, err)
		return
	}

	info := serviceDetail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2004, err)
		return
	}

	httpRule := serviceDetail.HTTPRule
	httpRule.NeedHttps = params.NeedHttps
	httpRule.NeedStripUri = params.NeedStripUri
	httpRule.NeedWebsocket = params.NeedWebsocket
	httpRule.UrlRewrite = params.UrlRewrite
	httpRule.HeaderTransfer = params.HeaderTransfer

	if err := httpRule.Updates(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}

	accessControl := serviceDetail.AccessControl
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.ClientIpFlowLimit = params.ClientIpFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := accessControl.Updates(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}

	loadBalance := serviceDetail.LoadBalance
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.UpstreamConnectTimeout = params.UpstreamConnectTimeout
	loadBalance.UpstreamHeaderTimeout = params.UpstreamHeaderTimeout
	loadBalance.UpstreamIdleTimeout = params.UpstreamIdleTimeout
	loadBalance.UpstreamMaxIdle = params.UpstreamMaxIdle
	if err := loadBalance.Updates(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	tx.Commit()
	middleware.ResponseSuccess(c, "更新HTTP服务成功")
}

// ServiceAddTcp godoc
// @Summary      添加TCP服务
// @Description  添加TCP服务
// @Tags         服务管理
// @Accept       json
// @Produce      json
// @Param        body	body		dto.ServiceAddTcpInput	true	"body"
// @Success      200	{object}	middleware.Response{data=string}
// @Router       /service/service_add_tcp	[post]
func (service *ServiceController) ServiceAddTcp(c *gin.Context) {
	params := &dto.ServiceAddTcpInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	serviceInfo := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	if _, err := serviceInfo.First(c, lib.GORMDefaultPool, serviceInfo); err == nil {
		middleware.ResponseError(c, 2002, errors.New("服务名被占用，请重新输入"))
		return
	}

	tcpPort := &dao.TcpRule{
		Port: params.Port,
	}
	if _, err := tcpPort.First(c, lib.GORMDefaultPool, tcpPort); err == nil {
		middleware.ResponseError(c, 2003, errors.New("服务端口被占用，请重新输入"))
		return
	}

	grpcPort := &dao.GrpcRule{
		Port: params.Port,
	}
	if _, err := grpcPort.First(c, lib.GORMDefaultPool, grpcPort); err == nil {
		middleware.ResponseError(c, 2004, errors.New("服务端口被占用，请重新输入"))
		return
	}

	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2005, errors.New("IP列表与权重列表数量不一致"))
		return
	}

	tx := lib.GORMDefaultPool.Begin()

	serviceModel := &dao.ServiceInfo{
		LoadType:    public.LoadTypeTCP,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := serviceModel.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}

	tcpRule := &dao.TcpRule{
		ServiceID: serviceModel.ID,
		Port:      params.Port,
	}
	if err := tcpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}

	accessControl := &dao.AccessControl{
		ServiceID:         serviceModel.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientIpFlowLimit: params.ClientIpFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2009, err)
		return
	}

	loadBalance := &dao.LoadBalance{
		ServiceID:  serviceModel.ID,
		RoundType:  params.RoundType,
		IpList:     params.IpList,
		WeightList: params.WeightList,
		ForbidList: params.ForbidList,
	}
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	tx.Commit()
	middleware.ResponseSuccess(c, "添加TCP服务成功")
}

// ServiceUpdateTcp godoc
// @Summary      更新TCP服务
// @Description  更新TCP服务
// @Tags         服务管理
// @Accept       json
// @Produce      json
// @Param        body	body		dto.ServiceUpdateTcpInput	true	"body"
// @Success      200	{object}	middleware.Response{data=string}
// @Router       /service/service_update_tcp	[post]
func (service *ServiceController) ServiceUpdateTcp(c *gin.Context) {
	params := &dto.ServiceUpdateTcpInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2002, errors.New("IP列表与权重列表数量不一致"))
		return
	}

	tx := lib.GORMDefaultPool.Begin()

	serviceInfo := &dao.ServiceInfo{ID: params.ID}
	serviceDetail, err := serviceInfo.ServiceDetail(c, lib.GORMDefaultPool, serviceInfo)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, err)
		return
	}

	info := serviceDetail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2004, err)
		return
	}

	tcpRule := &dao.TcpRule{}
	tcpRule = serviceDetail.TCPRule
	tcpRule.ServiceID = info.ID
	tcpRule.Port = params.Port
	if err := tcpRule.Updates(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}

	accessControl := &dao.AccessControl{}
	accessControl = serviceDetail.AccessControl
	accessControl.ServiceID = info.ID
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.WhiteHostName = params.WhiteHostName
	accessControl.ClientIpFlowLimit = params.ClientIpFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := accessControl.Updates(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}

	loadBalance := &dao.LoadBalance{}
	loadBalance = serviceDetail.LoadBalance
	loadBalance.ServiceID = info.ID
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList
	if err := loadBalance.Updates(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	tx.Commit()
	middleware.ResponseSuccess(c, "更新TCP服务成功")
}

// ServiceAddGrpc godoc
// @Summary      添加GRPC服务
// @Description  添加GRPC服务
// @Tags         服务管理
// @Accept       json
// @Produce      json
// @Param        body	body		dto.ServiceAddGrpcInput	true	"body"
// @Success      200	{object}	middleware.Response{data=string}
// @Router       /service/service_add_grpc	[post]
func (service *ServiceController) ServiceAddGrpc(c *gin.Context) {
	params := &dto.ServiceAddGrpcInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	serviceInfo := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	if _, err := serviceInfo.First(c, lib.GORMDefaultPool, serviceInfo); err == nil {
		middleware.ResponseError(c, 2002, errors.New("服务名被占用，请重新输入"))
		return
	}

	tcpPort := &dao.TcpRule{
		Port: params.Port,
	}
	if _, err := tcpPort.First(c, lib.GORMDefaultPool, tcpPort); err == nil {
		middleware.ResponseError(c, 2003, errors.New("服务端口被占用，请重新输入"))
		return
	}

	grpcPort := &dao.GrpcRule{
		Port: params.Port,
	}
	if _, err := grpcPort.First(c, lib.GORMDefaultPool, grpcPort); err == nil {
		middleware.ResponseError(c, 2004, errors.New("服务端口被占用，请重新输入"))
		return
	}

	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2005, errors.New("IP列表与权重列表数量不一致"))
		return
	}

	tx := lib.GORMDefaultPool.Begin()

	serviceModel := &dao.ServiceInfo{
		LoadType:    public.LoadTypeGRPC,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := serviceModel.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}

	grpcRule := &dao.GrpcRule{
		ServiceID:      serviceModel.ID,
		Port:           params.Port,
		HeaderTransfer: params.HeaderTransfer,
	}
	if err := grpcRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}

	accessControl := &dao.AccessControl{
		ServiceID:         serviceModel.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientIpFlowLimit: params.ClientIpFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2009, err)
		return
	}

	loadBalance := &dao.LoadBalance{
		ServiceID:  serviceModel.ID,
		RoundType:  params.RoundType,
		IpList:     params.IpList,
		WeightList: params.WeightList,
		ForbidList: params.ForbidList,
	}
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	tx.Commit()
	middleware.ResponseSuccess(c, "添加GRPC服务成功")
}

// ServiceUpdateGrpc godoc
// @Summary      更新GRPC服务
// @Description  更新GRPC服务
// @Tags         服务管理
// @Accept       json
// @Produce      json
// @Param        body	body		dto.ServiceUpdateGrpcInput	true	"body"
// @Success      200	{object}	middleware.Response{data=string}
// @Router       /service/service_update_grpc	[post]
func (service *ServiceController) ServiceUpdateGrpc(c *gin.Context) {
	params := &dto.ServiceUpdateGrpcInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2002, errors.New("IP列表与权重列表数量不一致"))
		return
	}

	tx := lib.GORMDefaultPool.Begin()

	serviceInfo := &dao.ServiceInfo{ID: params.ID}
	serviceDetail, err := serviceInfo.ServiceDetail(c, lib.GORMDefaultPool, serviceInfo)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, err)
		return
	}

	info := serviceDetail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2004, err)
		return
	}

	grpcRule := &dao.TcpRule{}
	grpcRule = serviceDetail.TCPRule
	grpcRule.ServiceID = info.ID
	grpcRule.Port = params.Port
	if err := grpcRule.Updates(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}

	accessControl := &dao.AccessControl{}
	accessControl = serviceDetail.AccessControl
	accessControl.ServiceID = info.ID
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.WhiteHostName = params.WhiteHostName
	accessControl.ClientIpFlowLimit = params.ClientIpFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := accessControl.Updates(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}

	loadBalance := &dao.LoadBalance{}
	loadBalance = serviceDetail.LoadBalance
	loadBalance.ServiceID = info.ID
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList
	if err := loadBalance.Updates(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	tx.Commit()
	middleware.ResponseSuccess(c, "更新GRPC服务成功")
}
