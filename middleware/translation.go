package middleware

import (
	"gateway-micro/public"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"regexp"
	"strings"
)

// TranslationMiddleware 设置Translation
func TranslationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//参照：https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go

		//设置支持语言
		en := en.New()
		zh := zh.New()

		//设置国际化翻译器
		uni := ut.New(zh, zh, en)
		val := validator.New()

		//根据参数取翻译器实例
		locale := c.DefaultQuery("locale", "zh")
		trans, _ := uni.GetTranslator(locale)

		//翻译器注册到validator
		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("en_comment")
			})
		default:
			zh_translations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("comment")
			})

			//自定义验证方法
			val.RegisterValidation("valid_username", func(fl validator.FieldLevel) bool {
				return fl.Field().String() == "admin"
			})
			val.RegisterValidation("valid_service_name", func(fl validator.FieldLevel) bool {
				match, _ := regexp.Match(`^[a-zA-Z0-9_]{6,128}$`, []byte(fl.Field().String()))
				return match
			})
			val.RegisterValidation("valid_rule", func(fl validator.FieldLevel) bool {
				match, _ := regexp.Match(`^\S+$`, []byte(fl.Field().String()))
				return match
			})
			val.RegisterValidation("valid_url_rewrite", func(fl validator.FieldLevel) bool {
				if fl.Field().String() == "" {
					return true
				}
				for _, s := range strings.Split(fl.Field().String(), ",") {
					if len(strings.Split(s, " ")) != 2 {
						return false
					}
				}
				return true
			})
			val.RegisterValidation("valid_header_transfer", func(fl validator.FieldLevel) bool {
				if fl.Field().String() == "" {
					return true
				}
				for _, s := range strings.Split(fl.Field().String(), ",") {
					split := strings.Split(s, " ")
					if split[0] == "add" && len(split) == 3 {
						continue
					}
					if split[0] == "del" && len(split) == 2 {
						continue
					}
					if split[0] == "edit" && len(split) == 2 {
						continue
					}
					return false
				}
				return true
			})
			val.RegisterValidation("valid_ip_list", func(fl validator.FieldLevel) bool {
				for _, s := range strings.Split(fl.Field().String(), ",") {
					if match, _ := regexp.Match(`^\S+:\d+$`, []byte(s)); !match {
						return false
					}
				}
				return true
			})
			val.RegisterValidation("valid_weight_list", func(fl validator.FieldLevel) bool {
				for _, s := range strings.Split(fl.Field().String(), ",") {
					if match, _ := regexp.Match(`^\d+$`, []byte(s)); !match {
						return false
					}
				}
				return true
			})
			val.RegisterValidation("valid_list", func(fl validator.FieldLevel) bool {
				if fl.Field().String() == "" {
					return true
				}
				for _, item := range strings.Split(fl.Field().String(), ",") {
					matched, _ := regexp.Match(`\S+`, []byte(item))
					if !matched {
						return false
					}
				}
				return true
			})

			//自定义翻译器
			val.RegisterTranslation("valid_username", trans, func(ut ut.Translator) error {
				return ut.Add("valid_username", "{0} 填写不正确哦", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_username", fe.Field())
				return t
			})
			val.RegisterTranslation("valid_service_name", trans, func(ut ut.Translator) error {
				return ut.Add("valid_service_name", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_service_name", fe.Field())
				return t
			})
			val.RegisterTranslation("valid_rule", trans, func(ut ut.Translator) error {
				return ut.Add("valid_rule", "{0} 必须是非空字符", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_rule", fe.Field())
				return t
			})
			val.RegisterTranslation("valid_url_rewrite", trans, func(ut ut.Translator) error {
				return ut.Add("valid_url_rewrite", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_url_rewrite", fe.Field())
				return t
			})
			val.RegisterTranslation("valid_header_transfer", trans, func(ut ut.Translator) error {
				return ut.Add("valid_header_transfer", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_header_transfer", fe.Field())
				return t
			})
			val.RegisterTranslation("valid_ip_list", trans, func(ut ut.Translator) error {
				return ut.Add("valid_ip_list", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_ip_list", fe.Field())
				return t
			})
			val.RegisterTranslation("valid_weight_list", trans, func(ut ut.Translator) error {
				return ut.Add("valid_weight_list", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_weight_list", fe.Field())
				return t
			})
			val.RegisterTranslation("valid_list", trans, func(ut ut.Translator) error {
				return ut.Add("valid_ip_list", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_ip_list", fe.Field())
				return t
			})
		}
		c.Set(public.TranslatorKey, trans)
		c.Set(public.ValidatorKey, val)
		c.Next()
	}
}
