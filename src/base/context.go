package base

import (
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yu31/snowflake"
)

//global config
var Context GlobalContext

type GlobalContext struct {
	Setting *GlobalSetting
	Logger  *logrus.Logger
}

// InitContext initialization the global variables
// TODO: Connection is established when invoking for mysql, redis, etc.
func InitContext(settings *GlobalSetting) (err error) {
	// Setup global settings
	Context.Setting = settings

	// Setup global default Logger
	Context.Logger = NewLogger(settings.Logger.Dir, settings.Logger.Name, settings.Logger.Level)
	// Setup global request id generator worker
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(100)
	reqIDWorker, err = snowflake.New(int64(num))
	if err != nil {
		Context.Logger.Infof("snowflake new err:%v", err)
		return err
	}

	return nil
}
