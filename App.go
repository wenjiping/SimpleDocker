package main

import (
	"SimpleDocker/api"
	_ "SimpleDocker/context"
	_ "SimpleDocker/routers"
	"SimpleDocker/utils"
	"encoding/json"
	"errors"
	"flag"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"net/http"
	"strconv"
)

var port = flag.Int("port", 4050, "help message for flagname")
var success = []byte("SUPPORT OPTIONS")

// 跨域配置
var corsFunc = func(ctx *context.Context) {
	origin := ctx.Input.Header("Origin")
	ctx.Output.Header("Access-Control-Allow-Methods", "OPTIONS,DELETE,POST,GET,PUT,PATCH")
	ctx.Output.Header("Access-Control-Max-Age", "3600")
	ctx.Output.Header("Access-Control-Allow-Headers", "x-requested-with,X-Custom-Header,accept,Content-Type,Access-Token,authorization")
	ctx.Output.Header("Access-Control-Allow-Credentials", "true")
	ctx.Output.Header("Access-Control-Allow-Origin", origin)
	if ctx.Input.Method() == http.MethodOptions {
		// options请求，返回200
		ctx.Output.SetStatus(http.StatusOK)
		_ = ctx.Output.Body(success)
	} else {
		url := ctx.Input.URL()
		if url != "/api/system/login" {
			header := ctx.Input.Header("Authorization")
			if header == "" {
				ctx.Output.Status = 200
				respData := utils.PackageError(errors.New("登录过期"))
				marshal, _ := json.Marshal(respData)
				_ = ctx.Output.Body(marshal)
			}
		}
	}
}

func main() {
	flag.Parse()

	beego.BConfig.CopyRequestBody = true
	beego.BConfig.WebConfig.Session.SessionOn = true
	// 配置静态资源

	beego.SetStaticPath("/", "./static")

	// 配置路由
	beego.Include(&api.DockerController{})
	beego.Include(&api.ContainerController{})
	beego.Include(&api.ImageController{})
	beego.Include(&api.VolumeController{})
	beego.Include(&api.NetworkController{})
	beego.Include(&api.LoginController{})

	// 添加CORS
	beego.InsertFilter("/*", beego.BeforeRouter, corsFunc)

	// 启动服务
	beego.Run(":" + strconv.Itoa(*port))
}
