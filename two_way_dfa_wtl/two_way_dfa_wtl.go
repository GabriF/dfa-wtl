package two_way_dfa_wtl

import (
	"dfawtl/utils"
	"strings"
)

type stateHeadDirection uint

const (
	R stateHeadDirection = iota
	L
)

type state int

type TwoWayDFAwtl struct {
	start       state
	direction   map[state]stateHeadDirection
	delta       map[state]map[rune]state
	translucent map[state]map[rune]bool
	stateName   map[state]string
	leftMarker  rune
	rightMarker rune
}

func makeConfiguration(state string, direction stateHeadDirection, inscription string, leftMarker rune, rightMarker rune) string {
	builder := strings.Builder{}

	builder.WriteRune(leftMarker)

	if direction == R {
		builder.WriteString(state + " " + inscription)
	} else {
		builder.WriteString(inscription + " " + state)
	}

	builder.WriteRune(rightMarker)

	return builder.String()
}

func (a *TwoWayDFAwtl) Exec(w string) []string {
	currentInscription := w
	currentState := a.start
	currentConfiguration := makeConfiguration(a.stateName[currentState], a.direction[currentState], w, a.leftMarker, a.rightMarker)
	configurationStory := []string{currentConfiguration}

	for currentConfiguration != "accept" && currentConfiguration != "reject" {
		var toRead string
		if a.direction[currentState] == R {
			toRead = currentInscription + string(a.rightMarker)
		} else {
			toRead = utils.Reverse(currentInscription) + string(a.leftMarker)
		}

		for i, l := range toRead {
			if !a.translucent[currentState][l] {
				if nextState, ok := a.delta[currentState][l]; !ok {
					currentConfiguration = "reject"
				} else if nextState == -1 {
					currentConfiguration = "accept"
				} else {
					removeIndex := i
					if a.direction[currentState] == L {
						removeIndex = len(currentInscription) - i - 1
					}
					currentInscription = currentInscription[:removeIndex] + currentInscription[removeIndex+1:]
					currentState = nextState
					currentConfiguration = makeConfiguration(a.stateName[currentState], a.direction[currentState], currentInscription, a.leftMarker, a.rightMarker)
				}
				break
			}
		}

		configurationStory = append(configurationStory, currentConfiguration)
	}

	return configurationStory
}
