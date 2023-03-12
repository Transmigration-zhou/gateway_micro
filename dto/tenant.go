package dto

import (
	"gateway-micro/public"
	"github.com/gin-gonic/gin"
	"time"
)

type TenantListInput struct {
	Info     string `json:"info" form:"info" comment:"关键词" validate:""`                                  //关键词
	PageNo   int    `json:"page_no" form:"page_no" comment:"页数" validate:"required,min=1,max=999"`       //页数
	PageSize int    `json:"page_size" form:"page_size" comment:"每页条数" validate:"required,min=1,max=999"` //每页条数
}

func (param *TenantListInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type TenantListOutput struct {
	Total int64                  `json:"total" form:"total" comment:"租户总数"` //租户总数
	List  []TenantListItemOutput `json:"list" form:"list" comment:"租户列表"`   //租户列表
}

type TenantListItemOutput struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	TenantID  string    `json:"tenant_id" gorm:"column:tenant_id" description:"租户id"`
	Name      string    `json:"name" gorm:"column:name" description:"租户名称"`
	Secret    string    `json:"secret" gorm:"column:secret" description:"密钥"`
	WhiteIPS  string    `json:"white_ips" gorm:"column:white_ips" description:"ip白名单，支持前缀匹配"`
	Qps       int64     `json:"qps" gorm:"column:qps" description:"每秒请求量限制"`
	Qpd       int64     `json:"qpd" gorm:"column:qpd" description:"日请求量限制"`
	RealQps   int64     `json:"real_qps" description:"实际每秒请求量"`
	RealQpd   int64     `json:"real_qpd" description:"实际日请求量"`
	UpdatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"添加时间"`
	CreatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	IsDelete  int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除 0：否 1：是"`
}

type TenantInput struct {
	ID int64 `json:"id" form:"id" comment:"服务ID" example:"56" validate:"required"` //租户ID
}

func (param *TenantInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type TenantStatisticsOutput struct {
	Today     []int64 `json:"today" form:"today" comment:"今日流量" validate:"required"`         //今日流量
	Yesterday []int64 `json:"yesterday" form:"yesterday" comment:"昨日流量" validate:"required"` //昨日流量
}

type TenantAddInput struct {
	TenantID string `json:"tenant_id" form:"tenant_id" comment:"租户id" example:"" validate:"required"`
	Name     string `json:"name" form:"name" comment:"租户名称" example:"" validate:"required"`
	Secret   string `json:"secret" form:"secret" comment:"密钥" example:"" validate:""`
	WhiteIPS string `json:"white_ips" form:"white_ips" comment:"ip白名单，支持前缀匹配" example:""`
	Qps      int64  `json:"qps" form:"qps" comment:"每秒请求量限制" validate:""`
	Qpd      int64  `json:"qpd" form:"qpd" comment:"日请求量限制" validate:""`
}

func (param *TenantAddInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type TenantUpdateInput struct {
	ID       int64  `json:"id" form:"id" gorm:"column:id" comment:"主键ID" validate:"required"`
	TenantID string `json:"tenant_id" form:"tenant_id" gorm:"column:tenant_id" comment:"租户id" example:"" validate:""`
	Name     string `json:"name" form:"name" gorm:"column:name" comment:"租户名称" example:"" validate:"required"`
	Secret   string `json:"secret" form:"secret" gorm:"column:secret" example:"" comment:"密钥" validate:"required"`
	WhiteIPS string `json:"white_ips" form:"white_ips" gorm:"column:white_ips" comment:"ip白名单，支持前缀匹配" example:""`
	Qps      int64  `json:"qps" form:"qps" gorm:"column:qps" comment:"每秒请求量限制"`
	Qpd      int64  `json:"qpd" form:"qpd" gorm:"column:qpd" comment:"日请求量限制"`
}

func (param *TenantUpdateInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}
