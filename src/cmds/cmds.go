package cmds

import (
	"fmt"
	"os"
)

func Execute() {
	app := App()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
