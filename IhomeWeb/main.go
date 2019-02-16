package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-log"
	"net/http"

	"github.com/micro/go-web"
	"gomicro/IhomeWeb/handler"
	_ "gomicro/IhomeWeb/models"
)

func main() {
	// create new web service
	service := web.NewService(
		web.Name("go.micro.web.IhomeWeb"),
		web.Version("latest"),
		web.Address(":8888"),
	)

	//使用路由中间件来映射页面
	rou := httprouter.New()
	rou.NotFound = http.FileServer(http.Dir("html"))

	//获取地区请求
	rou.GET("/api/v1.0/areas", handler.GetArea)

	// initialise service
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	// register html handler
	//service.Handle("/", http.FileServer(http.Dir("html")))
	service.Handle("/", rou)

	// register call handler
	service.HandleFunc("/example/call", handler.ExampleCall)

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
