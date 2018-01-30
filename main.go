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
		case ".quit":
			fmt.Println("< Goodbye")
			return

		case ".mode":
			fmt.Printf("< %s mode\n", mode)

		case ".help":
			fmt.Printf("< .mode: display or change evaluation mode to %s, %s, or %s.\n", scanMode, parseMode, evalMode)
			fmt.Println("< .keyboard: print a keyboard with valid operations and their ascii representation.")
			fmt.Println("< .help: view this help text.")
			fmt.Println("< .quit: exit program.")

		case ".keyboard":
			fmt.Printf("< conjunction: %s or %s\n", string(andRn), string(andAsciiRn))
			fmt.Printf("< disjunction: %s or %s\n", string(orRn), string(orAsciiRn))
			fmt.Printf("< negation: %s or %s\n", string(notRn), string(notAsciiRn))

		default:
			if strings.HasPrefix(text, setMode) {
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
				var errs []error
				ast := parse(scan(strings.TrimPrefix(text, evalLine)))

				switch v := ast.(type) {
				case expression:
					errs = v.errors()

				case binding:
					errs = v.value.errors()
				}

				if len(errs) > 0 {
					for _, err := range errs {
						fmt.Printf("< error: %s\n", err)
					}
				} else {
					fmt.Printf("= %t\n", ast.eval(env).internal)
				}
			}
		}

		fmt.Print("\n")
	}
}
