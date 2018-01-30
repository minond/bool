package main

import (
	"errors"
	"fmt"
)

type environment struct {
	bindings map[string]expression
	parent   *environment
}

type boolean struct {
	internal bool
}

type binding struct {
	label token
	value expression
}

// Kind of a catch-all structure for all types of expressions. Since it serves
// multiple purposes, it needs to be checked in a specific order:
//   - Check 1: err
//   - Check 2: lhs + op + rhs, this is a binary expression
//   - Check 3: op + rhs, this is a unary expression with an expression
//   - Check 3: lhs, this is a grouped expression
//   - Check 4: identifier, this is a plain identifier
//   - Check 5: literal, this is a plain literal
type expression struct {
	err        error
	lhs        *expression
	rhs        *expression
	op         *token
	identifier *token
	literal    *boolean
}

type evaluates interface {
	eval(env environment) (boolean, []error)
}

func (b binding) eval(env environment) (boolean, []error) {
	env.set(b.label.lexeme, b.value)
	return boolean{}, nil
}

func (b expression) eval(env environment) (boolean, []error) {
	if b.err != nil {
		return boolean{}, []error{fmt.Errorf("Cannot evaluate expression due to error: %s", b.err)}
	} else if b.lhs != nil && b.op != nil && b.rhs != nil {
		lhs, lhsErr := b.lhs.eval(env)
		rhs, rhsErr := b.rhs.eval(env)

		errs := append(lhsErr, rhsErr...)

		if len(errs) > 0 {
			return boolean{}, errs
		}

		switch b.op.id {
		case andTok:
			return boolean{lhs.internal && rhs.internal}, nil

		case orTok:
			return boolean{lhs.internal || rhs.internal}, nil

		case miTok:
			return boolean{!lhs.internal || rhs.internal}, nil

		case xorTok:
			return boolean{
				(lhs.internal || rhs.internal) && !(lhs.internal && rhs.internal),
			}, nil

		case eqTok:
			return boolean{lhs.internal == rhs.internal}, nil

		default:
			return boolean{}, []error{fmt.Errorf("Unknown unary operator: %s", b.op.lexeme)}
		}
	} else if b.op != nil && b.rhs != nil {
		val, errs := b.rhs.eval(env)

		if len(errs) > 0 {
			return boolean{}, errs
		}

		switch b.op.id {
		case notTok:
			return boolean{!val.internal}, nil

		default:
			return boolean{}, []error{fmt.Errorf("Unknown unary operator: %s", b.op.lexeme)}
		}
	} else if b.lhs != nil {
		return b.lhs.eval(env)
	} else if b.identifier != nil {
		val, set := env.get(b.identifier.lexeme)

		if !set {
			return boolean{}, []error{fmt.Errorf("Undefined identifier `%s`",
				b.identifier.lexeme)}
		} else {
			return val.eval(env)
		}
	} else if b.literal != nil {
		return *b.literal, nil
	} else {
		return boolean{}, []error{errors.New("Invalid evaluation path")}
	}
}

func (b boolean) eval(env environment) (boolean, []error) {
	return b, nil
}

func (e *environment) get(label string) (expression, bool) {
	val, ok := e.bindings[label]
	return val, ok
}

func (e *environment) set(label string, expr expression) expression {
	e.bindings[label] = expr
	return expr
}

func newEnvironment(parent *environment) environment {
	return environment{
		bindings: make(map[string]expression),
		parent:   parent,
	}
}

func (e expression) errors() []error {
	var errs []error

	if e.lhs != nil {
		errs = append(errs, e.lhs.errors()...)
	}

	if e.rhs != nil {
		errs = append(errs, e.rhs.errors()...)
	}

	if e.err != nil {
		errs = append(errs, e.err)
	}

	return errs
}
