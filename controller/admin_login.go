package controller

import (
	"encoding/json"
	"gateway-micro/common/lib"
	"gateway-micro/dao"
	"gateway-micro/dto"
	"gateway-micro/middleware"
	"gateway-micro/public"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
)

type AdminLoginController struct{}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.AdminLogin)
	group.GET("/logout", adminLogin.AdminLogOut)
}

// AdminLogin godoc
// @Summary      管理员登录
// @Description  管理员登录
// @Tags         管理员接口
// @Accept       json
// @Produce      json
// @Param        body 	body 		dto.AdminLoginInput	true	"body"
// @Success      200  	{object}  	middleware.Response{data=dto.AdminLoginOutput}
// @Router       /admin_login/login	[post]
func (adminLogin *AdminLoginController) AdminLogin(c *gin.Context) {
	params := &dto.AdminLoginInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	db, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	admin := &dao.Admin{}
	admin, err = admin.LoginCheck(c, db, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	//设置session
	sessInfo := &dto.AdminSessionInfo{
		ID:        admin.Id,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}
	sessBts, err := json.Marshal(sessInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	sess := sessions.Default(c)
	sess.Set(public.AdminSessionInfoKey, string(sessBts))
	sess.Save()

	out := &dto.AdminLoginOutput{Token: admin.UserName}
	middleware.ResponseSuccess(c, out)
}

// AdminLogOut godoc
// @Summary      管理员登出
// @Description  管理员登出
// @Tags         管理员接口
// @Accept       json
// @Produce      json
// @Success      200  	{object}	middleware.Response{data=string}
// @Router       /admin_login/logout	[get]
func (adminLogin *AdminLoginController) AdminLogOut(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Delete(public.AdminSessionInfoKey)
	sess.Save()
	middleware.ResponseSuccess(c, "登出成功")
}
