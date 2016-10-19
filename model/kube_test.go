package model

import (
	"testing"

	fakeclientset "k8s.io/client-go/1.4/kubernetes/fake"
)

func TestAllPredicateTrue(t *testing.T) {
	boolArray := []bool{true, true, true}
	allReady := allPredicate(boolArray, func(v bool) bool {
		return v
	})
	if !allReady {
		t.Error("Test failed")
	}
}

func TestAllPredicateWithOneFalse(t *testing.T) {
	boolArray := []bool{true, true, false}
	allReady := allPredicate(boolArray, func(v bool) bool {
		return v
	})
	if allReady {
		t.Error("Test failed, should be false")
	}
}

func TestCheckService(t *testing.T) {
	kubeResource := Service{Name: "myservice", Namespace: "mynamespace"}
	clientSet := fakeclientset.NewSimpleClientset()
	result, resultErr := kubeResource.CheckState(clientSet)
	t.Log("prout")
	t.Logf("%+v\n", result)
	if resultErr != nil {
		t.Error("Test failed")
	}
}
