from enum import Enum, auto
import sys


class _StateHeadDirection(Enum):
    LEFT = auto()
    RIGHT = auto()


ACCEPT_STRING = "accept"
REJECT_STRING = "reject"


class TwoWayDfaWtl:
    """Represents a Two-Way Deterministic Finite Automaton with Translucent 
    Input Letters (2DFAwtl). Autodetect whether it meets sweeping restriction
    (SDFAwtl) conditions."""
    class Computation:
        """
        Compute a given word with the given automaton, iterating over the
        configurations which the automaton goes through.
        """

        def __init__(self, automaton: TwoWayDfaWtl, word: str):
            """Create the computation given an automaton and a word."""
            self.__automaton = automaton
            self.__current_inscription = word
            self.__current_state = automaton._start
            self.__current_direction = automaton._state_direction[self.__current_state]
            self.__accept = False
            self.__reject = False

        def halt(self) -> bool:
            """Return whether the computation has halted."""
            return self.__accept or self.__reject

        def configuration(self) -> str:
            """Return che current configuration of the computation. The string
            "accept" means the computation has halted and the word is accepted;
            the string "reject" means the computation has halted and the word is
            rejected."""
            if self.__accept:
                return ACCEPT_STRING
            elif self.__reject:
                return REJECT_STRING

            state_str = self.__automaton._states[self.__current_state]

            if self.__current_direction == _StateHeadDirection.RIGHT:
                middle = state_str + " " + self.__current_inscription
            else:
                middle = self.__current_inscription + " " + state_str

            return self.__automaton._left_marker + middle + self.__automaton._right_marker

        def step(self):
            """Compute a step. This updates the current configuration and halt
            state. Raise RuntimeError if the computation is already halted."""
            if self.halt():
                raise RuntimeError("Computation already halted.")

            if self.__current_direction == _StateHeadDirection.RIGHT:
                to_iter = self.__current_inscription + self.__automaton._right_marker
            else:
                to_iter = reversed(
                    self.__automaton._left_marker + self.__current_inscription)

            for i, letter in enumerate(to_iter):
                if letter in self.__automaton._tau[self.__current_state]:
                    continue

                if letter not in self.__automaton._delta[self.__current_state]:
                    self.__reject = True
                    break

                if letter not in {self.__automaton._left_marker, self.__automaton._right_marker}:
                    if self.__current_direction == _StateHeadDirection.RIGHT:
                        remove_index = i
                    else:
                        remove_index = len(self.__current_inscription) - i - 1
                    self.__current_inscription = self.__current_inscription[
                        :remove_index] + self.__current_inscription[remove_index+1:]

                self.__current_state = self.__automaton._delta[self.__current_state][letter]
                self.__current_direction = self.__automaton._state_direction[self.__current_state]
                if self.__automaton._states[self.__current_state] == ACCEPT_STRING:
                    self.__accept = True

                break

    def __init__(self,
                 states: list[str],
                 state_direction: list[_StateHeadDirection],
                 delta: list[dict[str, int]],
                 tau: list[set[str]],
                 start: int,
                 left_marker: str,
                 right_marker: str):
        """Create a 2DFAwtl. User MUST NOT call this method; use factory method
        from_symbolic_definition() instaed."""
        self._states = states
        self._state_direction = state_direction
        self._delta = delta
        self._tau = tau
        self._start = start
        self._left_marker = left_marker
        self._right_marker = right_marker

        self.__sweeping = True
        if state_direction[start] != _StateHeadDirection.RIGHT:
            self.__sweeping = False
        else:
            for from_state, d in enumerate(delta):
                if not self.__sweeping:
                    break
                for letter, to_state in d.items():
                    if states[to_state] == ACCEPT_STRING:
                        continue
                    from_state_direction = state_direction[from_state]
                    to_state_direction = state_direction[to_state]
                    is_marker = letter in {left_marker, right_marker}

                    if is_marker and from_state_direction == to_state_direction:
                        self.__sweeping = False
                        break
                    elif (not is_marker) and from_state_direction != to_state_direction:
                        self.__sweeping = False
                        break

    def is_sweeping(self) -> bool:
        return self.__sweeping

    @classmethod
    def from_symbolic_definition(cls,
                                 right_states: set[str],
                                 left_states: set[str],
                                 delta: dict[str, dict[str, str]],
                                 tau: dict[str, set[str]],
                                 start: str,
                                 left_marker: str,
                                 right_marker: str) -> TwoWayDfaWtl:
        """
        Generate a 2DFAwtl with the given symbolic definition.

        The special state "accept" MUST NOT be in right_states or left_states.
        User may use "accept" in delta.

        Raise ValueError if start state is not contained in right_states or
        left_states, or delta contains a mapping from a state not contained in
        right_states or left_states.

        Arguments:

        right_states -- set of states leading the input head rightward (namely,
        the ones in Q_r).

        left_states -- set of states leading the input head leftward (namely,
        the ones in Q_l).

        delta -- representation of the delta function of the automaton. This is
        a dictionary which associates one state to a dictionary which, in turn,
        associates the previous state and the current letter to the next state
        (i.e. delta[q][l] = p iff delta(q, l) = p).

        tau -- representation of the translucency mapping. This is a dictionary
        which associates one state to a set of their invisible input letters.

        start -- initial state.

        left_marker -- left marker symbol.

        right_marker -- right marker symbol. 
        """
        states: list[str] = [ACCEPT_STRING] + list(left_states | right_states)

        state_direction: list[_StateHeadDirection] = [
            _StateHeadDirection.LEFT if s in left_states
            else _StateHeadDirection.RIGHT
            for s in states
        ]

        state_to_int = {s: i for i, s in enumerate(states)}

        try:
            start = state_to_int[start]
        except KeyError:
            raise ValueError("Start state does not exist.")

        try:
            delta: list[dict[str, int]] = [
                {} if s not in delta.keys()
                else {letter: state_to_int[to_state] for letter, to_state in delta[s].items()}
                for s in states
            ]
        except KeyError:
            raise ValueError("Destination state does not exist.")

        tau: list[set[str]] = [
            {} if s not in tau.keys()
            else {letter for letter in tau[s]}
            for s in states
        ]

        return TwoWayDfaWtl(states, state_direction, delta, tau, start, left_marker, right_marker)
