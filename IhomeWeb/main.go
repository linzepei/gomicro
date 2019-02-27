package main

import (
	"github.com/micro/go-log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-web"
	"gomicro/IhomeWeb/handler"
	_ "gomicro/IhomeWeb/models"
)

func main() {
	// create new web service
	//创建1个新的web服务
	service := web.NewService(
		web.Name("go.micro.web.IhomeWeb"),
		web.Version("latest"),
		web.Address(":8888"),
	)

	// initialise service
	//服务初始化
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	//使用路由中间件来映射页面
	rou := httprouter.New()

	rou.NotFound = http.FileServer(http.Dir("html"))
	//获取地区请求
	rou.GET("/api/v1.0/areas", handler.GetArea)
	//获取验证码图片
	rou.GET("/api/v1.0/imagecode/:uuid", handler.GetImageCd)
	//获取短信验证码
	rou.GET("/api/v1.0/smscode/:mobile", handler.GetSmscd)
	//获取短信验证码
	rou.POST("/api/v1.0/users", handler.PostRet)
	//获取session
	rou.GET("/api/v1.0/session", handler.GetSession)
	//登陆
	rou.POST("/api/v1.0/sessions", handler.PostLogin)
	//退出登陆
	rou.DELETE("/api/v1.0/session", handler.DeleteSession)
	//获取用户信息
	rou.GET("/api/v1.0/user", handler.GetUserInfo)
	//上传头像
	rou.POST("/api/v1.0/user/avatar", handler.PostAvatar)
	//用户认证检查
	rou.GET("/api/v1.0/user/auth", handler.GetUserAuth)
	//实名认证
	rou.POST("/api/v1.0/user/auth", handler.PostUserAuth)
	//获取用户房源信息
	rou.GET("/api/v1.0/user/houses", handler.GetUserHouses)
	//发布房源信息
	rou.POST("/api/v1.0/houses", handler.PostHouses)

	//获取首页轮播图
	rou.GET("/api/v1.0/house/index", handler.GetIndex)
	// register html handler
	//映射前端页面
	//service.Handle("/", http.FileServer(http.Dir("html")))
	service.Handle("/", rou)

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
