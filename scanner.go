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
	eqTok       tokenId = "EQ"
	errTok      tokenId = "err"
	falseTok    tokenId = "false"
	identTok    tokenId = "id"
	invldTok    tokenId = "invalid"
	miTok       tokenId = "matimp"
	notTok      tokenId = "not"
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
	miAsciiRn  = rune('>')
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

	ucAch = rune('A')
	ucZch = rune('Z')
	lcAch = rune('a')
	lcZch = rune('z')
	chUs  = rune('_')
	no0   = rune('0')
	no9   = rune('9')
)

var (
	tokenDict = map[rune]tokenId{
		andAsciiRn: andTok,
		andRn:      andTok,
		eqAsciiRn:  eqTok,
		eqRn:       eqTok,
		miAsciiRn:  miTok,
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
		"is":    bindTok,
		"where": bindContTok,
		"and":   bindContTok,
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

	case bindContTok:
		str = "WHERE"

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
		r := runes[i]

		if isWhitespace(r) {
			continue
		} else if isOp(r) {
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
		} else if isIdent(r) {
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
		} else {
			word := readWhile(runes, i, isDigit)
			str := string(word)

			if stringIsBoolean(str) {
				add(getBoolToken(str), str, nil)
			} else {
				add(errTok, str, fmt.Errorf("unknown word: `%s`", str))
			}

			i += len(word) - 1
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

func isIdent(r rune) bool {
	return (r >= ucAch && r <= ucZch) ||
		(r >= lcAch && r <= lcZch)
}

func isDigit(r rune) bool {
	return r >= no0 && r <= no9
}

func isIdentLike(r rune) bool {
	return r == chUs || isDigit(r) || isIdent(r)
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
