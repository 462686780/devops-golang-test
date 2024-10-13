package base

import (
	"fmt"
	"testing"
)

func TestLogger(t *testing.T) {
	confFile := ""
	cnf, err := LoadConfig(confFile)
	if err != nil {
		t.Error(`LoadConfig failed`)
	}
	fmt.Printf("conf dir:%v,name:%v,level:%v", cnf.Logger.Dir, cnf.Logger.Name, cnf.Logger.Level)
	InitLogger(cnf.Logger.Dir, cnf.Logger.Name, cnf.Logger.Level)
	Log.Debugf("this is debug")
}
