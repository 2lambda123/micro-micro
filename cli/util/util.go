// Package cliutil contains methods used across
// all cli commands
// @todo: get rid of os.Exits and use errors instread
package cliutil

import (
	"fmt"
	"os"

	"github.com/micro/go-micro/v2/util/config"
)

const (
	localAddress = "127.0.0.1:8081"
	liveAddress  = "proxy.micro.mu:443"
)

// SetupCommand includes things that should run for each command.
func SetupCommand() {
	os.Setenv("MICRO_PROXY", "service")
	env, err := config.Get("env")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(env) == 0 {
		os.Setenv("MICRO_PROXY_ADDRESS", localAddress)
		return
	}
	os.Setenv("MICRO_PROXY_ADDRESS", env)
}
