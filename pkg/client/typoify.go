package client

import (
	"bytes"
	"math/rand"
	"unicode"
)

// C = caps lock, E = enter, S = shift, D = delete
var keyboard = [][]rune{
	{'`', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '-', '=', 'D', 'D', 'D'},
	{'\t', 'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p', '[', ']', '\\', '\\'},
	{'C', 'C', 'a', 's', 'd', 'f', 'g', 'h', 'i', 'j', 'k', 'l', ';', '\'', 'E', 'E'},
	{'S', 'S', 'S', 'z', 'x', 'c', 'v', 'b', 'n', 'm', ',', '.', '/', 'S', 'S', 'S'},
	{0, 0, 0, 0, 0, ' ', ' ', ' ', ' ', ' ', ' ', 0, 0, 0, 0, 0},
}

// ignore rune = \0
var shiftKeyboard = [][]rune{
	{'~', '!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '_', '+', 0, 0, 0},
	{0, 'Q', 'W', 'E', 'R', 'T', 'Y', 'U', 'I', 'O', 'P', '{', '}', '|', '|'},
	{0, 0, 'A', 'S', 'D', 'F', 'G', 'H', 'I', 'J', 'K', 'L', ':', '"', 0, 0},
	{0, 0, 0, 'Z', 'X', 'C', 'V', 'B', 'N', 'M', '<', '>', '?', 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}

type keyboardPos struct {
	I int
	J int
}

var keyboardLookup map[rune][]keyboardPos
var shiftKeyboardLookup map[rune][]keyboardPos

func init() {
	keyboardLookup = make(map[rune][]keyboardPos)

	for i, arr := range keyboard {
		for j, rne := range arr {
			res, fnd := keyboardLookup[rne]
			if !fnd {
				res = make([]keyboardPos, 0, 1)
			}
			keyboardLookup[rne] = append(res, keyboardPos{I: i, J: j})
		}
	}

	shiftKeyboardLookup = make(map[rune][]keyboardPos)

	for i, arr := range shiftKeyboard {
		for j, rne := range arr {
			res, fnd := shiftKeyboardLookup[rne]
			if !fnd {
				res = make([]keyboardPos, 0, 1)
			}
			shiftKeyboardLookup[rne] = append(res, keyboardPos{I: i, J: j})
		}
	}
}

// Typoify inserts typos into the given message and returns the
// message with typos. The odds of typos depends on the message
// but is multiplied by the typo chanceFactor. The resulting message
// may have been split up with "enter"
func Typoify(msg string, chanceFactor float64) []string {
	chanceFactor = 1 / chanceFactor // easier to handle the inverse
	wantShiftPressed := false

	toPressStack := make([]interface{}, 0, len(msg))
	for _, key := range msg {
		if shiftPos, shiftFound := shiftKeyboardLookup[key]; shiftFound {
			if !wantShiftPressed {
				wantShiftPressed = true
				shiftKey := keyboardLookup['S']
				toPressStack = append(toPressStack, shiftKey[rand.Intn(len(shiftKey))])
			}

			toPressStack = append(toPressStack, shiftPos[rand.Intn(len(shiftPos))])
		} else if regPos, regFound := keyboardLookup[key]; regFound {
			if wantShiftPressed {
				wantShiftPressed = false
				shiftKey := keyboardLookup['S']
				toPressStack = append(toPressStack, shiftKey[rand.Intn(len(shiftKey))])
			}

			toPressStack = append(toPressStack, regPos[rand.Intn(len(regPos))])
		} else {
			toPressStack = append(toPressStack, key)
		}
	}

	for i := 0; i < len(toPressStack)/2; i++ {
		tmp := toPressStack[i]
		toPressStack[i] = toPressStack[len(toPressStack)-i-1]
		toPressStack[len(toPressStack)-i-1] = tmp
	}

	result := make([]string, 0, 1)
	line := bytes.NewBuffer(make([]byte, 0, len(msg)+5))

	// shift and caps lock can be handled identically at this point due
	// to preprocessing the extra shift presses. the capslock is always
	// unintentional. note that this isn't *perfect* but... good enough
	shiftPressed := false
	capsLock := false

	for len(toPressStack) > 0 {
		ele := toPressStack[len(toPressStack)-1]

		// 0.1% of the time, accidentally press the key an extra time
		if rand.Float64()*chanceFactor > 0.001 {
			toPressStack = toPressStack[:len(toPressStack)-1]
		}

		// 0.1% of the time, swap this with the next element
		if rand.Float64()*chanceFactor < 0.001 && len(toPressStack) > 0 {
			tmp := toPressStack[len(toPressStack)-1]
			toPressStack[len(toPressStack)-1] = ele
			ele = tmp
		}

		// 0.1% of the time, eat the key
		if rand.Float64()*chanceFactor < 0.001 {
			continue
		}

		switch v := ele.(type) {
		case rune:
			line.WriteRune(v)
		case keyboardPos:
			// 2% of the time jitter the key pressed
			if rand.Float64()*chanceFactor < 0.02 {
				seed := rand.Float64()
				if seed < 0.25 && v.I > 0 {
					v = keyboardPos{I: v.I - 1, J: v.J}
				} else if seed < 0.5 && v.J > 0 {
					v = keyboardPos{I: v.I, J: v.J - 1}
				} else if seed < 0.75 && v.I < len(keyboard)-1 {
					v = keyboardPos{I: v.I + 1, J: v.J}
				} else if v.J < len(keyboard[v.I])-1 {
					v = keyboardPos{I: v.I, J: v.J + 1}
				}
			}

			// check if we hit a special key first
			genKey := keyboard[v.I][v.J]
			handled := true
			switch genKey {
			case 'S':
				shiftPressed = !shiftPressed
			case 'C':
				capsLock = !capsLock
			case 'D':
				if line.Len() > 0 {
					line.Truncate(line.Len() - 1)
				}
			case 'E':
				if line.Len() > 0 {
					result = append(result, line.String())
					line.Truncate(0)
				}
			default:
				handled = false
			}
			if handled {
				break
			}

			if !shiftPressed {
				if capsLock {
					genKey = unicode.ToUpper(genKey)
				}
				if genKey != 0 {
					line.WriteRune(genKey)
				}
			} else {
				shiftKey := shiftKeyboard[v.I][v.J]
				if capsLock {
					shiftKey = unicode.ToLower(shiftKey)
				}
				if shiftKey != 0 {
					line.WriteRune(shiftKey)
				}
			}
		}
	}

	if line.Len() > 0 {
		result = append(result, line.String())
	}

	return result
}
