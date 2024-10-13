package base

import (
	"fmt"

	"github.com/urfave/cli"
)

func OnUsageError(ctx *cli.Context, err error, isSubcomand bool) error {
	fmt.Printf("%v\n", err)
	if isSubcomand {
		return cli.ShowSubcommandHelp(ctx)
	}
	return cli.ShowAppHelp(ctx)
}
