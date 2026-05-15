package controler

import (
	"os"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func MiddlewareAuth(r *ghttp.Request) {
	token := os.Getenv("AUTH_TOKEN")
	if token == "" {
		token = g.Cfg().MustGet(r.Context(), "server.authToken").String()
	}
	if token == "" {
		r.Middleware.Next()
		return
	}

	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		if strings.TrimPrefix(authHeader, "Bearer ") == token {
			r.Middleware.Next()
			return
		}
	}

	r.Response.WriteStatus(401)
	r.Response.WriteJson(g.Map{
		"success": false,
		"message": "unauthorized",
	})
}
