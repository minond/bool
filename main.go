package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

const (
	// For > .mode MODE
	scanMode  = "scan"
	parseMode = "parse"
	evalMode  = "eval"

	// For > MODE: expression
	scanLine  = "scan:"
	parseLine = "parse:"
	evalLine  = "eval:"

	setMode = ".mode "

	cmdHelp     = ".help"
	cmdKeyboard = ".keyboard"
	cmdMode     = ".mode"
	cmdQuit     = ".quit"
	cmdReset    = ".reset"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	env := newEnvironment(nil)
	mode := evalMode

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		switch text {
		case cmdQuit:
			fmt.Println("< Goodbye")
			return

		case cmdMode:
			fmt.Printf("< %s mode\n", mode)

		case cmdReset:
			fmt.Println("< clearing environment\n")
			env = newEnvironment(nil)

		case cmdHelp:
			fmt.Printf("< %s: reset current environment.\n", cmdReset)
			fmt.Printf("< %s: display or change evaluation mode to %s, %s, or %s.\n", cmdMode, scanMode, parseMode, evalMode)
			fmt.Printf("< %s: print a keyboard with valid operations and their ascii representation.\n", cmdKeyboard)
			fmt.Printf("< %s: view this help text.\n", cmdHelp)
			fmt.Printf("< %s: exit program.\n", cmdQuit)

		case cmdKeyboard:
			fmt.Printf("< conjunction: %s or %s\n", string(andRn), string(andAsciiRn))
			fmt.Printf("< disjunction: %s or %s\n", string(orRn), string(orAsciiRn))
			fmt.Printf("< negation: %s or %s\n", string(notRn), string(notAsciiRn))
			fmt.Printf("< exclusive or: %s or %s\n", string(xorRn), string(xorAsciiRn))
			fmt.Printf("< equivalence: %s or %s\n", string(eqRn), string(eqAsciiRn))
			fmt.Printf("< material implication: %s or %s\n", string(miRn), string(miAsciiRn))

		default:
			if text == "" {
				continue
			} else if strings.HasPrefix(text, setMode) {
				maybeMode := strings.TrimSpace(strings.TrimPrefix(text, setMode))
				switch maybeMode {
				case scanMode:
					mode = scanMode
				case parseMode:
					mode = parseMode
				case evalMode:
					mode = evalMode
				default:
					fmt.Printf("< error: Invalid mode `%s`\n", maybeMode)
					continue
				}

				fmt.Printf("< switching to %s mode\n", mode)
			} else if strings.HasPrefix(text, ".") {
				fmt.Printf("< error: Unknown command: `%s`. Enter `.help` for help.\n", text)
			} else if mode == scanMode || strings.HasPrefix(text, scanLine) {
				for _, t := range scan(strings.TrimPrefix(text, scanLine)) {
					fmt.Printf("< %04d %s\n", t.pos, t)
				}
			} else if mode == parseMode || strings.HasPrefix(text, parseLine) {
				spew.Dump(parse(scan(strings.TrimPrefix(text, parseLine))))
			} else if mode == evalMode || strings.HasPrefix(text, evalLine) {
				// FIXME This is really ugly. This to clean up:
				//
				//   1. I don't like that there is a need for an isBinding flag
				//   that gets used at the end. This should be somehow cleaner.
				//
				//   2. A lot of the error checking and printing should be in a
				//   separate function so main doesn't get messy.
				isBinding := false
				expr, parseErrors := parse(scan(strings.TrimPrefix(text, evalLine)))

				if len(parseErrors) > 0 {
					fmt.Println("< error: Cannot parse expression due to errors:")

					for _, err := range parseErrors {
						fmt.Printf("< error: %s\n", err)
					}

					fmt.Println()
					continue
				}

				switch expr.(type) {
				case binding:
					isBinding = true
				}

				ret, evalErrors := expr.eval(env)

				if len(evalErrors) > 0 {
					fmt.Println("< error: Cannot evaluate expression due to errors:")

					for _, err := range evalErrors {
						fmt.Printf("< error: %s\n", err)
					}

					fmt.Println()
					continue
				}

				if isBinding {
					fmt.Println("< ok")
				} else {
					fmt.Printf("= %t\n", ret.internal)
				}
			}
		}

		fmt.Print("\n")
	}
}
