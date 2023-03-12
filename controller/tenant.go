package controller

import (
	"errors"
	"gateway-micro/common/lib"
	"gateway-micro/dao"
	"gateway-micro/dto"
	"gateway-micro/middleware"
	"gateway-micro/public"
	"github.com/gin-gonic/gin"
	"time"
)

type TenantController struct{}

func TenantRegister(group *gin.RouterGroup) {
	tenant := &TenantController{}
	group.GET("/tenant_list", tenant.TenantList)
	group.GET("/tenant_detail", tenant.TenantDetail)
	group.GET("/tenant_statistics", tenant.TenantStatistics)
	group.GET("/tenant_delete", tenant.TenantDelete)
	group.POST("/tenant_add", tenant.TenantAdd)
	group.POST("/tenant_update", tenant.TenantUpdate)
}

// TenantList godoc
// @Summary      租户列表
// @Description  租户列表
// @Tags         租户管理
// @Accept       json
// @Produce      json
// @Param        info 		query 	string  false	"关键词"
// @Param        page_no 	query 	int  	true	"当前页数"
// @Param        page_size	query	int  	true	"每页条数"
// @Success      200		{object}	middleware.Response{data=dto.TenantListOutput}
// @Router       /tenant/tenant_list	[get]
func (tenant *TenantController) TenantList(c *gin.Context) {
	params := &dto.TenantListInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tenantInfo := &dao.Tenant{}
	list, total, err := tenantInfo.PageList(c, lib.GORMDefaultPool, params)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	var outputList []dto.TenantListItemOutput
	for _, item := range list {
		outputList = append(outputList, dto.TenantListItemOutput{
			ID:       item.ID,
			TenantID: item.TenantID,
			Name:     item.Name,
			Secret:   item.Secret,
			WhiteIPS: item.WhiteIPS,
			Qps:      item.Qpd,
			Qpd:      item.Qps,
			RealQps:  0,
			RealQpd:  0,
		})
	}
	output := dto.TenantListOutput{
		Total: total,
		List:  outputList,
	}
	middleware.ResponseSuccess(c, output)
}

// TenantDetail godoc
// @Summary      租户详情
// @Description  租户详情
// @Tags         租户管理
// @Accept       json
// @Produce      json
// @Param        id		query	int64	true	"关键词"
// @Success      200	{object}	middleware.Response{data=dao.Tenant}
// @Router       /tenant/tenant_detail	[get]
func (tenant *TenantController) TenantDetail(c *gin.Context) {
	params := &dto.TenantInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tenantSearch := &dao.Tenant{ID: params.ID}
	tenantDetail, err := tenantSearch.First(c, lib.GORMDefaultPool, tenantSearch)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	middleware.ResponseSuccess(c, tenantDetail)
}

// TenantStatistics godoc
// @Summary      租户统计
// @Description  租户统计
// @Tags         租户管理
// @Accept       json
// @Produce      json
// @Param        id		query	int64	true	"关键词"
// @Success      200	{object}	middleware.Response{data=dto.TenantStatisticsOutput}
// @Router       /tenant/tenant_statistics	[get]
func (tenant *TenantController) TenantStatistics(c *gin.Context) {
	params := &dto.TenantInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	var todayList []int64
	todayTime := time.Now()
	for i := 0; i <= todayTime.Hour(); i++ {
		todayList = append(todayList, 0)
	}

	var yesterdayList []int64
	//yesterdayTime := todayTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i <= 23; i++ {
		yesterdayList = append(yesterdayList, 0)
	}

	middleware.ResponseSuccess(c, &dto.TenantStatisticsOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	})
}

// TenantDelete godoc
// @Summary      租户删除
// @Description  租户删除
// @Tags         租户管理
// @Accept       json
// @Produce      json
// @Param        id		query	int64	true	"关键词"
// @Success      200	{object}	middleware.Response{data=string}
// @Router       /tenant/tenant_delete	[get]
func (tenant *TenantController) TenantDelete(c *gin.Context) {
	params := &dto.TenantInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tenantSearch := &dao.Tenant{ID: params.ID}
	tenantInfo, err := tenantSearch.First(c, lib.GORMDefaultPool, tenantSearch)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	tenantInfo.IsDelete = 1
	if err := tenantInfo.Save(c, lib.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	middleware.ResponseSuccess(c, "租戶删除成功")
}

// TenantAdd godoc
// @Summary      租户添加
// @Description  租户添加
// @Tags         租户管理
// @Accept       json
// @Produce      json
// @Param        body	body		dto.TenantAddInput 	true	"body"
// @Success      200	{object}	middleware.Response{data=string}
// @Router       /tenant/tenant_add	[post]
func (tenant *TenantController) TenantAdd(c *gin.Context) {
	params := &dto.TenantAddInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tenantSearch := &dao.Tenant{
		TenantID: params.TenantID,
		IsDelete: 0,
	}
	if _, err := tenantSearch.First(c, lib.GORMDefaultPool, tenantSearch); err == nil {
		middleware.ResponseError(c, 2001, errors.New("租户ID被占用，请重新输入"))
		return
	}

	if params.Secret == "" {
		params.Secret = public.MD5(params.TenantID)
	}

	tx := lib.GORMDefaultPool
	tenantInfo := &dao.Tenant{
		TenantID: params.TenantID,
		Name:     params.Name,
		Secret:   params.Secret,
		WhiteIPS: params.WhiteIPS,
		Qps:      params.Qps,
		Qpd:      params.Qpd,
	}
	if err := tenantInfo.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2002, err)
		return
	}

	middleware.ResponseSuccess(c, "租戶添加成功")
}

// TenantUpdate godoc
// @Summary      租户更新
// @Description  租户更新
// @Tags         租户管理
// @Accept       json
// @Produce      json
// @Param        body	body		dto.TenantUpdateInput	true	"body"
// @Success      200	{object}	middleware.Response{data=string}
// @Router       /tenant/tenant_update	[post]
func (tenant *TenantController) TenantUpdate(c *gin.Context) {
	params := &dto.TenantUpdateInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tenantSearch := &dao.Tenant{
		TenantID: params.TenantID,
	}

	tenantInfo, err := tenantSearch.First(c, lib.GORMDefaultPool, tenantSearch)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	if params.Secret == "" {
		params.Secret = public.MD5(params.TenantID)
	}

	tenantInfo.Name = params.Name
	tenantInfo.Secret = params.Secret
	tenantInfo.WhiteIPS = params.WhiteIPS
	tenantInfo.Qps = params.Qps
	tenantInfo.Qpd = params.Qpd
	if err := tenantInfo.Save(c, lib.GORMDefaultPool); err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	middleware.ResponseSuccess(c, "租户更新成功")
}
