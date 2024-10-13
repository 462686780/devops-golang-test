package base

import "testing"

func TestInitContext(t *testing.T) {
	confFile := ""
	cnf, err := LoadConfig(confFile)
	if err != nil {
		t.Error(`LoadConfig failed`)
	}
	if err := InitContext(cnf); err != nil {
		t.Error(`LoadConfig failed`)
	}
}
