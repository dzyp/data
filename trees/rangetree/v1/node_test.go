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

func TestShiftLeft(t *testing.T) {
	p := newOrderedList(2)
	p1 := newPoint(1, 1)
	p.insert(p1)
	n := &node{
		value: 1,
		p:     p,
	}

	tree := New(2)

	shiftLeft(n, 2, tree)

	if n.value != 2 {
		t.Errorf(`Expected: %d, received: %d`, 2, n.value)
	}

	if n.left.value != 1 {
		t.Errorf(`Expected: %d, received: %d`, 1, n.left.value)
	}

	if n.right.value != 2 {
		t.Errorf(`Expected: %d, received: %d`, 2, n.right.value)
	}

	if n.p.len() != 1 {
		t.Errorf(`Expected len: %d, received: %d`, 1, n.p.len())
	}

	if n.right.p.len() != 0 {
		t.Errorf(`Expected len: %d, received: %d`, 0, n.right.p.len())
	}

	if n.left.p.len() != 1 {
		t.Errorf(`Expected len: %d, received: %d`, 1, n.left.p.len())
	}
}

func TestShiftRight(t *testing.T) {
	p := newOrderedList(2)
	p1 := newPoint(1, 1)
	p.insert(p1)
	n := &node{
		value: 1,
		p:     p,
	}

	tree := New(2)

	shiftRight(n, 0, tree)

	if n.value != 1 {
		t.Errorf(`Expected: %d, received: %d`, 1, n.value)
	}

	if n.left.value != 0 {
		t.Errorf(`Expected: %d, received: %d`, 0, n.left.value)
	}

	if n.right.value != 1 {
		t.Errorf(`Expected: %d, received: %d`, 1, n.right.value)
	}

	if n.p.len() != 1 {
		t.Errorf(`Expected len: %d, received: %d`, 1, n.p.len())
	}

	if n.right.p.len() != 1 {
		t.Errorf(`Expected len: %d, received: %d`, 1, n.right.p.len())
	}

	if n.left.p.len() != 0 {
		t.Errorf(`Expected len: %d, received: %d`, 0, n.left.p.len())
	}
}
