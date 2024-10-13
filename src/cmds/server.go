package cmds

import (
	"statefulset/base"
	"statefulset/cmds/server"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/urfave/cli"
)

//CmdServer listen and run server
var CmdServer = cli.Command{
	Name:        "server",
	Usage:       "Listen and run server",
	Description: "Listen and run server",
	Action:      runServer,
}

func runServer(ctx *cli.Context) error {
	confFile := ctx.String("conf")
	cnf, err := base.LoadConfig(confFile)
	if err != nil {
		return err
	}

	err = base.InitContext(cnf)
	if err != nil {
		return err
	}

	base.Context.Logger.Debugf("NewMysqlDelivery Load conf")

	if cnf.RunMode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	//new server
	server := server.NewApiServer(base.Context.Logger)
	server.InitRoute()
	return server.Start()
}
