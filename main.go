package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

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
			return

		case "m":
			fallthrough
		case "mode":
			fmt.Printf("  %s on\n", currMode)

		case "k":
			fallthrough
		case "keyboard":
			fmt.Println("")
			fmt.Println("-----------------------------------")
			fmt.Println("Operation | Character | Name")
			fmt.Println("-----------------------------------")
			fmt.Printf("%5s     | %5s     | Conjunction\n", string(andRn), string(andAsciiRn))
			fmt.Printf("%5s     | %5s     | Disjunction\n", string(orRn), string(orAsciiRn))
			fmt.Printf("%5s     | %5s     | Negation\n", string(notRn), string(notAsciiRn))

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
					fmt.Printf("  error: invalid mode %s\n", req)
					fmt.Printf("  options: [%s %s %s]\n", scanMode, astMode, evalMode)
				}

				fmt.Printf("  %s on\n", currMode)
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
			} else if evaling || currMode == evalMode {
				fmt.Println("= ?")
			}
		}
	}
}
