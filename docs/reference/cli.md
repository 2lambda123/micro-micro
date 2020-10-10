---
title: CLI usage
keywords: micro
tags: [micro]
sidebar: home_sidebar
permalink: /reference/cli
summary: A CLI usage guide
---

## CLI

Micro is driven entirely through a CLI experience. This reference highlights the CLI design.

## Overview

The CLI speaks to the `micro server` through the gRPC proxy running locally by default on :8081. All requests are proxied based on your environment 
configuration. The CLI provides the sole interaction for controlling services and environments.

## Builtin Commands

Built in commands are system or configuration level commands for interacting with the server or 
changing user config. For the most part this is syntactic sugar for user convenience. Here's a 
subset of well known commands.

```
signup
login
run
update
kill
services
logs
status
env
user
```

The micro binary and each subcommand has a --help flag to provide a usage guide. The majority should be 
obvious to the user. We will go through a few in more detail.

### Signup

Signup is a command which attempts to query a "signup" to register a new account, this is env specific and requires a signup service to be 
running. By default locally this will not exist and we expect the user to use the admin/micro credentials to administrate the system. 
You can then choose to run your own signup service conforming to the proto in micro/proto or use `micro auth create account`. 

Signup is seen as a command for those who want to run their own micro server for others and potentially license the software to take payment.

### Login

Login authenticates the user and stores credentials locally in a .micro/tokens file. This calls the micro auth service to authenticate the 
user against existing accounts stored in the system. Login asks for a username and password at the prompt.

## Dynamic Commands

When issuing a command to the Micro CLI (ie. `micro command`), if the command is not a builtin, Micro will try to dynamically resolve this command and call
a service running. Let's take the `micro registry` command, because although the registry is a core service that's running by default on a local Micro setup,
the `registry` command is not a builtin one.

With the `--help` flag, we can get information about available subcommands and flags

```sh
$ micro registry --help
NAME:
	micro registry

VERSION:
	latest

USAGE:
	micro registry [command]

COMMANDS:
	deregister
	getService
	listServices
	register
	watch
```

The commands listed are endpoints of the `registry` service (see `micro services`).

To see the flags (which are essentially endpoint request parameters) for a subcommand:

```sh
$ micro registry getService --help
NAME:
	micro registry getService

USAGE:
	micro registry getService [flags]

FLAGS:
	--service string
	--options_ttl int64
	--options_domain string

```

At this point it is useful to have a look at the proto of the [registry service here](https://github.com/micro/micro/blob/master/proto/registry/registry.proto).

In particular, let's see the `GetService` endpoint definition to understand how request parameters map to flags:

```proto
message Options {
	int64 ttl = 1;
	string domain = 2;
}

message GetRequest {
	string service = 1;
	Options options = 2;
}
```

As the above definition tells us, the request of `GetService` has the field `service` at the top level, and fields `ttl` and `domain` in an options structure.
The dynamic CLI maps the underscored flagnames (ie. `options_domain`) to request fields, so the following request JSON:

```js
{
    "service": "serviceName",
    "options": {
        "domain": "domainExample"
    }
}
```

is equivalent to the following flags:

```sh
micro registry getService --service=serviceName --options_domain=domainExample
```
