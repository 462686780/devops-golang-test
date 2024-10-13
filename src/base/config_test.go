package base

import "testing"

func TestLoadConfig(t *testing.T) {
	tests := []string{"", "D:\\statefulset\\statefulset.yaml"}
	for _, test := range tests {
		if _, err := LoadConfig(test); err != nil {
			t.Error(`LoadConfig failed`)
		}
	}
}
