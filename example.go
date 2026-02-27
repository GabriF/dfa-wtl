package main

import (
	"dfawtl/two_way_dfa_wtl"
	"fmt"
)

func twoWayDfaWtl() {
	w := "aababbcc"

	b := two_way_dfa_wtl.NewBuilder([]string{"q_a", "q_b"}, []string{"q_c"}, "q_a", ']', '[')
	b.WithDelta("q_a", 'a', "q_b")
	b.WithDelta("q_b", 'b', "q_c")
	b.WithDelta("q_c", 'c', "q_a")
	b.WithDelta("q_c", ']', "accept")
	b.WithTau("q_a", 'b')
	b.WithTau("q_b", 'a')

	a, _ := b.Build()

	for _, s := range a.Exec(w) {
		fmt.Println(s)
	}
}

func main() {
	twoWayDfaWtl()
}
