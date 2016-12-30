# Micro Bot

The **micro bot** is a bot that sits inside your microservices environment which you can interact with via Slack, HipChat, XMPP, etc. 
It mimics the functions of the CLI via messaging.

<p align="center">
  <img src="https://raw.githubusercontent.com/micro/micro/master/doc/images/bot.png" />
</p>

## Supported Inputs

- Slack
- HipChat
- Discord

## Getting Started

### Install Micro

```go
go get github.com/micro/micro
```

### Run with Slack

```shell
micro bot --inputs=slack --slack_token=SLACK_TOKEN
```

<img src="https://raw.githubusercontent.com/micro/micro/master/doc/images/slack.png">
-

### Run with HipChat

```shell
micro bot --inputs=hipchat --hipchat_username=XMPP_USER --hipchat_password=XMPP_PASSWORD
```

<img src="https://raw.githubusercontent.com/micro/micro/master/doc/images/hipchat.png">
-

### Run with Discord

Basic: 
```shell
micro bot --inputs=discord --discord_token="Bot TOKEN"
```

With Prefix: (it also accepts mentions as prefix)
```shell
micro bot --inputs=discord --discord_token="Bot TOKEN" --discord_prefix="micro "
```

With Prefix & Whitelist: (only whitelisted users will be able to run any commands, seperated by `,`)
```shell
micro bot --inputs=discord --discord_token="Bot TOKEN" --discord_prefix="micro " --discord_whitelist=USERID,ANOTHERUSERID
```

<img src="https://raw.githubusercontent.com/micro/micro/master/doc/images/discord.png">
-

Use multiple inputs by specifying a comma separated list

```shell
micro bot --inputs=hipchat,slack --slack_token=SLACK_TOKEN --hipchat_username=XMPP_USER --hipchat_password=XMPP_PASSWORD
```

### Help

In slack
```shell
micro help

deregister service [definition] - Deregisters a service
echo [text] - Returns the [text]
get service [name] - Returns a registered service
health [service] - Returns health of a service
hello - Returns a greeting
list services - Returns a list of registered services
ping - Returns pong
query [service] [method] [request] - Returns the response for a service query
register service [definition] - Registers a service
the three laws - Returns the three laws of robotics
time - Returns the server time
```

## Adding new Commands

Commands are functions executed by the bot based on text based pattern matching.

### Write a Command

```go
import "github.com/micro/micro/bot/command"

func Ping() command.Command {
	usage := "ping"
	description := "Returns pong"

	return command.NewCommand("ping", usage, desc, func(args ...string) ([]byte, error) {
		return []byte("pong"), nil
	})
}
```

### Register the command

Add the command to the Commands map with a pattern key that can be matched by golang/regexp.Match

```go
import "github.com/micro/micro/bot/command"

func init() {
	command.Commands["^ping$"] = Ping()
}
```

### Link the Command

Drop a link to your command into the top level dir

link_command.go:
```go
import _ "path/to/import"
```

## Adding new Inputs

Inputs are plugins for communication e.g Slack, HipChat, XMPP, IRC, SMTP, etc, etc. 

New inputs can be added in the following way.

### Write an Input

Write an input that satisfies the Input interface.

```go
type Input interface {
	// Provide cli flags
	Flags() []cli.Flag
	// Initialise input using cli context
	Init(*cli.Context) error
	// Stream events from the input
	Stream() (Conn, error)
	// Start the input
	Start() error
	// Stop the input
	Stop() error
	// name of the input
	String() string
}
```

### Register the input

Add the input to the Inputs map.

```go
import "github.com/micro/micro/bot/input"

func init() {
	input.Inputs["name"] = MyInput
}
```

### Link the input

Drop a link to your input into the top level dir

link_input.go:
```go
import _ "path/to/import"
```
