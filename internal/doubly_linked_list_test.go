package internal_test

import (
	"testing"

	. "github.com/GabriF/dfa-wtl/internal"
)

func TestDLLOneElement(t *testing.T) {
	ll := &IntDoublyLinkedList{}
	ll.InsertEnd(1)
	if ll.First().Val() != 1 {
		t.FailNow()
	}
}

func TestLinkedListRemove(t *testing.T) {
	expected := "0 2"

	ll := &IntDoublyLinkedList{}
	ll.InsertEnd(0)
	ll.InsertEnd(1)
	ll.InsertEnd(2)

	ll.Remove(ll.AtIndex(1))
	actual := ll.String()
	if actual != expected {
		t.Fatalf("Expected %s found %s", expected, actual)
	}
}

func TestLinkedListRemoveHead(t *testing.T) {
	expected := "1 2"

	ll := &IntDoublyLinkedList{}
	ll.InsertEnd(0)
	ll.InsertEnd(1)
	ll.InsertEnd(2)
	ll.Remove(ll.First())

	actual := ll.String()
	if actual != expected {
		t.Fatalf("Expected %s found %s", expected, actual)
	}
}

func TestLinkedListRemoveTail(t *testing.T) {
	expected := "0 1"

	ll := &IntDoublyLinkedList{}
	ll.InsertEnd(0)
	ll.InsertEnd(1)
	ll.InsertEnd(2)
	ll.Remove(ll.Last())

	actual := ll.String()
	if actual != expected {
		t.Fatalf("Expected %s found %s", expected, actual)
	}
}
