package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func DemoMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		fmt.Println("进入 demo 中间件")
		c.Next()
		fmt.Println("离开 demo 中间件")
	}
}
