package controller

import (
	"gateway-micro/common/lib"
	"gateway-micro/dao"
	"gateway-micro/dto"
	"gateway-micro/middleware"
	"gateway-micro/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"time"
)

type DashboardController struct{}

func DashboardRegister(group *gin.RouterGroup) {
	dashboard := &DashboardController{}
	group.GET("/panel_group_data", dashboard.PanelGroupData)
	group.GET("/flow_statistics", dashboard.FlowStatistics)
	group.GET("/service_statistics", dashboard.ServiceStatistics)
}

// PanelGroupData godoc
// @Summary      指标统计
// @Description  指标统计
// @Tags         首页大盘
// @Accept       json
// @Produce      json
// @Success      200		{object}	middleware.Response{data=dto.PanelGroupDataOutput}
// @Router       /dashboard/panel_group_data	[get]
func (dashboard *DashboardController) PanelGroupData(c *gin.Context) {
	db, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	serviceInfo := &dao.ServiceInfo{}
	_, serviceNum, err := serviceInfo.PageList(c, db, &dto.ServiceListInput{PageSize: 1, PageNo: 1})
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	tenant := &dao.Tenant{}
	_, tenantNum, err := tenant.PageList(c, db, &dto.TenantListInput{PageNo: 1, PageSize: 1})
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	counter, err := public.FlowCounterHandler.GetCounter(public.FlowTotal)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	out := &dto.PanelGroupDataOutput{
		ServiceNum:      serviceNum,
		TenantNum:       tenantNum,
		TodayRequestNum: counter.TotalCount,
		CurrentQPS:      counter.QPS,
	}
	middleware.ResponseSuccess(c, out)
}

// FlowStatistics godoc
// @Summary      流量统计
// @Description  流量统计
// @Tags         首页大盘
// @Accept       json
// @Produce      json
// @Success      200		{object}	middleware.Response{data=dto.DashboardStatisticsOutput}
// @Router       /dashboard/flow_statistics	[get]
func (dashboard *DashboardController) FlowStatistics(c *gin.Context) {
	counter, err := public.FlowCounterHandler.GetCounter(public.FlowTotal)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
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

	middleware.ResponseSuccess(c, &dto.DashboardStatisticsOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	})
}

// ServiceStatistics godoc
// @Summary      服务统计
// @Description  服务统计
// @Tags         首页大盘
// @Accept       json
// @Produce      json
// @Success      200		{object}	middleware.Response{data=dto.DashServiceStatisticsOutput}
// @Router       /dashboard/service_statistics	[get]
func (dashboard *DashboardController) ServiceStatistics(c *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	serviceInfo := &dao.ServiceInfo{}
	list, err := serviceInfo.GroupByLoadType(c, tx)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	var legend []string
	for index, item := range list {
		name, ok := public.LoadTypeMap[item.LoadType]
		if !ok {
			middleware.ResponseError(c, 2002, errors.New("load_type not found"))
			return
		}
		list[index].Name = name
		legend = append(legend, name)
	}

	out := &dto.DashServiceStatisticsOutput{
		Legend: legend,
		Series: list,
	}
	middleware.ResponseSuccess(c, out)
}
