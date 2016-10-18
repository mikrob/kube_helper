package model

import (
	"fmt"
	"regexp"

	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/pkg/api"
)

// KubeResource is an abstract representation of kuberesource
type KubeResource interface {
	CheckState(clientSet *kubernetes.Clientset) (bool, error)
}

// Pod is our internal representation of Pod
type Pod struct {
	Namespace string
	Name      string
}

type Service struct {
	Namespace string
	Name      string
}

// ReplicationController is our internal replication of ReplicationController
type ReplicationController struct {
	Namespace string
	Name      string
}

// PetSet is our internal representation of PetSet
type PetSet struct {
	Namespace string
	Name      string
}

func allPredicate(vs []bool, f func(bool) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

// CheckState for service
func (svc Service) CheckState(clientSet *kubernetes.Clientset) (bool, error) {
	return true, nil
}

// CheckState for PetSet
func (ps PetSet) CheckState(clientSet *kubernetes.Clientset) (bool, error) {
	petSet, err := clientSet.Apps().PetSets(ps.Namespace).Get(ps.Name)

	if err != nil {
		panic(err.Error())
	}
	wantedReplicas := petSet.Spec.Replicas
	psStatusReplicas := petSet.Status.Replicas
	fmt.Println("Wanted replicas : ", *wantedReplicas)
	fmt.Println("Actual Replicas :", psStatusReplicas)

	petSetPods, errPods := clientSet.Core().Pods(ps.Namespace).List(api.ListOptions{})
	if errPods != nil {
		panic(errPods.Error())
	}
	//fmt.Println("TYPE :", reflect.TypeOf(petSetPods.Items))
	// fmt.Printf("%+v\n", petSetPods.Items())
	var podReadys []bool
	for _, pod := range petSetPods.Items {
		pName := string(pod.Name)
		match, _ := regexp.MatchString(ps.Name, pName)
		if match {
			fmt.Println(fmt.Sprintf("Pod : %s, Status : %s", pName, pod.Status.Phase))
			podReadys = append(podReadys, pod.Status.Phase == "Running")
		}
	}

	//podReadys = append(podReadys, false)
	allReady := allPredicate(podReadys, func(v bool) bool {
		return v
	})

	return allReady && *wantedReplicas == psStatusReplicas, nil
}

// CheckState for pod
func (p Pod) CheckState(clientSet *kubernetes.Clientset) (bool, error) {
	pod, err := clientSet.Core().Pods(p.Namespace).Get(p.Name)
	fmt.Println("Namespace :", p.Namespace)
	fmt.Println("Pod:", p.Name)
	if err != nil {
		panic(err.Error())
	}
	state := pod.Status.Phase
	//fmt.Printf("%+v\n", status)
	fmt.Println("Status :", state)
	return state == "Running", err
}

//CheckState for ReplicationController
func (rc ReplicationController) CheckState(clientSet *kubernetes.Clientset) (bool, error) {
	repcontroller, err := clientSet.Core().ReplicationControllers(rc.Namespace).Get(rc.Name)
	if err != nil {
		panic(err.Error())
	}
	state := repcontroller.Status
	//replicas := int32(4) //state.Replicas
	replicas := state.Replicas
	readyReplicas := state.ReadyReplicas
	fullyLabeledReplicas := state.FullyLabeledReplicas
	fmt.Println("Replicas :", replicas)
	fmt.Println("Ready replicas:", readyReplicas)
	return (replicas == readyReplicas) && (readyReplicas == fullyLabeledReplicas), err
}
