package response

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/yanshicheng/ikubeops-gin-demo/global"
	"github.com/yanshicheng/ikubeops-gin-demo/settings"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Data struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message interface{} `json:"message"`
}

func Success(c *gin.Context, data interface{}) {
	rd := &Data{
		Code:    0,
		Data:    data,
		Message: "",
	}
	c.JSON(http.StatusOK, rd)
	return
}

func SuccessCode(c *gin.Context, code int, data interface{}) {
	rd := &Data{
		Code:    code,
		Data:    data,
		Message: "",
	}
	c.JSON(http.StatusOK, rd)
	return
}

func Failed(c *gin.Context, msg interface{}) {
	rd := &Data{
		Code:    10001,
		Data:    "",
		Message: msg,
	}
	c.JSON(http.StatusOK, rd)
	return
}

func FailServerErr(c *gin.Context, msg interface{}) {
	rd := &Data{
		Code:    10001,
		Data:    "",
		Message: msg,
	}
	c.JSON(http.StatusInternalServerError, rd)
	return
}

func FailedCode(c *gin.Context, code int, msg interface{}) {
	rd := &Data{
		Code:    code,
		Data:    "",
		Message: msg,
	}
	c.JSON(http.StatusOK, rd)
	return
}

func FailedParam(c *gin.Context, err error) {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		rd := &Data{
			Code:    10005,
			Data:    "",
			Message: settings.RemoveTopStruct(validationErrs.Translate(global.IkubeopsTrans)),
		}
		c.JSON(http.StatusOK, rd)
		return
	} else {
		Failed(c, err.Error())
	}
}
