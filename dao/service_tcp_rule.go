package dao

import (
	"encoding/json"
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

func (t *TcpRule) First(c *gin.Context, db *gorm.DB, search *TcpRule) (*TcpRule, error) {
	model := &TcpRule{}
	err := db.WithContext(c).Where(search).First(model).Error
	return model, err
}

func (t *TcpRule) Save(c *gin.Context, db *gorm.DB) error {
	return db.WithContext(c).Save(t).Error
}

func (t *TcpRule) Updates(c *gin.Context, db *gorm.DB) error {
	data, _ := json.Marshal(&t)
	m := make(map[string]interface{})
	json.Unmarshal(data, &m)
	return db.WithContext(c).Model(&TcpRule{}).Where("id = ?", t.ID).Updates(&m).Error
}
