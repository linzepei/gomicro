package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"gomicro/PostHouses/handler"

	"github.com/micro/go-grpc"
	example "gomicro/PostHouses/proto/example"
)

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.PostHouses"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	example.RegisterExampleHandler(service.Server(), new(handler.Example))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
