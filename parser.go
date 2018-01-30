package main

import (
	"errors"
)

type parser struct {
	pos    int
	tokens []token
	errs   []error
}

/*

      =============
      == Grammar ==
      =============

BIN_OPERATOR = (* binary operators *)
UNI_OPERATOR = (* unary operators *)
LETTER       = 'a' | .. | 'z'
DIGIT        = '0' | .. | '9'
BOOLEAN      = "true" | "false" | "1" | "0"

MAIN         = { expression | binding } ;

identifier   = LETTER , { LETTER | DIGIT | "_" } ;
value        = "true" | "false" ;

binding      = identifier "=" expression ;
expression   = unary { BIN_OPERATOR unary } ;
unary        = UNI_OPERATOR unary
             | primary ;

primary      = BOOLEAN
             | identifier
             | "(" expression ")"

*/
func parse(tokens []token) evaluates {
	par := parser{
		pos:    0,
		tokens: tokens,
		errs:   []error{},
	}

	return par.main()
}

func (p *parser) main() evaluates {
	if p.curr().id == identTok && p.peek().id == eqTok {
		return p.binding()
	} else {
		return p.expression()
	}
}

func (p *parser) binding() binding {
	label := p.curr()

	p.expect(identTok)
	p.expect(eqTok)

	return binding{
		label: label,
		value: p.expression(),
	}
}

// XXX keep checking for binary expression
func (p *parser) expression() expression {
	return p.unary()
}

// XXX get real unary expression
func (p *parser) unary() expression {
	expr := expression{}

	if p.match(notTok) {
		// unary = UNI_OPERATOR unary
		tok := cloneToken(p.prev())
		rhs := p.unary()
		expr.op = &tok
		expr.rhs = &rhs
	} else if p.match(identTok) {
		// unary = primary = identifier
		tok := cloneToken(p.prev())
		expr.identifier = &tok
	} else if p.match(trueTok, falseTok) {
		// unary = primary = BOOLEAN
		if p.prev().id == trueTok {
			expr.literal = &boolean{true}
		} else {
			expr.literal = &boolean{false}
		}
	} else {
		expr.err = errors.New("Invalid expression.")
	}

	return expr
}

func (p *parser) expect(ids ...tokenId) error {
	if !p.match(ids...) {
		return errors.New("err")
	}

	return nil
}

func (p *parser) match(ids ...tokenId) bool {
	for _, id := range ids {
		if p.curr().id == id {
			p.eat()
			return true
		}
	}

	return false
}

func (p *parser) eat() {
	p.pos += 1
}

func (p parser) prev() token {
	if p.pos-1 >= len(p.tokens) {
		return token{id: eolTok}
	} else {
		return p.tokens[p.pos-1]
	}
}

func (p parser) curr() token {
	if p.pos >= len(p.tokens) {
		return token{id: eolTok}
	} else {
		return p.tokens[p.pos]
	}
}

func (p parser) peek() token {
	if p.pos+1 >= len(p.tokens) {
		return token{id: eolTok}
	} else {
		return p.tokens[p.pos+1]
	}
}
