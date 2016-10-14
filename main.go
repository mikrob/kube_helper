package main

// exemple of usage :
// go run main.go -c ~/.kube/config -k pod -n etcd0 -ns default

import (
	"errors"
	"flag"
	"fmt"
	"kube_helper/model"
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

func waitResource(resource model.KubeResource, clientSet *kubernetes.Clientset) (bool, error) {
	timeout := time.After(15 * time.Second)
	tick := time.Tick(2 * time.Second)
	for {
		select {
		case <-timeout:
			msg := fmt.Sprintf("Timeout after %d seconds", 15)
			return false, errors.New(msg)
		case <-tick:
			ok, err := resource.CheckState(clientSet)
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
		return model.Pod{Namespace: namespace, Name: resourceName}
	case "rc":
		return model.ReplicationController{Namespace: namespace, Name: resourceName}
	case "petset":
		return model.PetSet{Namespace: namespace, Name: resourceName}
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

	kubeResource := makeInstance(*resourceKind, *resourceName, *namespace).(model.KubeResource)

	_, podErr := waitResource(kubeResource, clientset)
	if podErr != nil {
		panic(podErr.Error())
	}
}
