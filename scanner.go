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
	andTok      tokenId = "and"
	bindContTok tokenId = "where"
	bindTok     tokenId = "is"
	cbrakTok    tokenId = "cbrak"
	commaTok    tokenId = "comma"
	cparenTok   tokenId = "cparen"
	eolTok      tokenId = "eol"
	eqTok       tokenId = "eq"
	errTok      tokenId = "err"
	falseTok    tokenId = "false"
	gateTok     tokenId = "gate"
	identTok    tokenId = "id"
	invldTok    tokenId = "invalid"
	miTok       tokenId = "matimp"
	notTok      tokenId = "not"
	numTok      tokenId = "num"
	obrakTok    tokenId = "obrak"
	oparenTok   tokenId = "oparen"
	orTok       tokenId = "or"
	trueTok     tokenId = "true"
	xorTok      tokenId = "xor"

	andAsciiRn = rune('^')
	andRn      = rune('∧')
	cbrakRn    = rune(']')
	commaRn    = rune(',')
	cparenRn   = rune(')')
	eqAsciiRn  = rune('=')
	eqRn       = rune('≡')
	miRn       = rune('→')
	nlRn       = rune('\n')
	notAsciiRn = rune('!')
	notRn      = rune('¬')
	obrakRn    = rune('[')
	oparenRn   = rune('(')
	orAsciiRn  = rune('v')
	orRn       = rune('∨')
	spaceRn    = rune(' ')
	xorAsciiRn = rune('*')
	xorRn      = rune('⊕')

	no0 = rune('0')
	no9 = rune('9')
)

var (
	tokenDict = map[rune]tokenId{
		andAsciiRn: andTok,
		andRn:      andTok,
		eqAsciiRn:  eqTok,
		eqRn:       eqTok,
		miRn:       miTok,
		notAsciiRn: notTok,
		notRn:      notTok,
		orAsciiRn:  orTok,
		orRn:       orTok,
		xorAsciiRn: xorTok,
		xorRn:      xorTok,
	}

	opDict = map[string]tokenId{
		"not": notTok,
	}

	keywordDict = map[string]tokenId{
		"and":   bindContTok,
		"gate":  gateTok,
		"is":    bindTok,
		"where": bindContTok,
	}

	boolDict = map[string]tokenId{
		"0":     falseTok,
		"1":     trueTok,
		"false": falseTok,
		"true":  trueTok,
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

	case eqTok:
		str = "EQ"

	case orTok:
		str = "OR"

	case xorTok:
		str = "XOR"

	case obrakTok:
		str = "OPEN-BRAKET"

	case cbrakTok:
		str = "CLOSE-BRAKET"

	case oparenTok:
		str = "OPEN-PAREN"

	case cparenTok:
		str = "CLOSE-PAREN"

	case miTok:
		str = "MATERIAL-IMPLICATION"

	case commaTok:
		str = "COMMA"

	case identTok:
		str = fmt.Sprintf("ID(%s)", t.lexeme)

	case numTok:
		str = fmt.Sprintf("NUM(%s)", t.lexeme)

	case bindContTok:
		str = "WHERE"

	case gateTok:
		str = "GATE"

	case bindTok:
		str = "BIND"

	case falseTok:
		str = "FALSE"

	case trueTok:
		str = "TRUE"

	case eolTok:
		str = "EOL"
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
		var n rune
		r := runes[i]

		if i+1 < max {
			n = runes[i+1]
		}

		if isWhitespace(r) {
			continue
		} else if isOp(r) && ((r == orAsciiRn && n == rune(' ')) || r != orAsciiRn) {
			add(getOpToken(r), string(r), nil)
		} else if r == oparenRn {
			add(oparenTok, "(", nil)
		} else if r == cparenRn {
			add(cparenTok, ")", nil)
		} else if r == obrakRn {
			add(obrakTok, "[", nil)
		} else if r == cbrakRn {
			add(cbrakTok, "]", nil)
		} else if r == commaRn {
			add(commaTok, ",", nil)
		} else if isDigit(r) {
			word := readWhile(runes, i, isDigit)
			str := string(word)

			if stringIsBoolean(str) {
				add(getBoolToken(str), str, nil)
			} else {
				add(numTok, str, nil)
			}

			i += len(word) - 1
		} else {
			ident := readWhile(runes, i, isIdentLike)
			str := string(ident)

			if stringIsKeyword(str) {
				add(getKeywordToken(str), str, nil)
			} else if stringIsOp(str) {
				add(getStrOpToken(str), str, nil)
			} else if stringIsBoolean(str) {
				add(getBoolToken(str), str, nil)
			} else {
				add(identTok, str, nil)
			}

			i += len(ident) - 1
		}
	}

	return tokens
}

func cloneToken(tok token) token {
	return token{
		id:     tok.id,
		lexeme: tok.lexeme,
		pos:    tok.pos,
		err:    tok.err,
	}
}

func getOpToken(r rune) tokenId {
	tok, ok := tokenDict[r]

	if !ok {
		return invldTok
	} else {
		return tok
	}
}

func getStrOpToken(s string) tokenId {
	tok, ok := opDict[s]

	if !ok {
		return invldTok
	} else {
		return tok
	}
}

func getKeywordToken(s string) tokenId {
	tok, ok := keywordDict[s]

	if !ok {
		return invldTok
	} else {
		return tok
	}
}

func getBoolToken(s string) tokenId {
	tok, ok := boolDict[s]

	if !ok {
		return invldTok
	} else {
		return tok
	}
}

func stringIsKeyword(s string) bool {
	_, ok := keywordDict[s]
	return ok
}

func stringIsOp(s string) bool {
	_, ok := opDict[s]
	return ok
}

func stringIsBoolean(s string) bool {
	_, ok := boolDict[s]
	return ok
}

func isOp(r rune) bool {
	_, ok := tokenDict[r]
	return ok
}

func isWhitespace(r rune) bool {
	return r == spaceRn ||
		r == nlRn
}

func isDigit(r rune) bool {
	return r >= no0 && r <= no9
}

func isIdentLike(r rune) bool {
	return r == orAsciiRn || (r != commaRn &&
		r != oparenRn &&
		r != cparenRn &&
		r != obrakRn &&
		r != cbrakRn &&
		!isWhitespace(r) &&
		!isOp(r))
}

func not(f tokenFn) tokenFn {
	return func(r rune) bool {
		return !f(r)
	}
}

func or(fs ...tokenFn) tokenFn {
	return func(r rune) bool {
		for _, f := range fs {
			if f(r) {
				return true
			}
		}

		return false
	}
}

func is(rs ...rune) tokenFn {
	return func(r rune) bool {
		for _, q := range rs {
			if r == q {
				return true
			}
		}

		return false
	}
}

func readWhile(runes []rune, pos int, f tokenFn) []rune {
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
