package main

import (
	"errors"
	"fmt"
)

type environment struct {
	bindings map[string]expression
	gates    map[string]gate
	parent   *environment
}

type boolean struct {
	internal bool
}

type binding struct {
	label token
	value expression
}

type gate struct {
	label token
	args  []token
	body  expression
}

// Kind of a catch-all structure for all types of expressions. Since it serves
// multiple purposes, it needs to be checked in a specific order:
//   - Check 1: err
//   - Check 2: lhs + op + rhs, this is a binary expression
//   - Check 3: op + rhs, this is a unary expression with an expression
//   - Check 3: lhs, this is a grouped expression
//   - Check 4: identifier + call, this is a gate call
//   - Check 5: identifier, this is a plain identifier
//   - Check 6: literal, this is a plain literal
type expression struct {
	err        error
	lhs        *expression
	rhs        *expression
	op         *token
	identifier *token
	call       bool
	args       []expression
	literal    *boolean
}

type evaluates interface {
	eval(env environment) (boolean, []error)
}

func (b binding) eval(env environment) (boolean, []error) {
	for _, id := range b.value.identifiers(env) {
		if id.lexeme == b.label.lexeme {
			return boolean{}, []error{fmt.Errorf("Detected circular reference in `%s` identifier",
				b.label.lexeme)}
		}
	}

	env.setBinding(b.label.lexeme, b.value)
	return boolean{}, nil
}

func (g gate) eval(env environment) (boolean, []error) {
	env.setGate(g.label.lexeme, g)
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
	} else if b.identifier != nil && b.call {
		gate, set := env.getGate(b.identifier.lexeme)

		if !set {
			return boolean{}, []error{fmt.Errorf("Undefined gate `%s`",
				b.identifier.lexeme)}
		} else {
			if len(gate.args) != len(b.args) {
				return boolean{}, []error{fmt.Errorf("Arity error, `%s` expects %d arguments but got %d instead.",
					b.identifier.lexeme, len(gate.args), len(b.args))}
			}

			subEnv := newEnvironment(&env)

			for i, arg := range gate.args {
				subEnv.setBinding(arg.lexeme, b.args[i])
			}

			return gate.body.eval(subEnv)
		}
	} else if b.identifier != nil {
		val, set := env.getBinding(b.identifier.lexeme)

		if !set {
			return boolean{}, []error{fmt.Errorf("Undefined identifier `%s`",
				b.identifier.lexeme)}
		} else if b.identifier != val.identifier {
			return val.eval(env)
		} else if env.parent != nil {
			// Ok, we're looking up a binding but when doing so we get back the
			// same identifier as the one we looked up, so we're probably
			// looking something up with the same name. If we keep looking for
			// it using the current environment, we'll get stuck. Let's try
			// looking it up in the parent environment instead. Note that this
			// should always have a parent environment since we're most likely
			// in a gate right now.
			return val.eval(*env.parent)
		} else {
			return boolean{}, []error{fmt.Errorf(
				"Internal error, detected a circular variable reference and "+
					"expected a parent environment for lookup but none found "+
					"for `%s` binding", b.identifier.lexeme)}
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

func (e *environment) getBinding(label string) (expression, bool) {
	val, ok := e.bindings[label]

	if !ok && e.parent != nil {
		return e.parent.getBinding(label)
	} else {
		return val, ok
	}
}

func (e *environment) setBinding(label string, expr expression) *environment {
	e.bindings[label] = expr
	return e
}

func (e *environment) getGate(label string) (gate, bool) {
	val, ok := e.gates[label]

	if !ok && e.parent != nil {
		return e.parent.getGate(label)
	} else {
		return val, ok
	}
}

func (e *environment) setGate(label string, g gate) *environment {
	e.gates[label] = g
	return e
}

func newEnvironment(parent *environment) environment {
	return environment{
		bindings: make(map[string]expression),
		gates:    make(map[string]gate),
		parent:   parent,
	}
}

func (e expression) identifiers(env environment) []token {
	var tokens []token

	if e.identifier != nil {
		tokens = append(tokens, *e.identifier)
		ident, ok := env.getBinding(e.identifier.lexeme)

		if ok {
			tokens = append(tokens, ident.identifiers(env)...)
		}
	}

	if e.lhs != nil {
		tokens = append(tokens, e.lhs.identifiers(env)...)
	}

	if e.rhs != nil {
		tokens = append(tokens, e.rhs.identifiers(env)...)
	}

	return tokens
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
