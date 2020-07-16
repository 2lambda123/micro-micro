package store

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	pb "github.com/micro/go-micro/v2/store/service/proto"
	mcli "github.com/micro/micro/v2/client/cli"
	"github.com/micro/micro/v2/internal/helper"
	"github.com/micro/micro/v2/service/store/handler"
)

var (
	// Name of the store service
	Name = "go.micro.store"
	// Address is the store address
	Address = ":8002"
)

// Run runs the micro server
func Run(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Init(log.WithFields(map[string]interface{}{"service": "store"}))

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}

	// Initialise service
	service := micro.NewService(
		micro.Name(Name),
		micro.Address(Address),
	)

	// the store handler
	h := handler.New(service.Options().Store)
	pb.RegisterStoreHandler(service.Server(), h)

	// start the service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

// Commands is the cli interface for the store service
func Commands(options ...micro.Option) []*cli.Command {
	command := &cli.Command{
		Name:  "store",
		Usage: "Run the micro store service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "address",
				Usage:   "Set the micro tunnel address :8002",
				EnvVars: []string{"MICRO_SERVER_ADDRESS"},
			},
		},
		Action: func(ctx *cli.Context) error {
			if err := helper.UnexpectedSubcommand(ctx); err != nil {
				return err
			}
			Run(ctx, options...)
			return nil
		},
		Subcommands: mcli.StoreCommands(),
	}

	for _, p := range Plugins() {
		if cmds := p.Commands(); len(cmds) > 0 {
			command.Subcommands = append(command.Subcommands, cmds...)
		}

		if flags := p.Flags(); len(flags) > 0 {
			command.Flags = append(command.Flags, flags...)
		}
	}

	return []*cli.Command{command}
}
