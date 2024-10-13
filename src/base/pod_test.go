package base

import "testing"

func TestMyStatefulSet(t *testing.T) {
	confFile := ""
	cnf, err := LoadConfig(confFile)
	if err != nil {
		t.Error(`LoadConfig failed`)
	}
	InitLogger(cnf.Logger.Dir, cnf.Logger.Name, cnf.Logger.Level)
	Log.Debugf("this is debug")
	statefulSet := NewMyStatefulSet("my-statefulset", 3)

	for i := 0; i < statefulSet.Replicas; i++ {
		statefulSet.CreatePod(i)
	}

	statefulSet.ListPods()

	statefulSet.DeletePod("my-statefulset-1")

	statefulSet.ListPods()
}
