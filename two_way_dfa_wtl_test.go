package gotwodfawtl

import "testing"

func TestLinkedListOneElement(t *testing.T) {
	ll := fromStringDLL("a")
	if ll.String() != "a" {
		t.FailNow()
	}
}

func TestLinkedListRemove(t *testing.T) {
	ll := fromStringDLL("aba")
	toRemove := ll.head.next
	ll.remove(toRemove)
	if ll.String() != "aa" {
		t.FailNow()
	}
}

func TestLinkedListRemoveHead(t *testing.T) {
	ll := fromStringDLL("aba")
	toRemove := ll.head
	ll.remove(toRemove)
	if ll.String() != "ba" {
		t.FailNow()
	}
}

func TestLinkedListRemoveTail(t *testing.T) {
	ll := fromStringDLL("aba")
	toRemove := ll.tail
	ll.remove(toRemove)
	if ll.String() != "ab" {
		t.FailNow()
	}
}
