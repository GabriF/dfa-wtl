package gotwodfawtl

import "strings"

type runeDoublyLinkedList struct {
	head *runeNode
	tail *runeNode
}

func (r *runeDoublyLinkedList) String() string {
	str := &strings.Builder{}
	current := r.head
	for current != nil {
		str.WriteRune(current.val)
		current = current.next
	}
	return str.String()
}

type runeNode struct {
	prev *runeNode
	next *runeNode
	val  rune
}

func newRuneNode(val rune) *runeNode {
	return &runeNode{nil, nil, val}
}

func fromStringDLL(s string) *runeDoublyLinkedList {
	ll := &runeDoublyLinkedList{nil, nil}
	for _, r := range s {
		newNode := newRuneNode(r)
		ll.insertEnd(newNode)
	}
	return ll
}

func (l *runeDoublyLinkedList) remove(node *runeNode) {
	if l.head == node && l.tail == node {
		l.head = nil
		l.tail = nil
	} else if l.head == node {
		l.head = node.next
		node.next.prev = nil
		node.next = nil
	} else if l.tail == node {
		l.tail = node.prev
		node.prev.next = nil
		node.prev = nil
	} else {
		node.prev.next = node.next
		node.next.prev = node.prev
		node.next = nil
		node.prev = nil
	}
}

func (l *runeDoublyLinkedList) insertEnd(node *runeNode) {
	switch l.tail {
	case nil:
		l.head = node
	case l.head:
		node.prev = l.head
		l.head.next = node
	default:
		node.prev = l.tail
		l.tail.next = node
	}
	l.tail = node
}
