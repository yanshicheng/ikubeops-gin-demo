package settings

import (
	"fmt"
	"github.com/yanshicheng/ikubeops-gin-demo/global"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// InitTrans 初始化翻译器
func InitTrans() (err error) {
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		//// 为SignUpparam注册自定义校验方法
		//for _, val := range global.SOG_VALIDATORREG {
		//	v.RegisterStructValidation(val.Func, val.Type)
		//}

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 handler 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		global.IkubeopsTrans, ok = uni.GetTranslator(global.C.App.Language)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", global.C.App.Language)
		}

		// 注册翻译器
		switch global.C.App.Language {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, global.IkubeopsTrans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, global.IkubeopsTrans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, global.IkubeopsTrans)
		}
		return
	}
	return
}

func RemoveTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
