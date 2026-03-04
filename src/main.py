from two_way_dfa_wtl import TwoWayDfaWtl
import sys


def main():
    usage = '''Usage:
    python main.py <1 | 2> <word>
    '''

    try:
        n = int(sys.argv[1])
        word = sys.argv[2]
    except ValueError, IndexError:
        print(usage)
        exit(1)

    if n == 1:
        automaton = TwoWayDfaWtl.from_symbolic_definition(
            {"q_a", "q_b"},
            {"q_c"},
            {
                "q_a": {
                    "a": "q_b"
                },
                "q_b": {
                    "b": "q_c"
                },
                "q_c": {
                    "c": "q_a",
                    "]": "accept"
                }
            },
            {
                "q_a": {"b"},
                "q_b": {"a"}
            },
            "q_a",
            "]",
            "["
        )
    elif n == 2:
        alphabet = {"a", "b"}
        automaton = TwoWayDfaWtl.from_symbolic_definition(
            {"q_0", "q_a", "q_b"},
            {"p_a", "p_b", "p_r"},
            {
                "q_0": {
                    "[": "accept",
                    "a": "q_a",
                    "b": "q_b"
                },
                "q_a": {
                    "[": "p_a",
                },
                "q_b": {
                    "[": "p_b"
                },
                "p_a": {
                    "a": "p_r"
                },
                "p_b": {
                    "b": "p_r"
                },
                "p_r": {
                    "]": "q_0"
                }
            },
            {
                "q_a": alphabet,
                "q_b": alphabet,
                "p_r": alphabet
            },
            "q_0",
            "]",
            "["
        )
    else:
        print(usage)
        exit(1)

    computation = automaton.Computation(automaton, word)
    while not computation.halt():
        print(computation.configuration())
        computation.step()
    else:
        print(computation.configuration())


if __name__ == "__main__":
    main()
