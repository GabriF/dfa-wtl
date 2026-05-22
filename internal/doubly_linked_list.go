package internal

import (
	"fmt"
	"strings"
)

type Node struct {
	prev *Node
	next *Node
	val  int
}

func (n *Node) Val() int {
	return n.val
}

func newNode(val int) *Node {
	return &Node{nil, nil, val}
}

type IntDoublyLinkedListIter struct {
	current         *Node
	isDirectionNext bool
}

func (t *IntDoublyLinkedListIter) Next() *Node {
	if t.current == nil {
		return nil
	}
	toReturn := t.current
	if t.isDirectionNext {
		t.current = t.current.next
	} else {
		t.current = t.current.prev
	}
	return toReturn
}

type IntDoublyLinkedList struct {
	head *Node
	tail *Node
}

func (r *IntDoublyLinkedList) String() string {
	str := &strings.Builder{}
	it := r.TraverseFromStart()
	current := it.Next()
	for current != nil {
		str.WriteString(fmt.Sprint(current.Val()))
		current = it.Next()
		if current != nil {
			str.WriteString(" ")
		}
	}
	return str.String()
}

func (l *IntDoublyLinkedList) AtIndex(index int) *Node {
	current := l.head
	for range index {
		current = current.next
	}
	return current
}

func (l *IntDoublyLinkedList) First() *Node {
	return l.head
}

func (l *IntDoublyLinkedList) Last() *Node {
	return l.tail
}

func (l *IntDoublyLinkedList) Remove(node *Node) {
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

func (l *IntDoublyLinkedList) InsertEnd(val int) {
	node := newNode(val)
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

func (l *IntDoublyLinkedList) TraverseFromStart() *IntDoublyLinkedListIter {
	return &IntDoublyLinkedListIter{current: l.head, isDirectionNext: true}
}

func (l *IntDoublyLinkedList) TraverseFromEnd() *IntDoublyLinkedListIter {
	return &IntDoublyLinkedListIter{current: l.tail, isDirectionNext: false}
}
