package main

import "fmt"

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
