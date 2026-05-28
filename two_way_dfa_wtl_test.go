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
		{automatons()[0], "ab", true},
		{automatons()[0], "abaabbcc", true},
		{automatons()[0], "ababbccc", false},

		{automatons()[1], "", true},
		{automatons()[1], "ababbbbaba", true},
		{automatons()[1], "ababbbaba", false},
	}

	for _, tt := range tests {
		id := NewComputation(tt.m, tt.word)
		for !id.Halt {
			t.Log(id.String())
			ComputeNext(id)
		}
		correct := (id.Accept && tt.accept) || ((!id.Accept) && (!tt.accept))
		if !correct {
			t.Errorf("Expected %t on input %s with automaton %+v", tt.accept, tt.word, tt.m)
		}
	}
}

func makeBalancedWord(n int, rng *rand.Rand) string {
	str := &strings.Builder{}
	aCount, bCount := 0, 0
	for i := 0; i < n*2; i++ {
		if aCount == n {
			str.WriteRune('b')
		} else if bCount == n {
			str.WriteRune('a')
		} else if rng.Intn(2) == 0 {
			str.WriteRune('a')
			aCount++
		} else {
			str.WriteRune('b')
			bCount++
		}
	}
	for i := 0; i < n-1; i++ {
		str.WriteRune('c')
	}
	return str.String()
}

func BenchmarkAutomatonComputation(b *testing.B) {
	m := automatons()[0]
	n := 1000000000

	rng := rand.New(rand.NewSource(1))
	inputs := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		inputs[i] = makeBalancedWord(n, rng)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := NewComputation(m, inputs[i])
		for !id.Halt {
			ComputeNext(id)
		}
		if !id.Accept {
			b.FailNow()
		}
	}
}
