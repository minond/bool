package main

import (
	"errors"
	"fmt"
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

MAIN         = { ["where"|"and"] binding | expression } ;

identifier   = LETTER , { LETTER | DIGIT | "_" } ;
value        = "true" | "false" ;

binding      = identifier "is" expression ;
expression   = unary { BIN_OPERATOR unary } ;
unary        = UNI_OPERATOR unary
             | primary ;

primary      = BOOLEAN
             | identifier

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
	if p.match(bindContTok) {
		// We have to be a binding here
		return p.binding()
	} else if p.curr().id == identTok && p.peek().id == bindTok {
		return p.binding()
	} else {
		return p.expression()
	}
}

func (p *parser) binding() binding {
	label := p.curr()

	p.expect(identTok)
	p.expect(bindTok)

	return binding{
		label: label,
		value: p.expression(),
	}
}

func (p *parser) expression() expression {
	expr := p.unary()

	for p.match(andTok, orTok, miTok, xorTok, eqTok) {
		lhs := expr
		op := cloneToken(p.prev())
		rhs := p.unary()

		expr = expression{}
		expr.lhs = &lhs
		expr.op = &op
		expr.rhs = &rhs
	}

	return expr
}

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
	} else if p.curr().id == eolTok {
		expr.err = errors.New("Unexpected end of line.")
	} else {
		expr.err = fmt.Errorf("Invalid expression starting in position %d with character `%s`.",
			p.curr().pos, p.curr().lexeme)
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
