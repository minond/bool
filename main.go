package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	env := newEnvironment(nil)

	scanMode := "scan:"
	modeMode := "mode:"
	evalMode := "eval:"
	astMode := "ast:"
	currMode := evalMode

	for {
		fmt.Print("\n> ")
		text, _ := reader.ReadString('\n')

		switch strings.TrimSpace(text) {
		case "quit":
			fallthrough
		case "exit":
			fmt.Println("< Goodbye")
			return

		case "m":
			fallthrough
		case "mode":
			fmt.Printf("< %s on\n", currMode)

		case "h":
			fallthrough
		case "help":
			fmt.Printf("< mode: display or change evaluation mode to one of [%s %s %s]\n", scanMode, astMode, evalMode)
			fmt.Println("< keyboard: print a keyboard with valid operations and their ascii representation.")

		case "k":
			fallthrough
		case "keyboard":
			fmt.Printf("< conjunction: %s or %s\n", string(andRn), string(andAsciiRn))
			fmt.Printf("< disjunction: %s or %s\n", string(orRn), string(orAsciiRn))
			fmt.Printf("< negation: %s or %s\n", string(notRn), string(notAsciiRn))

		default:
			input := strings.TrimSpace(text)
			scanning := strings.HasPrefix(input, scanMode)
			modding := strings.HasPrefix(input, modeMode)
			asting := strings.HasPrefix(input, astMode)
			evaling := strings.HasPrefix(input, evalMode)

			if modding {
				req := strings.TrimSpace(strings.TrimPrefix(text, modeMode))

				switch req {
				case "scan":
					currMode = scanMode

				case "eval":
					currMode = evalMode

				case "ast":
					currMode = astMode

				default:
					fmt.Printf("< error: invalid mode %s\n", req)
					fmt.Printf("< options: [%s %s %s]\n", scanMode, astMode, evalMode)
				}

				fmt.Printf("< %s on\n", currMode)
				continue
			} else if scanning {
				input = strings.TrimSpace(strings.TrimPrefix(text, scanMode))
			} else if asting {
				input = strings.TrimSpace(strings.TrimPrefix(text, astMode))
			} else if evaling {
				input = strings.TrimSpace(strings.TrimPrefix(text, evalMode))
			}

			if scanning || currMode == scanMode {
				toks := scan(input)

				for _, t := range toks {
					fmt.Printf("< %04d %s\n", t.pos, t)
				}
			} else if asting || currMode == astMode {
				spew.Dump(parse(scan(input)))
			} else if evaling || currMode == evalMode {
				var errs []error
				ast := parse(scan(input))

				switch v := ast.(type) {
				case expression:
					errs = v.errors()

				case binding:
					errs = v.value.errors()
				}

				if len(errs) > 0 {
					for _, err := range errs {
						fmt.Printf("> error: %s\n", err)
					}
				} else {
					fmt.Printf("= %t\n", ast.eval(env).internal)
				}
			}
		}
	}
}
