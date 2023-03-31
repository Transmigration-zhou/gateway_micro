package controller

import (
	"encoding/base64"
	"errors"
	"gateway-micro/common/lib"
	"gateway-micro/dao"
	"gateway-micro/dto"
	"gateway-micro/middleware"
	"gateway-micro/public"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

type OAuthController struct{}

func OAuthRegister(group *gin.RouterGroup) {
	oAuth := &OAuthController{}
	group.POST("/token", oAuth.OAuthToken)
}

func (oauth *OAuthController) OAuthToken(c *gin.Context) {
	params := &dto.OAuthTokenInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	split := strings.Split(c.GetHeader("Authorization"), " ")
	if len(split) != 2 {
		middleware.ResponseError(c, 2001, errors.New("用户名或密码格式错误"))
		return
	}
	secret, err := base64.StdEncoding.DecodeString(split[1])
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	//取出 tenant_id secret
	//生成 tenant_list
	//匹配 tenant_id
	//基于 jwt生成token
	//生成 output
	part := strings.Split(string(secret), ":")
	if len(part) != 2 {
		middleware.ResponseError(c, 2003, errors.New("用户名或密码格式错误"))
		return
	}

	tenantList := dao.TenantManagerHandler.GetTenantList()
	for _, tenant := range tenantList {
		if tenant.TenantID == part[0] && tenant.Secret == part[1] {
			claims := jwt.RegisteredClaims{
				Issuer:    tenant.TenantID,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(public.JwtExpires * time.Second).In(lib.TimeLocation)),
			}
			token, err := public.JwtEncode(claims)
			if err != nil {
				middleware.ResponseError(c, 2004, err)
				return
			}
			out := &dto.OAuthTokenOutput{
				AccessToken: token,
				ExpiresIn:   public.JwtExpires,
				TokenType:   "Bearer",
				Scope:       "read_write",
			}
			middleware.ResponseSuccess(c, out)
			return
		}
	}

	middleware.ResponseError(c, 2005, errors.New("未匹配正确tenant信息"))
}
