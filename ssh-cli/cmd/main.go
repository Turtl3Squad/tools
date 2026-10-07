package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lmittmann/tint"

	"github.com/turtl3squad/tools/console/pkg/cli"
)

var version string

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.Kitchen,

		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.Attr{}
			}
			return a
		},
	})))

	app := initApp()

	if err := app.Run(ctx, os.Args); err != nil {
		cli.Fatal(err)
	}
}

func initApp() cli.Command {
	app := cli.Command{
		Name:  config.Name(),
		Usage: "Platform Plane CLI",

		Suggest: true,
		Version: version,

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "space",
				Usage: "The name of the space to use",
			},

			&cli.StringFlag{
				Name:  "cluster",
				Usage: "The name of the cluster to use",
			},

			&cli.StringFlag{
				Name:  "registry",
				Usage: "The name of the registry to use",
			},
		},

		Commands: []*cli.Command{
			login.Command,
			logout.Command,
			token.Command,

			// Platform
			cluster.Command,
			hypervisor.Command,

			app.Command,
			api.Command,

			ai.Command,
			openai.Command,
			companion.Command,

			// Artifacts
			git.Command,
			docker.Command,
		},
	}

	app.Commands = append(app.Commands, catalog.Commands...)

	return app
}
