package v1

import (
	"reflect"
	"testing"
)

func TestReverse(t *testing.T) {
	n1 := &node{value: 1}
	n2 := &node{value: 2}
	n3 := &node{value: 3}

	nodes := Nodes([]*node{n3, n2, n1}).Reverse()
	expected := Nodes([]*node{n1, n2, n3})

	if !reflect.DeepEqual(nodes, expected) {
		t.Errorf(`Expected: %+v, received: %+v`, expected, nodes)
	}

	nodes = Nodes([]*node{n2, n1}).Reverse()
	expected = Nodes([]*node{n1, n2})

	if !reflect.DeepEqual(nodes, expected) {
		t.Errorf(`Expected: %+v, received: %+v`, expected, nodes)
	}
}
