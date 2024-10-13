package cmds

import (
	"fmt"
	"runtime"
	"statefulset/base"

	"github.com/urfave/cli"
)

var CmdVersion = cli.Command{
	Name:         "version",
	Description:  "Print version information and quit",
	Action:       runVersion,
	OnUsageError: base.OnUsageError,
}

func runVersion(ctx *cli.Context) error {
	fmt.Printf("Version: 	%s\n", base.Version)
	fmt.Printf("Go: 		%s\n", runtime.Version())
	fmt.Printf("Compiled: 	%s\n", ctx.App.Compiled.Format("2006-01-02 15:04:05"))
	return nil
}
