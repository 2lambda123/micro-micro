// Package runtime is the micro runtime
package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	golog "github.com/micro/go-micro/v3/logger"
	"github.com/micro/go-micro/v3/runtime/local/source/git"
	"github.com/micro/go-micro/v3/util/file"
	"github.com/micro/micro/v3/client/cli/namespace"
	"github.com/micro/micro/v3/client/cli/util"
	cliutil "github.com/micro/micro/v3/client/cli/util"
	muclient "github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime"
	"github.com/micro/micro/v3/service/runtime/server"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/status"
)

const (
	// RunUsage message for the run command
	RunUsage = "Run a service: micro run [source]"
	// KillUsage message for the kill command
	KillUsage = "Kill a service: micro kill [source]"
	// UpdateUsage message for the update command
	UpdateUsage = "Update a service: micro update [source]"
	// GetUsage message for micro get command
	GetUsage = "Get the status of services"
	// ServicesUsage message for micro services command
	ServicesUsage = "micro services"
	// LogUsage message for logs command
	LogUsage = "Required usage: micro log example"
	// CannotWatch message for the run command
	CannotWatch = "Cannot watch filesystem on this runtime"
	// credentialsKey is the key for the secret in which git credentials should be passed when
	// creating or updating a service
	credentialsKey = "GIT_CREDENTIALS"
)

var (
	// DefaultRetries which should be attempted when starting a service
	DefaultRetries = 3
	// DefaultImage which should be run
	DefaultImage = "micro/cells:micro"
	// GitOrgs we currently support for credentials
	GitOrgs = []string{"github", "bitbucket", "gitlab"}
)

func runService(ctx *cli.Context) error {
	// we need some args to run
	if ctx.Args().Len() == 0 {
		fmt.Println(RunUsage)
		return nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	source, err := git.ParseSourceLocal(wd, appendSourceBase(ctx, wd, ctx.Args().Get(0)))
	if err != nil {
		return err
	}
	var newSource string
	if source.Local {
		if cliutil.IsPlatform(ctx) {
			fmt.Println("Local sources are not yet supported on m3o. It's coming soon though!")
			os.Exit(1)
		}
		newSource, err = upload(ctx, source)
		if err != nil {
			return err
		}
	} else {
		err := sourceExists(source)
		if err != nil {
			return err
		}
	}

	typ := ctx.String("type")
	command := strings.TrimSpace(ctx.String("command"))
	args := strings.TrimSpace(ctx.String("args"))

	runtimeSource := source.RuntimeSource()
	if source.Local {
		runtimeSource = newSource
	}

	var retries = DefaultRetries
	if ctx.IsSet("retries") {
		retries = ctx.Int("retries")
	}

	var image = DefaultImage
	if ctx.IsSet("image") {
		image = ctx.String("image")
	}

	// when using the micro/cells:go image, we pass the source as the argument
	args = runtimeSource
	if len(source.Ref) > 0 {
		args += "@" + source.Ref
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
	for _, evar := range ctx.StringSlice("env_vars") {
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

	// determine the namespace
	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return err
	}
	opts = append(opts, runtime.CreateNamespace(ns))
	gitCreds, ok := getGitCredentials(source.Repo)
	if ok {
		opts = append(opts, runtime.WithSecret(credentialsKey, gitCreds))
	}

	// run the service
	service := &runtime.Service{
		Name:     source.RuntimeName(),
		Source:   runtimeSource,
		Version:  source.Ref,
		Metadata: make(map[string]string),
	}

	if err := runtime.Create(service, opts...); err != nil {
		return err
	}

	if runtime.DefaultRuntime.String() == "local" {
		// we need to wait
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		<-ch
		// delete the service
		return runtime.Delete(service)
	}

	return nil
}

func killService(ctx *cli.Context) error {
	// we need some args to run
	if ctx.Args().Len() == 0 {
		fmt.Println(KillUsage)
		return nil
	}

	name := ctx.Args().Get(0)
	ref := ""
	if parts := strings.Split(name, "@"); len(parts) > 1 {
		name = parts[0]
		ref = parts[1]
	}
	if ref == "" {
		ref = "latest"
	}
	service := &runtime.Service{
		Name:    name,
		Version: ref,
	}

	// determine the namespace
	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return err
	}

	if err := runtime.Delete(service, runtime.DeleteNamespace(ns)); err != nil {
		return err
	}

	return nil
}

func upload(ctx *cli.Context, source *git.Source) (string, error) {
	if err := grepMain(source.FullPath); err != nil {
		return "", err
	}
	uploadedFileName := filepath.Base(source.Folder) + ".tar.gz"
	path := filepath.Join(os.TempDir(), uploadedFileName)

	var err error
	if len(source.LocalRepoRoot) > 0 {
		// @todo currently this uploads the whole repo all the time to support local dependencies
		// in parents (ie service path is `repo/a/b/c` and it depends on `repo/a/b`).
		// Optimise this by only uploading things that are needed.
		err = server.Compress(source.LocalRepoRoot, path)
	} else {
		err = server.Compress(source.FullPath, path)
	}

	if err != nil {
		return "", err
	}
	cli := muclient.DefaultClient
	err = file.New("server", cli, file.WithContext(context.DefaultContext)).Upload(uploadedFileName, path)
	if err != nil {
		return "", err
	}
	// ie. if relative folder path to repo root is `test/service/example`
	// file name becomes `example.tar.gz/test/service`
	parts := strings.Split(source.Folder, "/")
	if len(parts) == 1 {
		return uploadedFileName, nil
	}
	allButLastDir := parts[0 : len(parts)-1]
	return filepath.Join(append([]string{uploadedFileName}, allButLastDir...)...), nil
}

func updateService(ctx *cli.Context) error {
	// we need some args to run
	if ctx.Args().Len() == 0 {
		fmt.Println(RunUsage)
		return nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	source, err := git.ParseSourceLocal(wd, appendSourceBase(ctx, wd, ctx.Args().Get(0)))
	if err != nil {
		return err
	}
	var newSource string
	if source.Local {
		newSource, err = upload(ctx, source)
		if err != nil {
			return err
		}
	}

	runtimeName := source.RuntimeName()
	runtimeSource := source.RuntimeSource()
	ref := source.Ref
	if source.Local {
		runtimeSource = newSource
	} else {
		runtimeSource = ""
		name := ctx.Args().Get(0)
		if parts := strings.Split(name, "@"); len(parts) > 1 {
			runtimeName = parts[0]
			ref = parts[1]
		}
	}
	if ref == "" {
		ref = "latest"
	}
	service := &runtime.Service{
		Name:    runtimeName,
		Source:  runtimeSource,
		Version: ref,
	}

	// determine the namespace
	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return err
	}

	opts := []runtime.UpdateOption{runtime.UpdateNamespace(ns)}
	gitCreds, ok := getGitCredentials(source.Repo)
	if ok {
		opts = append(opts, runtime.UpdateSecret(credentialsKey, gitCreds))
	}
	return runtime.Update(service, runtime.UpdateNamespace(ns))
}

func getService(ctx *cli.Context) error {
	name := ""
	version := "latest"
	typ := ctx.String("type")

	if ctx.Args().Len() > 0 {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		source, err := git.ParseSourceLocal(wd, ctx.Args().Get(0))
		if err != nil {
			return err
		}
		name = source.RuntimeName()
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
			return nil
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

	// determine the namespace
	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return err
	}
	readOpts = append(readOpts, runtime.ReadNamespace(ns))

	// read the service
	services, err = runtime.Read(readOpts...)
	if err != nil {
		return err
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
		return nil
	}

	sort.Slice(services, func(i, j int) bool { return services[i].Name < services[j].Name })

	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintln(writer, "NAME\tVERSION\tSOURCE\tSTATUS\tBUILD\tUPDATED\tMETADATA")
	for _, service := range services {
		// cut the commit down to first 7 characters
		build := parse(service.Metadata["build"])
		if len(build) > 7 {
			build = build[:7]
		}

		// if there is an error, display this in metadata (there is no error field)
		metadata := fmt.Sprintf("owner=%s, group=%s", parse(service.Metadata["owner"]), parse(service.Metadata["group"]))
		if service.Status == runtime.Error {
			metadata = fmt.Sprintf("%v, error=%v", metadata, parse(service.Metadata["error"]))
		}

		// parse when the service was started
		updated := parse(timeAgo(service.Metadata["started"]))

		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			service.Name,
			parse(service.Version),
			parse(service.Source),
			humanizeStatus(service.Status),
			build,
			updated,
			metadata)
	}
	writer.Flush()
	return nil
}

const ()

func getLogs(ctx *cli.Context) error {
	logger.DefaultLogger.Init(golog.WithFields(map[string]interface{}{"service": "runtime"}))
	if ctx.Args().Len() == 0 {
		fmt.Println("Service name is required")
		return nil
	}

	name := ctx.Args().Get(0)

	fmt.Println(LogUsage)
	return nil
	// must specify service name
	if len(name) == 0 {
	}

	// get the args
	options := []runtime.LogsOption{}

	count := ctx.Int("lines")
	if count > 0 {
		options = append(options, runtime.LogsCount(int64(count)))
	} else {
		options = append(options, runtime.LogsCount(int64(15)))
	}

	follow := ctx.Bool("follow")

	if follow {
		options = append(options, runtime.LogsStream(follow))
	}

	// @todo reintroduce since
	//since := ctx.String("since")
	//var readSince time.Time
	//d, err := time.ParseDuration(since)
	//if err == nil {
	//	readSince = time.Now().Add(-d)
	//}

	// determine the namespace
	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return err
	}
	options = append(options, runtime.LogsNamespace(ns))

	logs, err := runtime.Log(&runtime.Service{Name: name}, options...)

	if err != nil {
		return err
	}

	output := ctx.String("output")
	for {
		select {
		case record, ok := <-logs.Chan():
			if !ok {
				if err := logs.Error(); err != nil {
					fmt.Printf("Error reading logs: %s\n", status.Convert(err).Message())
					os.Exit(1)
				}
				return nil
			}
			switch output {
			case "json":
				b, _ := json.Marshal(record)
				fmt.Printf("%v\n", string(b))
			default:
				fmt.Printf("%v\n", record.Message)

			}
		}
	}
}

func humanizeStatus(status runtime.ServiceStatus) string {
	switch status {
	case runtime.Pending:
		return "pending"
	case runtime.Building:
		return "building"
	case runtime.Starting:
		return "starting"
	case runtime.Running:
		return "running"
	case runtime.Stopping:
		return "stopping"
	case runtime.Stopped:
		return "stopped"
	case runtime.Error:
		return "error"
	default:
		return "unknown"
	}
}
