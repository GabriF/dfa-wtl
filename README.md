# dfa-wtl

This repository contains implementations for some finite automata with
translucent input letters.

## Two-Way Finite Automata with Translucent Input Letters (2DFAwtl)

The package `two_way_dfa_wtl` contains the functions to create and execute a
2DFAwtl. The remainder of this paragraph explains how to use it with the
following automata as example.

Let $`M = (\{q_a, q_b, q_c \}, \{a, b, c\}, q_a, ], [, \tau, \delta)`$ with:

- $`Q_r = \{q_a, q_b\}`$
- $`Q_l = \{q_c\}`$
- ] left endmarker
- [ right endmarker
- $\tau(q_a) = \{b\}, \tau(q_b) = \{a\}$
- $\delta(q_a, a) = q_b, \delta(q_b, b) = q_c, \delta(q_c, c) = q_a, \delta(q_c,
  ]) = \text{accept}$

which accepts the language $`L = \{uc^{n-1} : u \in \{a,b \}^{\star}, |u|_a = n,
|u|_b = n, n > 0\}`$.

First, initialize a builder:

```
b := two_way_dfa_wtl.NewBuilder([]string{"q_a", "q_b"}, []string{"q_c"}, "q_a", ']', '[')
```

where the first and second parameters are the slice containing the name of the
state in $Q_r$ and $Q_l$, respectively; the third parameter is name of the
initial state; the last two parameters are the left and right endmarkers,
respectively. Please note that `accept` is a special state used later, thus the
user must not define a state named `accept`.

Then define $\delta$ as follows:

```
b.WithDelta("q_a", 'a', "q_b")
b.WithDelta("q_b", 'b', "q_c")
b.WithDelta("q_c", 'c', "q_a")
b.WithDelta("q_c", ']', "accept")
```

Where the first paramter is the current state, the second parameter is the
current input letter and the last paramter is the next state. The string
`accept` may be used to indicate acceptance in the given state on the given
input letter.

Then define $\tau$ as follows:

```
b.WithTau("q_a", 'b')
b.WithTau("q_b", 'a')
```

Where the first parameter is a state and the second parameter is a list of
runes. (e.g. not shown in the current example, one could use `b.WithTau("q",
'a', 'b', 'c')`)

Finally the automata can be obtained with:

```
automata, error := b.Build()
```

Note that `Build()` resets the builder.

Once the automata is obtained it can be executed on input string `w` as follow

```
automata.Exec(w)
```

Which returns a slice of strings, where the last string is either `æccept` or
`reject` and the other strings are the configurations the automata went
through during the execution.

The above example can be seen in `example.go` in function `twoWayDfaWtl()`
