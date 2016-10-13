package main

// exemple of usage :
// go run main.go -c ~/.kube/config -k pod -n etcd0 -ns default

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/tools/clientcmd"
)

var (
	kubeconfig   = flag.String("c", "./config", "absolute path to the kubeconfig file")
	resourceKind = flag.String("k", "rc", "resource kind name : pod, svc, rc, petset ...")
	resourceName = flag.String("n", "mysuperpod", "resource name e.g. pod name, rc name ...")
	namespace    = flag.String("ns", "namespace", "namespace of the resource")
)

// KubeResource is an abstract representation of kuberesource
type KubeResource interface {
	checkState(clientSet *kubernetes.Clientset) (bool, error)
}

//Pod is our internal representation of Pod
type Pod struct {
	namespace string
	name      string
}

//ReplicationController is our internal replication of ReplicationController
type ReplicationController struct {
	namespace string
	name      string
}

type PetSet struct {
	namespace string
	name      string
}

func (p Pod) checkState(clientSet *kubernetes.Clientset) (bool, error) {
	pod, err := clientSet.Core().Pods(p.namespace).Get(p.name)
	fmt.Println("Namespace :", p.namespace)
	fmt.Println("Pod:", p.name)
	if err != nil {
		panic(err.Error())
	}
	state := pod.Status.Phase
	//fmt.Printf("%+v\n", status)
	fmt.Println("Status :", state)
	return state == "Prout", err
}

func (rc ReplicationController) checkState(clientSet *kubernetes.Clientset) (bool, error) {
	repcontroller, err := clientSet.Core().ReplicationControllers(rc.namespace).Get(rc.name)
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

func (ps PetSet) checkState(clientSet *kubernetes.Clientset) (bool, error) {
	// ps, err := clientSet.Core().PetSets(ps.namespace).Get(ps.name)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// fmt.Printf("%+v\n", ps)
	return true, nil
}

func waitResource(resource KubeResource, clientSet *kubernetes.Clientset) (bool, error) {
	timeout := time.After(15 * time.Second)
	tick := time.Tick(2 * time.Second)
	for {
		select {
		case <-timeout:
			msg := fmt.Sprintf("Timeout after %d seconds", 15)
			return false, errors.New(msg)
		case <-tick:
			ok, err := resource.checkState(clientSet)
			if err != nil {
				return false, err
			} else if ok {
				return true, nil
			}
		}
	}
}

func makeInstance(kind string, resourceName string, namespace string) interface{} {
	switch kind {
	case "pod":
		return Pod{namespace: namespace, name: resourceName}
	case "rc":
		return ReplicationController{namespace: namespace, name: resourceName}
	default:
		fmt.Println(fmt.Errorf("Don't know what to do with kind : %s", kind))
		return nil
	}
}

func main() {
	flag.Parse()

	fmt.Println("Resource Kind is :", *resourceKind)
	fmt.Println("Resource Name is : ", *resourceName)
	fmt.Println("Namespace is : ", *namespace)

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	kubeResource := makeInstance(*resourceKind, *resourceName, *namespace).(KubeResource)

	_, podErr := waitResource(kubeResource, clientset)
	if podErr != nil {
		panic(podErr.Error())
	}
}
