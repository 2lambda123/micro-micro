package main

//go:generate ./scripts/generate.sh

import (
	"fmt"
	"os"
	"unicode"

	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/micro/v3/cmd"

	// internal packages
	_ "github.com/micro/micro/v3/internal/usage"

	// load packages so they can register commands
	_ "github.com/micro/micro/v3/client/cli"
	_ "github.com/micro/micro/v3/client/cli/init"
	_ "github.com/micro/micro/v3/client/cli/new"
	_ "github.com/micro/micro/v3/client/cli/signup"
	_ "github.com/micro/micro/v3/client/cli/user"
	_ "github.com/micro/micro/v3/server"
	_ "github.com/micro/micro/v3/service/auth/cli"
	_ "github.com/micro/micro/v3/service/cli"
	_ "github.com/micro/micro/v3/service/config/cli"
	_ "github.com/micro/micro/v3/service/network/cli"
	_ "github.com/micro/micro/v3/service/runtime/cli"
	_ "github.com/micro/micro/v3/service/store/cli"
)

func main() {
	if err := cmd.DefaultCmd.Run(); err != nil {
		fmt.Println(formatErr(err))
		os.Exit(1)
	}
}

func formatErr(err error) string {
	switch v := err.(type) {
	case *errors.Error:
		return upcaseInitial(v.Detail)
	default:
		return upcaseInitial(err.Error())
	}
}

func upcaseInitial(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}
