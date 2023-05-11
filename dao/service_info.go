package dao

import (
	"gateway-micro/dto"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type ServiceInfo struct {
	ID          int64     `json:"id" gorm:"primary_key"`
	LoadType    int       `json:"load_type" gorm:"column:load_type" description:"负载类型 0=http 1=tcp 2=grpc"`
	ServiceName string    `json:"service_name" gorm:"column:service_name" description:"服务名称"`
	ServiceDesc string    `json:"service_desc" gorm:"column:service_desc" description:"服务描述"`
	UpdatedAt   time.Time `json:"create_at" gorm:"column:create_at" description:"更新时间"`
	CreatedAt   time.Time `json:"update_at" gorm:"column:update_at" description:"添加时间"`
	IsDelete    int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除 0：否 1：是"`
}

func (t *ServiceInfo) TableName() string {
	return "gateway_service_info"
}

func (t *ServiceInfo) ServiceDetail(c *gin.Context, db *gorm.DB, search *ServiceInfo) (*ServiceDetail, error) {
	if search.ServiceName == "" {
		info, err := t.First(c, db, search)
		if err != nil {
			return nil, err
		}
		search = info
	}
	httpRule := &HttpRule{ServiceID: search.ID}
	httpRule, err := httpRule.First(c, db, httpRule)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	tcpRule := &TcpRule{ServiceID: search.ID}
	tcpRule, err = tcpRule.First(c, db, tcpRule)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	grpcRule := &GrpcRule{ServiceID: search.ID}
	grpcRule, err = grpcRule.First(c, db, grpcRule)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	accessControl := &AccessControl{ServiceID: search.ID}
	accessControl, err = accessControl.First(c, db, accessControl)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	loadBalance := &LoadBalance{ServiceID: search.ID}
	loadBalance, err = loadBalance.First(c, db, loadBalance)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	detail := &ServiceDetail{
		Info:          search,
		HTTPRule:      httpRule,
		TCPRule:       tcpRule,
		GRPCRule:      grpcRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}
	return detail, nil
}

func (t *ServiceInfo) PageList(c *gin.Context, db *gorm.DB, param *dto.ServiceListInput) ([]ServiceInfo, int64, error) {
	total := int64(0)
	var list []ServiceInfo
	offset := (param.PageNo - 1) * param.PageSize

	query := db.WithContext(c).Table(t.TableName()).Where("is_delete = 0")
	if param.Info != "" {
		query = query.Where("(service_name like ? or service_desc like ?)", "%"+param.Info+"%", "%"+param.Info+"%")
	}
	query.Count(&total)
	if err := query.Limit(param.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	query.Limit(param.PageSize).Offset(offset)
	return list, total, nil
}

func (t *ServiceInfo) First(c *gin.Context, db *gorm.DB, search *ServiceInfo) (*ServiceInfo, error) {
	out := &ServiceInfo{}
	err := db.WithContext(c).Where(search).First(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *ServiceInfo) Save(c *gin.Context, db *gorm.DB) error {
	return db.WithContext(c).Save(t).Error
}

func (t *ServiceInfo) GroupByLoadType(c *gin.Context, db *gorm.DB) ([]dto.DashServiceStatisticsItemOutput, error) {
	var list []dto.DashServiceStatisticsItemOutput
	query := db.WithContext(c)
	if err := query.Table(t.TableName()).Where("is_delete = 0").Select("load_type, count(*) as value").Group("load_type").Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
