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

func parse(tokens []token) (evaluates, []error) {
	par := parser{
		pos:    0,
		tokens: tokens,
		errs:   []error{},
	}

	expr := par.main()
	errs := par.errs

	switch v := expr.(type) {
	case expression:
		errs = append(errs, v.errors()...)

	case binding:
		errs = append(errs, v.value.errors()...)

	case *gate:
		errs = append(errs, v.body.errors()...)
	}

	return expr, errs
}

func (p *parser) main() evaluates {
	if p.match(gateTok) {
		return p.gateDecl()
	} else if p.match(bindContTok) || (p.curr().id == identTok && p.peek().id == bindTok) {
		return p.binding()
	} else {
		expr := p.expression()

		if !p.done() {
			p.errs = append(p.errs, fmt.Errorf(
				"Unexpected word a position %d `%s`",
				p.curr().pos, p.curr().lexeme))
		}

		return expr
	}
}

func (p *parser) binding() binding {
	label := p.curr()

	if p.expect(identTok) != nil {
		p.errs = append(p.errs, errors.New("Expecting a binding label."))
		return binding{}
	}

	if p.expect(bindTok) != nil {
		p.errs = append(p.errs, errors.New("Expecting `is` keyword."))
		return binding{}
	}

	return binding{
		label: label,
		value: p.expression(),
	}
}

func (p *parser) gateDecl() *gate {
	g := &gate{}
	g.env = nil

	if p.expect(identTok) != nil {
		p.errs = append(p.errs, errors.New("Expecting a gate label."))
		return g
	}

	g.label = p.prev()

	if p.expect(oparenTok) != nil {
		p.errs = append(p.errs, errors.New(
			"Expecting an open paren after the gate label."))
		return g
	}

	if p.curr().id == identTok {
		for {
			if !p.match(identTok) {
				p.errs = append(p.errs, fmt.Errorf("Expecting an identifier "+
					"in position %d but found %s instead.",
					p.curr().pos, p.curr()))
			}

			g.args = append(g.args, p.prev())

			if p.match(identTok) {
				p.errs = append(p.errs, fmt.Errorf("Expecting a comma to "+
					"separate gate arguments. Found identity in position %d "+
					"instead.", p.prev().pos))
				return g
			}

			if !p.match(commaTok) {
				break
			}
		}
	}

	if p.expect(cparenTok) != nil {
		p.errs = append(p.errs, fmt.Errorf(
			"Expecting a close paren after the gate arguments but found %s "+
				"in position %d instead..", p.curr(), p.curr().pos))
		return g
	}

	if p.expect(eqTok) != nil {
		p.errs = append(p.errs, fmt.Errorf(
			"Expecting an equal sign after gate arguments but found %s in "+
				"position %d instead.", p.curr(), p.curr().pos))
		return g
	}

	g.body = p.expression()

	return g
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

		// unary = primary = gate-call
		if p.match(oparenTok) {
			expr.call = true

			// This is matching id(arg,) since getting to the command restarts
			// and immediatelly ends because of the !p.match(cparenTok). Maybe
			// this is ok. Maybe it's not. Just noting it here.
			for !p.match(cparenTok) {
				arg := p.expression()

				if arg.err != nil {
					expr.err = arg.err
					return expr
				}

				expr.args = append(expr.args, arg)

				if p.match(commaTok) {
					continue
				}

				if p.match(cparenTok) {
					break
				} else {
					expr.err = fmt.Errorf(
						"Expecting a closing paren but found %s in position %d instead.",
						p.curr(), p.curr().pos)
				}
			}
		}
	} else if p.match(trueTok, falseTok) {
		// unary = primary = BOOLEAN
		if p.prev().id == trueTok {
			expr.literal = &boolean{true}
		} else {
			expr.literal = &boolean{false}
		}
	} else if p.match(oparenTok) {
		// unary = "(" expression ")"
		lhs := p.expression()
		expr.lhs = &lhs
		expr.err = p.expect(cparenTok)
	} else if p.curr().id == eolTok {
		expr.err = errors.New("Unexpected end of line.")
	} else {
		expr.err = fmt.Errorf(
			"Invalid expression starting in position %d with character `%s`.",
			p.curr().pos, p.curr().lexeme)
	}

	return expr
}

func (p *parser) expect(ids ...tokenId) error {
	if !p.match(ids...) {
		return fmt.Errorf("Expecting one of the following tokens %v but found %s",
			ids, p.curr().id)
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

func (p parser) done() bool {
	return p.pos >= len(p.tokens)
}
