package controler

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type Infos struct {
	Name       string `json:"name"`
	Systemtype string `json:"system_type"`
	Note       string `json:"note"`
}

func Router() {
	s := g.Server()
	s.Use(MiddlewareCORS)

	s.Group("/fir", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareAuth)

		group.GET("/info", GetInfo)
		group.POST("/delete", DeleteInfo)
		group.POST("/upload", Upload)
		group.GET("/page", GetPage)
		group.POST("/update", UpdateInfo)
	})

	s.Run()
}
