import unittest

from two_way_dfa_wtl import TwoWayDfaWtl


class TwoWayDfaWtlTest(unittest.TestCase):
    def setUp(self):
        self.automaton1 = TwoWayDfaWtl.from_symbolic_definition(
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

        automaton2_alphabet = {"a", "b"}
        self.automaton2 = TwoWayDfaWtl.from_symbolic_definition(
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
                "q_a": automaton2_alphabet,
                "q_b": automaton2_alphabet,
                "p_r": automaton2_alphabet
            },
            "q_0",
            "]",
            "["
        )

    def test_sweeping_detection_true(self):
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
        self.assertTrue(automaton.is_sweeping())

    def test_sweeping_detection_false_case_is_marker(self):
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
                    "[": "q_b",
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
        self.assertFalse(automaton.is_sweeping())

    def test_sweeping_detection_false_case_not_is_marker(self):
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
                    "[": "p_b",
                },
                "q_b": {
                    "[": "p_b"
                },
                "p_a": {
                    "a": "q_a"
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
        self.assertFalse(automaton.is_sweeping())

    def test_ex1_automaton_accept(self):
        computation = self.automaton1.Computation(self.automaton1, "ababc")
        actualResult: list[str] = [computation.configuration()]
        while not computation.halt():
            computation.step()
            actualResult.append(computation.configuration())

        expectedResult = [
            "]q_a ababc[",
            "]q_b babc[",
            "]abc q_c[",
            "]q_a ab[",
            "]q_b b[",
            "] q_c[",
            "accept"
        ]

        self.assertEqual(actualResult, expectedResult)

    def test_ex1_automaton_reject(self):
        computation = self.automaton1.Computation(self.automaton1, "aba")
        actualResult: list[str] = [computation.configuration()]
        while not computation.halt():
            computation.step()
            actualResult.append(computation.configuration())

        expectedResult = [
            "]q_a aba[",
            "]q_b ba[",
            "]a q_c[",
            "reject"
        ]

        self.assertEqual(actualResult, expectedResult)

    def test_ex2_automaton_accept(self):
        computation = self.automaton2.Computation(self.automaton2, "abba")
        actualResult: list[str] = [computation.configuration()]
        while not computation.halt():
            computation.step()
            actualResult.append(computation.configuration())
        expectedResult = [
            "]q_0 abba[",
            "]q_a bba[",
            "]bba p_a[",
            "]bb p_r[",
            "]q_0 bb[",
            "]q_b b[",
            "]b p_b[",
            "] p_r[",
            "]q_0 [",
            "accept"
        ]

        self.assertEqual(actualResult, expectedResult)

    def test_ex2_automaton_reject(self):
        computation = self.automaton2.Computation(self.automaton2, "abab")
        actualResult: list[str] = [computation.configuration()]
        while not computation.halt():
            computation.step()
            actualResult.append(computation.configuration())
        expectedResult = [
            "]q_0 abab[",
            "]q_a bab[",
            "]bab p_a[",
            "reject"
        ]

        self.assertEqual(actualResult, expectedResult)


if __name__ == "__main__":
    unittest.main()
