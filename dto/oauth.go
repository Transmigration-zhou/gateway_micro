package dto

import (
	"gateway-micro/public"
	"github.com/gin-gonic/gin"
)

type OAuthTokenInput struct {
	GrantType string `json:"grant_type" form:"grant_type" comment:"授权类型" validate:"required"` //授权类型
	Scope     string `json:"scope" form:"scope" comment:"权限范围" validate:"required"`           //权限范围
}

func (param *OAuthTokenInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type OAuthTokenOutput struct {
	AccessToken string `json:"access_token" form:"access_token"`
	ExpiresIn   int    `json:"expires_in" form:"expires_in"`
	TokenType   string `json:"token_type" form:"token_type"`
	Scope       string `json:"scope" form:"scope"`
}
