package dao

import (
	"gateway-micro/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type Tenant struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	TenantID  string    `json:"tenant_id" gorm:"column:tenant_id" description:"租户id"`
	Name      string    `json:"name" gorm:"column:name" description:"租户名称"`
	Secret    string    `json:"secret" gorm:"column:secret" description:"密钥"`
	WhiteIPS  string    `json:"white_ips" gorm:"column:white_ips" description:"ip白名单，支持前缀匹配"`
	Qps       int64     `json:"qps" gorm:"column:qps" description:"每秒请求量限制"`
	Qpd       int64     `json:"qpd" gorm:"column:qpd" description:"日请求量限制"`
	UpdatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"更新时间"`
	CreatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"添加时间"`
	IsDelete  int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除 0：否 1：是"`
}

func (t *Tenant) TableName() string {
	return "gateway_tenant"
}

func (t *Tenant) PageList(c *gin.Context, db *gorm.DB, param *dto.TenantListInput) ([]Tenant, int64, error) {
	total := int64(0)
	var list []Tenant
	offset := (param.PageNo - 1) * param.PageSize

	query := db.WithContext(c).Table(t.TableName()).Where("is_delete = 0")
	if param.Info != "" {
		query = query.Where("(tenant_name like ? or tenant_desc like ?)", "%"+param.Info+"%", "%"+param.Info+"%")
	}
	if err := query.Limit(param.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	query.Limit(param.PageSize).Offset(offset).Count(&total)
	return list, total, nil
}

func (t *Tenant) First(c *gin.Context, db *gorm.DB, search *Tenant) (*Tenant, error) {
	out := &Tenant{}
	err := db.WithContext(c).Where(search).First(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *Tenant) Save(c *gin.Context, db *gorm.DB) error {
	return db.WithContext(c).Save(t).Error
}
