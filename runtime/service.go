// Package runtime is the micro runtime
package runtime

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/runtime"
	srvRuntime "github.com/micro/go-micro/v2/runtime/service"
	cliutil "github.com/micro/micro/v2/cli/util"
	"github.com/micro/micro/v2/internal/git"
)

const (
	// RunUsage message for the run command
	RunUsage = "Required usage: micro run [source]"
	// KillUsage message for the kill command
	KillUsage = "Require usage: micro kill [source]"
	// UpdateUsage message for the update command
	UpdateUsage = "Require usage: micro update [source]"
	// GetUsage message for micro get command
	GetUsage = "Require usage: micro ps [service] [version]"
	// ServicesUsage message for micro services command
	ServicesUsage = "Require usage: micro services"
	// CannotWatch message for the run command
	CannotWatch = "Cannot watch filesystem on this runtime"
)

var (
	// DefaultRetries which should be attempted when starting a service
	DefaultRetries = 3
	// Image to specify if none is specified
	Image = "docker.pkg.github.com/micro/services"
	// Source where we get services from
	Source = "github.com/micro/services"
)

// timeAgo returns the time passed
func timeAgo(v string) string {
	if len(v) == 0 {
		return "unknown"
	}
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return v
	}
	return fmt.Sprintf("%v ago", time.Since(t).Truncate(time.Second))
}

func runtimeFromContext(ctx *cli.Context) runtime.Runtime {
	if cliutil.IsLocal() {
		return *cmd.DefaultCmd.Options().Runtime
	}
	return srvRuntime.NewRuntime()
}

// exists returns whether the given file or directory exists
func dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func runService(ctx *cli.Context, srvOpts ...micro.Option) {
	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	// we need some args to run
	if ctx.Args().Len() == 0 {
		fmt.Println(RunUsage)
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	source, err := git.ParseSourceLocal(wd, ctx.Args().Get(0))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	typ := ctx.String("type")
	image := ctx.String("image")
	command := strings.TrimSpace(ctx.String("command"))
	args := strings.TrimSpace(ctx.String("args"))

	// load the runtime
	r := runtimeFromContext(ctx)

	var retries = DefaultRetries
	if ctx.IsSet("retries") {
		retries = ctx.Int("retries")
	}

	if cliutil.IsPlatform() && len(image) == 0 {
		if source.Local {
			fmt.Println("Can't run local code on platform")
			os.Exit(1)
		}

		formattedName := strings.ReplaceAll(source.Folder, "/", "-")
		// eg. docker.pkg.github.com/micro/services/users-api
		image = fmt.Sprintf("%v/%v", Image, formattedName)
	}

	// specify the options
	opts := []runtime.CreateOption{
		runtime.WithOutput(os.Stdout),
		runtime.WithRetries(retries),
		runtime.CreateImage(image),
		runtime.CreateType(typ),
	}

	// add environment variable passed in via cli
	var environment []string
	for _, evar := range ctx.StringSlice("env") {
		for _, e := range strings.Split(evar, ",") {
			if len(e) > 0 {
				environment = append(environment, strings.TrimSpace(e))
			}
		}
	}

	if len(environment) > 0 {
		opts = append(opts, runtime.WithEnv(environment))
	}

	if len(command) > 0 {
		opts = append(opts, runtime.WithCommand(strings.Split(command, " ")...))
	}

	if len(args) > 0 {
		opts = append(opts, runtime.WithArgs(strings.Split(args, " ")...))
	}

	// run the service
	service := &runtime.Service{
		Name:     source.RuntimeName(),
		Source:   source.RuntimeSource(),
		Version:  source.Ref,
		Metadata: make(map[string]string),
	}

	if err := r.Create(service, opts...); err != nil {
		fmt.Println(err)
		return
	}

	if r.String() == "local" {
		// we need to wait
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		<-ch
		// delete the service
		r.Delete(service)
	}
}

func killService(ctx *cli.Context, srvOpts ...micro.Option) {
	// we need some args to run
	if ctx.Args().Len() == 0 {
		fmt.Println(RunUsage)
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	source, err := git.ParseSourceLocal(wd, ctx.Args().Get(0))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	service := &runtime.Service{
		Name:    source.RuntimeName(),
		Source:  source.RuntimeSource(),
		Version: source.Ref,
	}

	if err := runtimeFromContext(ctx).Delete(service); err != nil {
		fmt.Println(err)
		return
	}
}

func updateService(ctx *cli.Context, srvOpts ...micro.Option) {
	// we need some args to run
	if ctx.Args().Len() == 0 {
		fmt.Println(RunUsage)
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	source, err := git.ParseSourceLocal(wd, ctx.Args().Get(0))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	service := &runtime.Service{
		Name:    source.RuntimeName(),
		Source:  source.RuntimeSource(),
		Version: source.Ref,
	}

	if err := runtimeFromContext(ctx).Update(service); err != nil {
		fmt.Println(err)
		return
	}
}

func getService(ctx *cli.Context, srvOpts ...micro.Option) {
	name := ctx.Args().Get(0)
	version := "latest"
	typ := ctx.String("type")
	r := runtimeFromContext(ctx)

	if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "/") {
		fmt.Println(GetUsage)
		return
	}

	// set version as second arg
	if ctx.Args().Len() > 1 {
		version = ctx.Args().Get(1)
	}

	// should we list sevices
	var list bool

	// zero args so list all
	if ctx.Args().Len() == 0 {
		list = true
	}

	var services []*runtime.Service
	var readOpts []runtime.ReadOption

	// return a list of services
	switch list {
	case true:
		// return specific type listing
		if len(typ) > 0 {
			readOpts = append(readOpts, runtime.ReadType(typ))
		}
	// return one service
	default:
		// check if service name was passed in
		if len(name) == 0 {
			fmt.Println(GetUsage)
			return
		}

		// get service with name and version
		readOpts = []runtime.ReadOption{
			runtime.ReadService(name),
			runtime.ReadVersion(version),
		}

		// return the runtime services
		if len(typ) > 0 {
			readOpts = append(readOpts, runtime.ReadType(typ))
		}

	}

	// read the service
	services, err := r.Read(readOpts...)
	if err != nil {
		fmt.Println(err)
		return
	}

	// make sure we return UNKNOWN when empty string is supplied
	parse := func(m string) string {
		if len(m) == 0 {
			return "n/a"
		}
		return m
	}

	// don't do anything if there's no services
	if len(services) == 0 {
		return
	}

	sort.Slice(services, func(i, j int) bool { return services[i].Name < services[j].Name })

	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintln(writer, "NAME\tVERSION\tSOURCE\tSTATUS\tBUILD\tUPDATED\tMETADATA")
	for _, service := range services {
		status := parse(service.Metadata["status"])
		if status == "error" {
			status = service.Metadata["error"]
		}

		// cut the commit down to first 7 characters
		build := parse(service.Metadata["build"])
		if len(build) > 7 {
			build = build[:7]
		}

		// parse when the service was started
		updated := parse(timeAgo(service.Metadata["started"]))

		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			service.Name,
			parse(service.Version),
			parse(service.Source),
			strings.ToLower(status),
			build,
			updated,
			fmt.Sprintf("owner=%s,group=%s", parse(service.Metadata["owner"]), parse(service.Metadata["group"])))
	}
	writer.Flush()
}
