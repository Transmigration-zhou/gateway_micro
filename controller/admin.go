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
)

type AdminController struct{}

func AdminRegister(group *gin.RouterGroup) {
	admin := &AdminController{}
	group.GET("/admin_info", admin.AdminInfo)
	group.POST("/change_pwd", admin.ChangePwd)
}

// AdminInfo godoc
// @Summary      管理员信息
// @Description  管理员信息
// @Tags         管理员接口
// @Accept       json
// @Produce      json
// @Success      200	{object}	middleware.Response{data=dto.AdminInfoOutput}
// @Router       /admin/admin_info	[get]
func (adminLogin *AdminController) AdminInfo(c *gin.Context) {
	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(sessInfo.(string)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	out := &dto.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		Name:         adminSessionInfo.UserName,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "",
		Introduction: "我是管理员",
		Roles:        []string{"admin"},
	}
	middleware.ResponseSuccess(c, out)
}

// ChangePwd godoc
// @Summary      修改密码
// @Description  修改密码
// @Tags         管理员接口
// @Accept       json
// @Produce      json
// @Param        body 	body	dto.ChangePwdInput	true	"body"
// @Success      200  	{object}	middleware.Response{data=string}
// @Router       /admin/change_pwd	[post]
func (adminLogin *AdminController) ChangePwd(c *gin.Context) {
	params := &dto.ChangePwdInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(sessInfo.(string)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	db, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	admin := &dao.Admin{}
	admin, err = admin.First(c, db, &dao.Admin{UserName: adminSessionInfo.UserName})
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	//生成新密码
	saltPassword := public.GenSaltPassword(admin.Salt, params.Password)
	admin.Password = saltPassword

	//执行数据保存
	if err := admin.Save(c, db); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	middleware.ResponseSuccess(c, "修改密码成功")
}
