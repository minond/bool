package main

import "errors"

type environment struct{}

type parser struct {
	pos    int
	tokens []token
	errs   []error
}

type boolean struct {
	internal bool
}

type binding struct {
	label token
	value expression
}

type expression struct {
	lhs     *expression
	rhs     *expression
	op      token
	literal boolean
}

type ast interface {
	eval(env environment) boolean
}

func (b binding) eval(env environment) boolean {
	return boolean{}
}

func (b boolean) eval(env environment) boolean {
	return boolean{}
}

func (b expression) eval(env environment) boolean {
	return boolean{}
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
expression   = unaryOp { BIN_OPERATOR unaryOp } ;
unaryOp      = UNI_OPERATOR unaryOp
             | primary ;

primary      = BOOLEAN
             | "(" expression ")"

*/
func parse(tokens []token) ast {
	par := parser{
		pos:    0,
		tokens: tokens,
		errs:   []error{},
	}

	return par.main()
}

func (p *parser) main() ast {
	if p.curr().id == identTok && p.peek().id == eqTok {
		return p.binding()
	} else {
		return p.expression()
	}
}

func (p *parser) expression() expression {
	return p.unary()
}

// XXX
func (p *parser) binding() binding {
	label := p.curr()

	p.expect(identTok)
	p.expect(eqTok)

	value := p.expression()

	return binding{
		label: label,
		value: value,
	}
}

// XXX
func (p *parser) unary() expression {
	return expression{
		literal: boolean{},
	}
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
	if p.pos-1 > len(p.tokens) {
		return token{id: eolTok}
	} else {
		return p.tokens[p.pos-1]
	}
}

func (p parser) curr() token {
	if p.pos > len(p.tokens) {
		return token{id: eolTok}
	} else {
		return p.tokens[p.pos]
	}
}

func (p parser) peek() token {
	if p.pos+1 > len(p.tokens) {
		return token{id: eolTok}
	} else {
		return p.tokens[p.pos+1]
	}
}
