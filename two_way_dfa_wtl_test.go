package gotwodfawtl_test

import (
	"math/rand"
	"strings"
	"testing"

	. "github.com/GabriF/dfa-wtl"
)

func automatons() []*TwoWayDfaWtl {
	return []*TwoWayDfaWtl{FromSymbolDescription(
		[]string{"qa", "qb"},
		[]string{"qc"},
		"qa",
		[]rune{'a', 'b', 'c'},
		']',
		'[',
		map[string]map[rune]string{
			"qa": {
				'a': "qb",
			},
			"qb": {
				'b': "qc",
			},
			"qc": {
				'c': "qa",
				']': "accept",
			},
		},
		map[string]map[rune]bool{
			"qa": {
				'b': true,
			},
			"qb": {
				'a': true,
			},
		},
	),
		FromSymbolDescription(
			[]string{"q0", "qa", "qb"},
			[]string{"pa", "pb", "pr"},
			"q0",
			[]rune{'a', 'b'},
			']',
			'[',
			map[string]map[rune]string{
				"q0": {
					'[': "accept",
					'a': "qa",
					'b': "qb",
				},
				"qa": {
					'[': "pa",
				},
				"qb": {
					'[': "pb",
				},
				"pa": {
					'a': "pr",
				},
				"pb": {
					'b': "pr",
				},
				"pr": {
					']': "q0",
				},
			},
			map[string]map[rune]bool{
				"qa": {
					'a': true,
					'b': true,
				},
				"qb": {
					'a': true,
					'b': true,
				},
				"pr": {
					'a': true,
					'b': true,
				},
			},
		),
	}
}

func TestAutomatonComputation(t *testing.T) {
	tests := []struct {
		m      *TwoWayDfaWtl
		word   string
		accept bool
	}{
		{automatons()[0], "abaabbcc", true},
		{automatons()[0], "ababbccc", false},

		{automatons()[1], "ababbbbaba", true},
		{automatons()[1], "ababbbaba", false},
	}

	for _, tt := range tests {
		id := NewComputation(tt.m, tt.word)
		for !id.Halt {
			ComputeNext(tt.m, id)
		}
		correct := (id.Accept && tt.accept) || ((!id.Accept) && (!tt.accept))
		if !correct {
			t.Errorf("Expected %t on input %s with automaton %+v", tt.accept, tt.word, tt.m)
		}
	}
}

func BenchmarkAutomatonComputationBenchmark(b *testing.B) {
	m := automatons()[0]
	n := 100000000

	for b.Loop() {
		b.StopTimer()

		str := &strings.Builder{}
		a_count, b_count := 0, 0
		for range n * 2 {
			if a_count == n {
				str.WriteRune('b')
			} else if b_count == n {
				str.WriteRune('a')
			} else if rand.Intn(2) == 0 {
				str.WriteRune('a')
				a_count++
			} else {
				str.WriteRune('b')
				b_count++
			}
		}
		for range n - 1 {
			str.WriteRune('c')
		}
		w := str.String()

		b.StartTimer()

		id := NewComputation(m, w)
		for !id.Halt {
			ComputeNext(m, id)
		}
		if !id.Accept {
			b.FailNow()
		}
	}
}
