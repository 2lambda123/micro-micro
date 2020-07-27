// Package cliutil contains methods used across all cli commands
// @todo: get rid of os.Exits and use errors instread
package util

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/micro/v2/client/cli/namespace"
	clitoken "github.com/micro/micro/v2/client/cli/token"
	"github.com/micro/micro/v2/internal/config"
	"github.com/micro/micro/v2/internal/platform"
	muauth "github.com/micro/micro/v2/service/auth"
)

const (
	// EnvLocal is a builtin environment, it represents your local `micro server`
	EnvLocal = "local"
	// EnvPlatform is a builtin environment, the One True Micro Live(tm) environment.
	EnvPlatform = "platform"
)

const (
	// localProxyAddress is the default proxy address for environment server
	localProxyAddress = "127.0.0.1:8081"
	// platformProxyAddress is teh default proxy address for environment platform
	platformProxyAddress = "proxy.m3o.com"
)

var defaultEnvs = map[string]Env{
	EnvLocal: {
		Name:         EnvLocal,
		ProxyAddress: localProxyAddress,
	},
	EnvPlatform: {
		Name:         EnvPlatform,
		ProxyAddress: platformProxyAddress,
	},
}

func isBuiltinService(command string) bool {
	for _, service := range platform.Services {
		if command == service {
			return true
		}
	}
	return false
}

// SetProxyAddress includes things that should run for each command.
func SetProxyAddress(ctx *cli.Context) {
	// This makes `micro [command name] --help` work without a server
	for _, arg := range os.Args {
		if arg == "--help" || arg == "-h" {
			return
		}
	}
	switch ctx.Args().First() {
	case "new", "server", "help", "env":
		return
	}

	// fix for "micro service [command]", e.g "micro service auth"
	if ctx.Args().First() == "service" && isBuiltinService(ctx.Args().Get(1)) {
		return
	}

	// don't set the proxy address on the proxy
	if ctx.Args().First() == "proxy" {
		return
	}

	env := GetEnv(ctx)
	if len(env.ProxyAddress) == 0 {
		return
	}

	// Set the proxy. TODO: Pass this as an option to the client instead.
	setFlags(ctx, []string{"MICRO_PROXY=" + env.ProxyAddress})
}

type Env struct {
	Name         string
	ProxyAddress string
}

func AddEnv(env Env) {
	envs := getEnvs()
	envs[env.Name] = env
	setEnvs(envs)
}

func getEnvs() map[string]Env {
	envsJSON, err := config.Get("envs")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	envs := map[string]Env{}
	if len(envsJSON) > 0 {
		err := json.Unmarshal([]byte(envsJSON), &envs)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	for k, v := range defaultEnvs {
		envs[k] = v
	}
	return envs
}

func setEnvs(envs map[string]Env) {
	envsJSON, err := json.Marshal(envs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = config.Set(string(envsJSON), "envs")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// GetEnv returns the current selected environment
// Does not take
func GetEnv(ctx *cli.Context) Env {
	var envName string
	if len(ctx.String("env")) > 0 {
		envName = ctx.String("env")
	} else {
		env, err := config.Get("env")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if env == "" {
			env = EnvLocal
		}
		envName = env
	}

	return GetEnvByName(envName)
}

func GetEnvByName(env string) Env {
	envs := getEnvs()

	envir, ok := envs[env]
	if !ok {
		fmt.Println(fmt.Sprintf("Env \"%s\" not found. See `micro env` for available environments.", env))
		os.Exit(1)
	}

	if len(envir.ProxyAddress) == 0 {
		return envir
	}

	// default to :8081 (the proxy port)
	if _, port, _ := net.SplitHostPort(envir.ProxyAddress); len(port) == 0 {
		envir.ProxyAddress = net.JoinHostPort(envir.ProxyAddress, "8081")
	}

	return envir
}

func GetEnvs() []Env {
	envs := getEnvs()
	ret := []Env{defaultEnvs[EnvLocal], defaultEnvs[EnvPlatform]}
	nonDefaults := []Env{}
	for _, env := range envs {
		if _, isDefault := defaultEnvs[env.Name]; !isDefault {
			nonDefaults = append(nonDefaults, env)
		}
	}
	// @todo order nondefault envs alphabetically
	ret = append(ret, nonDefaults...)
	return ret
}

// SetEnv selects an environment to be used.
func SetEnv(envName string) {
	envs := getEnvs()
	_, ok := envs[envName]
	if !ok {
		fmt.Printf("Environment '%v' does not exist\n", envName)
		os.Exit(1)
	}
	config.Set(envName, "env")
}

// DelEnv deletes an env from config
func DelEnv(envName string) {
	envs := getEnvs()
	_, ok := envs[envName]
	if !ok {
		fmt.Printf("Environment '%v' does not exist\n", envName)
		os.Exit(1)
	}
	delete(envs, envName)
	setEnvs(envs)
}

func IsPlatform(ctx *cli.Context) bool {
	return GetEnv(ctx).Name == EnvPlatform
}

type Exec func(*cli.Context, []string) ([]byte, error)

func Print(e Exec) func(*cli.Context) error {
	return func(c *cli.Context) error {
		rsp, err := e(c, c.Args().Slice())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(rsp) > 0 {
			fmt.Printf("%s\n", string(rsp))
		}
		return nil
	}
}

func toFlag(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, "MICRO_", ""))
}

func setFlags(ctx *cli.Context, envars []string) {
	for _, envar := range envars {
		// setting both env and flags here
		// as the proxy settings for example did not take effect
		// with only flags
		parts := strings.Split(envar, "=")
		key := toFlag(parts[0])
		os.Setenv(parts[0], parts[1])
		ctx.Set(key, parts[1])
	}
}

// SetAuthToken handles exchanging refresh tokens to access tokens
// The structure of the local micro userconfig file is the following:
// micro.auth.[envName].token: temporary access token
// micro.auth.[envName].refresh-token: long lived refresh token
// micro.auth.[envName].expiry: expiration time of the access token, seconds since Unix epoch.
func SetAuthToken(ctx *cli.Context) error {
	env := GetEnv(ctx)
	ns, err := namespace.Get(env.Name)
	if err != nil {
		return err
	}

	tok, err := clitoken.Get(env.Name)
	if err != nil {
		return err
	}

	// If there is no refresh token, do not try to refresh it
	if len(tok.RefreshToken) == 0 {
		return nil
	}

	// Check if token is valid
	if time.Now().Before(tok.Expiry.Add(-15 * time.Second)) {
		muauth.DefaultAuth.Init(auth.ClientToken(tok))
		return nil
	}

	// Get new access token from refresh token if it's close to expiry
	tok, err = muauth.DefaultAuth.Token(
		auth.WithToken(tok.RefreshToken),
		auth.WithTokenIssuer(ns),
	)
	if err != nil {
		clitoken.Remove(env.Name)
		return nil
	}

	// Save the token to user config file
	muauth.DefaultAuth.Init(auth.ClientToken(tok))
	return clitoken.Save(env.Name, tok)
}
