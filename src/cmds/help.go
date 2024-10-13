package cmds

import (
	"statefulset/base"

	"github.com/urfave/cli"
)

var CmdHelp = cli.Command{
	Name:        "help",
	Usage:       "Show a list of comands or help for one command",
	Description: "Show a list of comands or help for one command",
	ArgsUsage:   "[command]",
	Action: func(c *cli.Context) error {
		args := c.Args()
		if args.Present() {
			return cli.ShowCommandHelp(c, args.First())
		}
		return cli.ShowAppHelp(c)
	},
	OnUsageError: base.OnUsageError,
}
