package main

// exemple of usage :
// go run main.go -c ~/.kube/config -k pod -n etcd0 -ns default

import (
	"errors"
	"flag"
	"fmt"
	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/tools/clientcmd"
	"time"
	//  "k8s.io/client-go/1.4/pkg/api/v1"
)

var (
	kubeconfig    = flag.String("c", "./config", "absolute path to the kubeconfig file")
	resource_kind = flag.String("k", "rc", "resource kind name : pod, svc, rc, petset ...")
	resource_name = flag.String("n", "mysuperpod", "resource name e.g. pod name, rc name ...")
	namespace     = flag.String("ns", "namespace", "namespace of the resource")
)

func waitPod(clientset *kubernetes.Clientset, namespace string, podname string) (bool, error) {
	timeout := time.After(15 * time.Second)
	tick := time.Tick(2 * time.Second)
	for {
		select {
		case <-timeout:
			msg := fmt.Sprintf("Timeout after %d seconds", 15)
			return false, errors.New(msg)
		case <-tick:
			ok, err := podRunning(clientset, namespace, podname)
			if err != nil {
				return false, err
			} else if ok {
				return true, nil
			}
		}
	}
}

func podRunning(clientset *kubernetes.Clientset, namespace string, podname string) (bool, error) {
	pod, err := clientset.Core().Pods(namespace).Get(podname)
	fmt.Println("Namespace :", namespace)
	fmt.Println("Pod:", podname)
	if err != nil {
		panic(err.Error())
	}
	state := pod.Status.Phase
	//fmt.Printf("%+v\n", status)
	fmt.Println("Status :", state)
	return state == "Prout", err
}

func main() {
	flag.Parse()

	fmt.Println("Resource Kind is :", *resource_kind)
	fmt.Println("Resource Name is : ", *resource_name)
	fmt.Println("Namespace is : ", *namespace)
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	_, pod_err := waitPod(clientset, *namespace, *resource_name)
	if pod_err != nil {
		panic(pod_err.Error())
	}
}
