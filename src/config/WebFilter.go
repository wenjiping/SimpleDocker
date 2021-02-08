package config

import (
	"SimpleDocker/src/auth"
	"SimpleDocker/src/utils"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/context"
	"net/http"
	"strings"
)

// 跨域配置
func Handler(ctx *context.Context) {
	origin := ctx.Input.Header("Origin")
	ctx.Output.Header("Access-Control-Allow-Methods", "OPTIONS,DELETE,POST,GET,PUT,PATCH")
	ctx.Output.Header("Access-Control-Max-Age", "3600")
	ctx.Output.Header("Access-Control-Allow-Headers", "x-requested-with,X-Custom-Header,accept,Content-Type,Access-Token,authorization,responsetype")
	ctx.Output.Header("Access-Control-Allow-Credentials", "true")
	ctx.Output.Header("Access-Control-Allow-Origin", origin)
	CheckAuth(ctx)
}

func CheckAuth(ctx *context.Context) {
	// options请求，返回200
	if ctx.Input.Method() == http.MethodOptions {
		ctx.Output.SetStatus(http.StatusOK)
		_ = ctx.Output.Body([]byte("SUPPORT OPTIONS"))
	} else {
		url := ctx.Input.URL()
		if url != "/api/system/login" && !strings.HasPrefix(url, "/ws") {
			header := ctx.Input.Header("Authorization")
			err := auth.ParseToken(header)
			if header == "" || err != nil {
				ctx.Output.Status = 403
				respData := utils.PackageError(errors.New("无效的Token"))
				marshal, _ := json.Marshal(respData)
				_ = ctx.Output.Body(marshal)
			}
		} else if strings.HasPrefix(url, "/ws") {
			header := ctx.Input.Query("token")
			err := auth.ParseToken(header)
			if header == "" || err != nil {
				ctx.Output.Status = 403
				respData := utils.PackageError(errors.New("无效的Token"))
				marshal, _ := json.Marshal(respData)
				_ = ctx.Output.Body(marshal)
			}
		}
	}
}
