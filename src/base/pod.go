package base

import (
	"fmt"
	"sync"
)

// PersistentVolume simulates a persistent volume
type PersistentVolume struct {
	Name  string
	Claim string
}

// simulates a Pod
type Pod struct {
	Name                 string
	PersistentVolumeName string
}

// MyStatefulSet controller
type MyStatefulSet struct {
	Name              string
	Replicas          int
	Pods              []Pod
	PersistentVolumes map[string]*PersistentVolume
	PodMutex          sync.Mutex
}

// new MyStatefulSet
func NewMyStatefulSet(name string, replicas int) *MyStatefulSet {
	return &MyStatefulSet{
		Name:              name,
		Replicas:          replicas,
		PersistentVolumes: make(map[string]*PersistentVolume),
	}
}

// CreatePod
func (s *MyStatefulSet) CreatePod(index int) {
	s.PodMutex.Lock()
	defer s.PodMutex.Unlock()

	// add volume
	pvName := fmt.Sprintf("%s-pv-%d", s.Name, index)
	pv := &PersistentVolume{Name: pvName, Claim: ""}
	s.PersistentVolumes[pvName] = pv

	pod := Pod{
		Name:                 fmt.Sprintf("%s-%d", s.Name, index),
		PersistentVolumeName: pvName,
	}
	s.Pods = append(s.Pods, pod)

	Log.Debugf("Created Pod %s with PersistentVolume %s\n", pod.Name, pv.Name)
}

// DeletePod
func (s *MyStatefulSet) DeletePod(name string) {
	s.PodMutex.Lock()
	defer s.PodMutex.Unlock()

	for i, pod := range s.Pods {
		if pod.Name == name {
			pvName := pod.PersistentVolumeName
			delete(s.PersistentVolumes, pvName)
			s.Pods = append(s.Pods[:i], s.Pods[i+1:]...)
			Log.Debugf("Deleted Pod %s and released PersistentVolume %s\n", name, pvName)
			return
		}
	}
	Log.Debugf("Pod %s not found\n", name)
}

// ListPods
func (s *MyStatefulSet) ListPods() {
	s.PodMutex.Lock()
	defer s.PodMutex.Unlock()

	for _, pod := range s.Pods {
		Log.Debugf("Pod %s with PersistentVolume %s\n", pod.Name, pod.PersistentVolumeName)
	}
}
