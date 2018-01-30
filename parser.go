package main

import "errors"

type environment struct {
	bindings map[string]boolean
	parent   *environment
}

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

// Check 1: lhs + op + rhs
// Check 2: op + rhs
// Check 3: op + identifier
// Check 4: op + literal
// Check 5: identifier
// Check 6: literal
type expression struct {
	lhs        *expression
	rhs        *expression
	op         *token
	identifier *token
	literal    *boolean
}

type ast interface {
	eval(env environment) boolean
}

func (b binding) eval(env environment) boolean {
	value := b.value.eval(env)
	env.set(b.label.lexeme, value)
	return value
}

func (b expression) eval(env environment) boolean {
	return boolean{}
}

func (b boolean) eval(env environment) boolean {
	return b
}

// XXX
func (e *environment) get(label string) {}

// XXX
func (e *environment) set(label string, value boolean) {}

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
	return expression{
		literal: &boolean{},
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
