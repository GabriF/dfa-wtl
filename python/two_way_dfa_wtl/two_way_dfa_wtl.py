from enum import Enum, auto
import sys


class _StateHeadDirection(Enum):
    LEFT = auto()
    RIGHT = auto()


ACCEPT_STRING = "accept"
REJECT_STRING = "reject"
class TwoWayDfaWtl:
    def __init__(self,
                 states: list[str],
                 state_direction: list[_StateHeadDirection],
                 delta: list[dict[str, int]],
                 tau: list[set[str]],
                 start: int,
                 left_marker: str,
                 right_marker: str):
        """
        User MUST NOT call this constructor
        """
        self.__states = states
        self.__state_direction = state_direction
        self.__delta = delta
        self.__tau = tau
        self.__start = start
        self.__left_marker = left_marker
        self.__right_marker = right_marker

    @classmethod
    def from_symbolic_definition(cls,
                                 right_state: set[str],
                                 left_state: set[str],
                                 delta: dict[str, dict[str, str]],
                                 tau: dict[str, set[str]],
                                 start: str,
                                 left_marker: str,
                                 right_marker: str) -> TwoWayDfaWtl:
        states: list[str] = [ACCEPT_STRING] + list(left_state | right_state)

        state_direction: list[_StateHeadDirection] = [
            _StateHeadDirection.LEFT if s in left_state
            else _StateHeadDirection.RIGHT
            for s in states
        ]

        state_to_int = {s: i for i, s in enumerate(states)}

        try:
            start = state_to_int[start]
        except KeyError:
            raise ValueError("Start state does not exist")

        try:
            delta: list[dict[str, int]] = [
                {} if s not in delta.keys()
                else {letter: state_to_int[to_state] for letter, to_state in delta[s].items()}
                for s in states
            ]
        except KeyError:
            raise ValueError("Destination state does not exist")

        tau: list[set[str]] = [
            {} if s not in tau.keys()
            else {letter for letter in tau[s]}
            for s in states
        ]

        return TwoWayDfaWtl(states, state_direction, delta, tau, start, left_marker, right_marker)

    @staticmethod
    def __make_configuration(state: str,
                             direction: _StateHeadDirection,
                             inscription: str,
                             left_marker: str,
                             right_marker: str) -> str:
        if state == ACCEPT_STRING:
            return ACCEPT_STRING
        
        if direction == _StateHeadDirection.RIGHT:
            middle = state + " " + inscription
        else:
            middle = inscription + " " + state

        return left_marker + middle + right_marker

    def exec(self, w: str) -> list[str]:
        inscription = w
        current_state = self.__start
        current_configuration = self.__make_configuration(
            self.__states[current_state],
            self.__state_direction[current_state],
            inscription,
            self.__left_marker,
            self.__right_marker)
        configuration_story: list[str] = [current_configuration]

        while current_configuration != ACCEPT_STRING and current_configuration != REJECT_STRING:
            current_direction = self.__state_direction[current_state]

            if current_direction == _StateHeadDirection.RIGHT:
                to_iter = inscription + self.__right_marker
            else:
                to_iter = reversed(self.__left_marker + inscription)

            for i, letter in enumerate(to_iter):
                if letter not in self.__tau[current_state]:
                    if letter not in self.__delta[current_state]:
                        current_configuration = REJECT_STRING
                    else:
                        if current_direction == _StateHeadDirection.RIGHT:
                            inscription = inscription[:i] + inscription[i+1:]
                        else:
                            inscription = inscription[:len(
                                inscription) - i - 1] + inscription[len(inscription) - i:]

                        current_state = self.__delta[current_state][letter]
                        current_configuration = self.__make_configuration(
                            self.__states[current_state],
                            self.__state_direction[current_state],
                            inscription,
                            self.__left_marker,
                            self.__right_marker
                        )
                    break

            configuration_story.append(current_configuration)

        return configuration_story


def main():
    a = TwoWayDfaWtl.from_symbolic_definition(
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

    for s in a.exec(sys.argv[1]):
        print(s)


if __name__ == "__main__":
    main()
