package main

import (
	"faststartup_go/src/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter(listen string) {

	routes := gin.Default()

	// 添加中间件
	routes.Use(utils.CORS())

	gin.SetMode(gin.ReleaseMode)

	routes.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, "hello, .. (｡◕ˇ∀ˇ◕）)")
	})

	//err := r.RunTLS(listen, s.cfg.Sec.CertPath, s.cfg.Sec.Key)
	err := routes.Run(listen)
	if err != nil {
		panic(err)
	}

}


