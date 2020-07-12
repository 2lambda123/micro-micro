package main

import (
	"fmt"

	"github.com/micro/go-micro/v2"
)

func main() {
	// New Service
	service := service.NewService(
		service.Name("go.micro.service.config-read"),
		service.Version("latest"),
	)
	service.Init()

	// create a new config
	c := service.Options().Config

	// set a value
	fmt.Println("Value of key.subkey: ", c.Get("key", "subkey").String(""))
}
