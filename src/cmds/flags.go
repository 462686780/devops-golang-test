package cmds

import "github.com/urfave/cli"

var (
	HelpFlag = cli.BoolFlag{Name: "help,h", Usage: "Show help message", Hidden: true}
)
