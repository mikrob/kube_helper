package main

// exemple of usage :
// go run main.go -c ~/.kube/config -k pod -n etcd0 -ns default

import (
	"errors"
	"flag"
	"fmt"
	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/tools/clientcmd"
	"reflect"
	"time"
)

var (
	kubeconfig     = flag.String("c", "./config", "absolute path to the kubeconfig file")
	resourceKind   = flag.String("k", "rc", "resource kind name : pod, svc, rc, petset ...")
	resourceName   = flag.String("n", "mysuperpod", "resource name e.g. pod name, rc name ...")
	namespace      = flag.String("ns", "namespace", "namespace of the resource")
	resourceMapper = map[string]reflect.Type{
		"pod": reflect.TypeOf(Pod{}),
		"rc":  reflect.TypeOf(ReplicationController{}),
	}
)

type KubeResource interface {
	checkState() (bool, error)
	init(clientSet *kubernetes.Clientset, namespace string, name string)
}

type Pod struct {
	clientSet *kubernetes.Clientset
	namespace string
	name      string
}

func (p Pod) init(clientSet *kubernetes.Clientset, namespace string, name string) {
	p.clientSet = clientSet
	p.namespace = namespace
	p.name = name
}

func (p Pod) checkState() (bool, error) {
	a := p.clientSet.Core()
	fmt.Printf("----> %+v\n", a)

	pod, err := p.clientSet.Core().Pods(p.namespace).Get(p.name)
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

type ReplicationController struct {
}

func (c ReplicationController) init(clientSet *kubernetes.Clientset, namespace string, name string) {
	// TODO: implement this
}

func (c ReplicationController) checkState() (bool, error) {
	// TODO: implement this
	return true, nil
}

func waitResource(resource KubeResource) (bool, error) {
	timeout := time.After(15 * time.Second)
	tick := time.Tick(2 * time.Second)
	for {
		select {
		case <-timeout:
			msg := fmt.Sprintf("Timeout after %d seconds", 15)
			return false, errors.New(msg)
		case <-tick:
			ok, err := resource.checkState()
			if err != nil {
				return false, err
			} else if ok {
				return true, nil
			}
		}
	}
}

func makeInstance(name string) interface{} {
	v := reflect.New(resourceMapper[name]).Elem()
	return v.Interface()
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

	var kubeResource KubeResource
	if _, ok := resourceMapper[*resourceKind]; ok {
		kubeResource = makeInstance(*resourceKind).(KubeResource)
		kubeResource.init(clientset, *namespace, *resourceName)
	} else {
		panic(fmt.Sprintf("Dont known what to do with kind : %s", *resourceKind))
	}

	_, podErr := waitResource(kubeResource)
	if podErr != nil {
		panic(podErr.Error())
	}
}
