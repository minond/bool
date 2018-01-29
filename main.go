package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type tokenId string
type tokenFn func(rune) bool

type token struct {
	id     tokenId
	lexeme string
	pos    int
	err    error
}

const (
	errTok   tokenId = "err"
	invldTok tokenId = "invalid"
	andTok   tokenId = "and"
	notTok   tokenId = "not"
	orTok    tokenId = "or"
	identTok tokenId = "id"
	eqTok    tokenId = "eq"

	andAsciiRn = rune('^')
	andRn      = rune('∧')
	eqRn       = rune('=')
	nlRn       = rune('\n')
	notRn      = rune('¬')
	notAsciiRn = rune('!')
	orAsciiRn  = rune('v')
	orRn       = rune('∨')
	spaceRn    = rune(' ')
)

var (
	tokenDict = map[rune]tokenId{
		andAsciiRn: andTok,
		andRn:      andTok,
		eqRn:       eqTok,
		notRn:      notTok,
		notAsciiRn: notTok,
		orAsciiRn:  orTok,
		orRn:       orTok,
	}
)

func (t token) String() string {
	if t.err != nil {
		return fmt.Sprintf("ERROR(%s)", t.err)
	}

	str := ""

	switch t.id {
	case invldTok:
		str = fmt.Sprintf("INVALID(%s)", t.lexeme)

	case andTok:
		str = "AND"

	case notTok:
		str = "NOT"

	case orTok:
		str = "OR"

	case identTok:
		str = fmt.Sprintf("ID(%s)", t.lexeme)

	case eqTok:
		str = "EQ"
	}

	return fmt.Sprintf("%s", str)
}

func main() {
	repl()
}

func repl() {
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

func scan(raw string) []token {
	var tokens []token

	runes := []rune(raw)
	max := len(runes)
	i := 0

	add := func(id tokenId, lexeme string, err error) {
		tokens = append(tokens, token{
			id:     id,
			lexeme: lexeme,
			pos:    i,
			err:    err,
		})
	}

	for ; i < max; i++ {
		r := runes[i]

		if isWhitespace(r) {
			continue
		} else if isOp(r) {
			add(getOpToke(r), string(r), nil)
		} else if isIdent(r) {
			ident := readUntil(runes, i, isIdent)
			add(identTok, string(ident), nil)
			i += len(ident) - 1
		} else {
			word := readUntil(runes, i, not(isWhitespace))
			add(errTok, string(word), fmt.Errorf("unknown word: `%s`", string(word)))
			i += len(word) - 1
		}
	}

	return tokens
}

func getOpToke(r rune) tokenId {
	tok, ok := tokenDict[r]

	if !ok {
		return invldTok
	} else {
		return tok
	}
}

func isOp(r rune) bool {
	_, ok := tokenDict[r]
	return ok
}

func isWhitespace(r rune) bool {
	return r == spaceRn ||
		r == nlRn
}

func isIdent(r rune) bool {
	return r >= rune('a') && r <= rune('z')
}

func not(f tokenFn) tokenFn {
	return func(r rune) bool {
		return !f(r)
	}
}

func readUntil(runes []rune, pos int, f tokenFn) []rune {
	var buff []rune
	max := len(runes)

	for ; pos < max; pos++ {
		if f(runes[pos]) {
			buff = append(buff, runes[pos])
		} else {
			break
		}
	}

	return buff
}
