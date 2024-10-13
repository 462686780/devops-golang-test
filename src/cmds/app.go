package cmds

import (
	"fmt"
	"os"
	"statefulset/base"

	"github.com/urfave/cli"
)

func App() *cli.App {
	app := cli.NewApp()
	app.Name = base.App
	app.Usage = base.Usage
	app.Description = base.Description
	app.Version = base.Version

	app.Commands = []cli.Command{
		CmdHelp,
		CmdVersion,
		CmdServer,
	}

	app.Before = Before
	app.After = After
	app.Flags = []cli.Flag{
		HelpFlag,
	}
	return app
}

// Before execute before any subcommands
func Before(ctx *cli.Context) error {
	return nil
}

// After execute after any subcommands
func After(ctx *cli.Context) error {
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	// log.Shutdown()
	return nil
}
