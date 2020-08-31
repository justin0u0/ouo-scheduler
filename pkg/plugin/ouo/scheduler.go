package ouo

import (
	"context"
	"strconv"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/api/v1/pod"
	"k8s.io/kubernetes/pkg/apis/core/v1/helper/qos"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

// Name ... the custom shceduler name
const Name = "ouo-scheduler"

// CustomScheduler ... The type CustomScheduler implement the interface of the kube-scheduler framework
type CustomScheduler struct {
	handle framework.FrameworkHandle
}

// Let the type CustomScheduler implement the QueueSortPlugin, PreFilterPlugin interface
var _ framework.QueueSortPlugin = &CustomScheduler{}
var _ framework.PreFilterPlugin = &CustomScheduler{}

// Name ... Implement Plugin interface Name() @pkg/scheduler/framework/v1alpha1/interface.go
func (*CustomScheduler) Name() string {
	return Name
}

// Less ... Implement QueueSortPlugin interface Less() @pkg/scheduler/framework/v1alpha1/interface.go
func (*CustomScheduler) Less(pInfo1, pInfo2 *framework.PodInfo) bool {
	p1 := pod.GetPodPriority(pInfo1.Pod)
	p2 := pod.GetPodPriority(pInfo2.Pod)

	klog.V(3).Infof("[queue sort] [Less]: %v: %v, %v: %v", pInfo1.Pod.Name, p1, pInfo2.Pod.Name, p2)

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

	klog.V(3).Infof("[queue sort] [comparePodQOS]: %v: %v, %v, %v", pod1.Name, pod1QOS, pod2.Name, pod2QOS)

	if pod1QOS == v1.PodQOSGuaranteed {
		return true
	}
	if pod1QOS == v1.PodQOSBurstable {
		return pod2QOS != v1.PodQOSGuaranteed
	}
	return pod1QOS == pod2QOS
}

// PreFilter ... Implement PreFilterPlugin interface PreFilter()
func (s *CustomScheduler) PreFilter(_ context.Context, _ *framework.CycleState, pod *v1.Pod) *framework.Status {
	nodeInfos, err := s.handle.SnapshotSharedLister().NodeInfos().List()
	if err != nil {
		return framework.NewStatus(framework.Unschedulable, "Failed to pass pre filter, get nodeInfos failed.")
	}

	podGroup, podGroupExist := pod.Labels["podGroup"]
	if !podGroupExist {
		return framework.NewStatus(framework.Success, "Pass pre filter successfully, pod has no label podGroup.")
	}

	minAvailable, minAvailableExist := pod.Labels["minAvailable"]
	if !minAvailableExist {
		return framework.NewStatus(framework.Success, "Pass pre filter successfully, pod has no label minAvailable.")
	}
	minAvailableNum, atoiErr := strconv.Atoi(minAvailable)
	if atoiErr != nil {
		return framework.NewStatus(framework.Unschedulable, "Failed to pass pre filter, pod label minAvailable is not a valid number")
	}

	totalPodsInPodGroup := getTotalPodsByPodGroup(nodeInfos, pod.Namespace, podGroup)
	if totalPodsInPodGroup < minAvailableNum {
		klog.V(3).Infof("The count of PodGroup %v (%v) is less than minAvailable(%d) in PreFilter: %d", podGroup, pod.Name, minAvailableNum, totalPodsInPodGroup)
		return framework.NewStatus(framework.Unschedulable, "Failed to pass pre filter, less than min available")
	}

	return framework.NewStatus(framework.Success, "Pass pre filter successfully")
}

// PreFilterExtensions ...
func (*CustomScheduler) PreFilterExtensions() framework.PreFilterExtensions {
	return nil
}

func getTotalPodsByPodGroup(nodeInfos []*nodeinfo.NodeInfo, ns string, pg string) int {
	total := 0
	for _, nodeInfo := range nodeInfos {
		for _, pod := range nodeInfo.Pods() {
			if podGroup, ok := pod.Labels["podGroup"]; ok && podGroup == pg && pod.Namespace == ns {
				total++
			}
		}
	}
	return total
}

// New ... Create an scheduler instance
// New() is type PluginFactory = func(configuration runtime.Object, f v1alpha1.FrameworkHandle) (v1alpha1.Plugin, error)
// mentioned in https://github.com/kubernetes/kubernetes/blob/master/pkg/scheduler/framework/runtime/registry.go
func New(_ runtime.Object, handle framework.FrameworkHandle) (framework.Plugin, error) {
	return &CustomScheduler{
		handle: handle,
	}, nil
}
