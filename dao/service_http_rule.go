package dao

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HttpRule struct {
	ID             int64  `json:"id" gorm:"primary_key"`
	ServiceID      int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	RuleType       int    `json:"rule_type" gorm:"column:rule_type" description:"匹配类型 domain=域名, url_prefix=url前缀"`
	Rule           string `json:"rule" gorm:"column:rule" description:"type=domain表示域名，type=url_prefix时表示url前缀"`
	NeedHttps      int    `json:"need_https" gorm:"column:need_https" description:"type=支持https 1=支持"`
	NeedWebsocket  int    `json:"need_websocket" gorm:"column:need_websocket" description:"启用websocket 1=启用"`
	NeedStripUri   int    `json:"need_strip_uri" gorm:"column:need_strip_uri" description:"启用strip_uri 1=启用"`
	UrlRewrite     string `json:"url_rewrite" gorm:"column:url_rewrite" description:"url重写功能，每行一个	"`
	HeaderTransfer string `json:"header_transfer" gorm:"column:header_transfer" description:"header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue	"`
}

func (t *HttpRule) TableName() string {
	return "gateway_service_http_rule"
}

func (t *HttpRule) First(c *gin.Context, db *gorm.DB, search *HttpRule) (*HttpRule, error) {
	model := &HttpRule{}
	err := db.WithContext(c).Where(search).First(model).Error
	return model, err
}

func (t *HttpRule) Save(c *gin.Context, db *gorm.DB) error {
	return db.WithContext(c).Save(t).Error
}

func (t *HttpRule) Updates(c *gin.Context, db *gorm.DB) error {
	data, _ := json.Marshal(&t)
	m := make(map[string]interface{})
	json.Unmarshal(data, &m)
	return db.WithContext(c).Model(&HttpRule{}).Where("id = ?", t.ID).Updates(&m).Error
}
