package ouo

import (
  "k8s.io/api/core/v1"
  "k8s.io/kubernetes/pkg/api/v1/pod"
  "k8s.io/kubernetes/apis/core/v1/helper/qos"
  framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
  "k8s.io/apimachinery/pkg/runtime"
)

const (
  Name = "ouo-scheduler"
)

// The type CustomScheduler implement the interface of the kube-scheduler framework
type CustomScheduler struct {}

// Let the type CustomScheduler implement the QueueSortPlugin interface
var _ framework.QueueSortPlugin = &CustomScheduler{}

// Implement Plugin interface Name() @pkg/scheduler/framework/v1alpha1/interface.go
func (_ *CustomScheduler) Name() string {
  return Name
}

// Implement QueueSortPlugin interface Less() @pkg/scheduler/framework/v1alpha1/interface.go
func (_ *CustomScheduler) Less(pInfo1, pInfo2 *framework.PodInfo) bool {
  p1 := pod.GetPodPriority(pInfo1.Pod)
  p2 := pod.GetPodPriority(pInfo2.Pod)
  return (p1 > p2) || (p1 == p2 && comparePodQOS(pInfo1.Pod, pInfo2.Pod))
}

// As https://github.com/kubernetes/kubernetes/blob/master/pkg/scheduler/framework/v1alpha1/types.go,
// we can know that the type of framework.PodInfo.Pod is *v1.Pod.
// The GetPodQOS function is implement in https://github.com/kubernetes/kubernetes/blob/master/pkg/apis/core/v1/helper/qos/qos.go.
// PodQOSGuaranteed is the Guaranteed qos class, PodQOSBurstable is the Burstable qos class, PodQOSBestEffort is the BestEffort qos class.
// Mentioned in https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/api/core/v1/types.go
func comparePodQOS(pod1, pod2 *v1.Pod) bool {
  pod1QOS := qos.GetPodQOS(pod1)
  pod2QOS := qos.GetPodQOS(pod2)
  if pod1QOS == PodQOSGuaranteed {
    return true
  }
  if pod1QOS == PodQOSBurstable {
    return pod2QOS != PodQOSGuaranteed
  }
  return pod1QOS == pod2QOS
}

// Create an scheduler instance
// New() is type PluginFactory = func(configuration runtime.Object, f v1alpha1.FrameworkHandle) (v1alpha1.Plugin, error)
// mentioned in https://github.com/kubernetes/kubernetes/blob/master/pkg/scheduler/framework/runtime/registry.go
func New(_ runtime.Object, _ framework.FrameworkHandle) (framework.Plugin, error) {
	return &CustomScheduler{}, nil
}
