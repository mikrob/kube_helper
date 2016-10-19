package model

import "testing"

func TestKube(t *testing.T) {
	expected := "Hello Go!"
	actual := "Hello"
	if actual != expected {
		t.Error("Test failed")
	}
}
