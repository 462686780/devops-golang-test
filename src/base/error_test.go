package base

import (
	"flag"
	"fmt"
	"testing"

	"github.com/urfave/cli"
)

func TestOnUsageError(t *testing.T) {
	tests := []string{"err is test"}
	for _, test := range tests {
		app := cli.NewApp()
		set := flag.NewFlagSet(app.Name, flag.ContinueOnError)
		context := cli.NewContext(app, set, nil)
		err := OnUsageError(context, fmt.Errorf(test), true)
		if err != nil {
			t.Error(`LoadConfig failed`)
		}
	}
}
