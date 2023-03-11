package dao

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccessControl struct {
	ID                int64  `json:"id" gorm:"primary_key"`
	ServiceID         int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	OpenAuth          int    `json:"open_auth" gorm:"column:open_auth" description:"是否开启权限 1=开启"`
	BlackList         string `json:"black_list" gorm:"column:black_list" description:"黑名单ip"`
	WhiteList         string `json:"white_list" gorm:"column:white_list" description:"白名单ip"`
	WhiteHostName     string `json:"white_host_name" gorm:"column:white_host_name" description:"白名单主机"`
	ClientIpFlowLimit int    `json:"client_ip_flow_limit" gorm:"column:client_ip_flow_limit" description:"客户端ip限流"`
	ServiceFlowLimit  int    `json:"service_flow_limit" gorm:"column:service_flow_limit" description:"服务端限流"`
}

func (t *AccessControl) TableName() string {
	return "gateway_service_access_control"
}

func (t *AccessControl) First(c *gin.Context, db *gorm.DB, search *AccessControl) (*AccessControl, error) {
	model := &AccessControl{}
	err := db.WithContext(c).Where(search).First(model).Error
	return model, err
}

func (t *AccessControl) Save(c *gin.Context, db *gorm.DB) error {
	return db.WithContext(c).Save(t).Error
}

func (t *AccessControl) Updates(c *gin.Context, db *gorm.DB) error {
	data, _ := json.Marshal(&t)
	m := make(map[string]interface{})
	json.Unmarshal(data, &m)
	return db.WithContext(c).Model(&AccessControl{}).Where("id = ?", t.ID).Updates(&m).Error
}
