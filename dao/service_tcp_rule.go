package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TcpRule struct {
	ID        int64 `json:"id" gorm:"primary_key"`
	ServiceID int64 `json:"service_id" gorm:"column:service_id" description:"服务id	"`
	Port      int   `json:"port" gorm:"column:port" description:"端口	"`
}

func (t *TcpRule) TableName() string {
	return "gateway_service_tcp_rule"
}

func (t *TcpRule) Find(c *gin.Context, db *gorm.DB, search *TcpRule) (*TcpRule, error) {
	model := &TcpRule{}
	err := db.WithContext(c).Where(search).Find(model).Error
	return model, err
}

func (t *TcpRule) Save(c *gin.Context, db *gorm.DB) error {
	if err := db.WithContext(c).Save(t).Error; err != nil {
		return err
	}
	return nil
}
