package model

import (
	"fmt"

	"k8s.io/client-go/1.4/kubernetes"
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

// CheckState for PetSet
func (ps PetSet) CheckState(clientSet *kubernetes.Clientset) (bool, error) {
	petset, err := clientSet.Apps().PetSets(ps.Namespace).Get(ps.Name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", petset)
	return true, nil
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
