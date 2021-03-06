package main

import (
	"errors"
	"fmt"
	"strconv"
)

type typeId string
type method func(env environment, args ...value) (value, []error)

type environment struct {
	bindings map[string]expression
	methods  map[string]method
	gates    map[string]gate
	parent   *environment
}

type value struct {
	boolean  *boolean
	sequence *sequence
	number   int
}

type boolean struct {
	internal bool
}

type sequence struct {
	internal []expression
}

type binding struct {
	label token
	value expression
}

type gate struct {
	label token
	args  []token
	body  expression
	env   *environment
}

// NOTE Well, I guess this is the ugly side of Go's type system. Note to my
// future self: use separate structs for all of these expressions and move
// towards something like a visitor pattern that could be automated with `go
// generate`.
//
// Kind of a
// catch-all structure for all types of expressions. Since it serves multiple
// purposes, it needs to be checked in a specific order:
//
//   - Check 1: err
//   - Check 2: lhs + op + rhs, this is a binary expression
//   - Check 3: op + rhs, this is a unary expression with an expression
//   - Check 3: lhs, this is a grouped expression
//   - Check 4: identifier + call, this is a gate call
//   - Check 5: identifier, this is a plain identifier
//   - Check 6: literal, this is a plain literal
//   - Check 7: sequence, this is a sequence
//   - Check 8: num, this is a number
type expression struct {
	err        error
	lhs        *expression
	rhs        *expression
	op         *token
	identifier *token
	num        *token
	call       bool
	args       []expression
	literal    *boolean
	sequence   *sequence
	env        *environment
}

type evaluates interface {
	eval(env environment) (value, []error)
}

const (
	typeInvalid  typeId = "invalid"
	typeBoolean  typeId = "boolean"
	typeSequence typeId = "sequence"
	typeNumber   typeId = "number"
)

func (b binding) eval(env environment) (value, []error) {
	for _, id := range b.value.identifiers(env) {
		if id.lexeme == b.label.lexeme {
			return value{}, []error{fmt.Errorf(
				"Detected circular reference in `%s` identifier",
				b.label.lexeme)}
		}
	}

	env.setBinding(b.label.lexeme, b.value)
	return value{}, nil
}

func (g *gate) eval(env environment) (value, []error) {
	env.setGate(g.label.lexeme, *g)
	return value{}, nil
}

func (b expression) eval(env environment) (value, []error) {
	if b.env != nil {
		env = *b.env
	}

	if b.err != nil {
		return value{}, []error{fmt.Errorf(
			"Cannot evaluate expression due to error: %s",
			b.err)}
	} else if b.lhs != nil && b.op != nil && b.rhs != nil {
		lhs, lhsErr := b.lhs.eval(env)
		rhs, rhsErr := b.rhs.eval(env)

		errs := append(lhsErr, rhsErr...)

		if len(errs) > 0 {
			return value{}, errs
		}

		switch b.op.id {
		case andTok:
			if fn, ok := env.getMethod("and"); !ok {
				return value{}, []error{fmt.Errorf("Unknown binary operator: %s",
					b.op.lexeme)}
			} else {
				return fn(env, lhs, rhs)
			}

		case orTok:
			if fn, ok := env.getMethod("or"); !ok {
				return value{}, []error{fmt.Errorf("Unknown binary operator: %s",
					b.op.lexeme)}
			} else {
				return fn(env, lhs, rhs)
			}

		case miTok:
			if fn, ok := env.getMethod("mi"); !ok {
				return value{}, []error{fmt.Errorf("Unknown binary operator: %s",
					b.op.lexeme)}
			} else {
				return fn(env, lhs, rhs)
			}

		case xorTok:
			if fn, ok := env.getMethod("xor"); !ok {
				return value{}, []error{fmt.Errorf("Unknown binary operator: %s",
					b.op.lexeme)}
			} else {
				return fn(env, lhs, rhs)
			}

		case eqTok:
			if fn, ok := env.getMethod("eq"); !ok {
				return value{}, []error{fmt.Errorf("Unknown binary operator: %s",
					b.op.lexeme)}
			} else {
				return fn(env, lhs, rhs)
			}

		case geTok:
			if fn, ok := env.getMethod("ge"); !ok {
				return value{}, []error{fmt.Errorf("Unknown binary operator: %s",
					b.op.lexeme)}
			} else {
				return fn(env, lhs, rhs)
			}

		case gtTok:
			if fn, ok := env.getMethod("gt"); !ok {
				return value{}, []error{fmt.Errorf("Unknown binary operator: %s",
					b.op.lexeme)}
			} else {
				return fn(env, lhs, rhs)
			}

		case leTok:
			if fn, ok := env.getMethod("le"); !ok {
				return value{}, []error{fmt.Errorf("Unknown binary operator: %s",
					b.op.lexeme)}
			} else {
				return fn(env, lhs, rhs)
			}

		case ltTok:
			if fn, ok := env.getMethod("lt"); !ok {
				return value{}, []error{fmt.Errorf("Unknown binary operator: %s",
					b.op.lexeme)}
			} else {
				return fn(env, lhs, rhs)
			}

		default:
			return value{}, []error{fmt.Errorf("Unknown binary operator: %s",
				b.op.lexeme)}
		}
	} else if b.op != nil && b.rhs != nil {
		val, errs := b.rhs.eval(env)

		if len(errs) > 0 {
			return value{}, errs
		}

		switch b.op.id {
		case notTok:
			if fn, ok := env.getMethod("not"); !ok {
				return value{}, []error{fmt.Errorf("Unknown unary operator: %s",
					b.op.lexeme)}
			} else {
				return fn(env, val)
			}

		default:
			return value{}, []error{fmt.Errorf("Unknown unary operator: %s",
				b.op.lexeme)}
		}
	} else if b.lhs != nil {
		return b.lhs.eval(env)
	} else if b.identifier != nil && b.call {
		gate, set := env.getGate(b.identifier.lexeme)

		if !set {
			if len(b.args) != 1 {
				return value{}, []error{fmt.Errorf("Undefined gate `%s`",
					b.identifier.lexeme)}
			}

			var seq value
			var errs []error
			val, set := env.getBinding(b.identifier.lexeme)

			if b.identifier != val.identifier {
				seq, errs = val.eval(env)
			} else if env.parent != nil {
				seq, errs = val.eval(*env.parent)
			}

			if len(errs) > 0 {
				return value{}, errs
			}

			idx, errs := b.args[0].eval(env)

			// NOTE Ok this is totally a hack. Clearly there's something wrong
			// with both the AST data structures and the evalution ones as
			// well.
			if idx.isBoolean() && idx.boolean.internal {
				idx.boolean = nil
				idx.number = 1
			} else if idx.isBoolean() && !idx.boolean.internal {
				idx.boolean = nil
				idx.number = 0
			}

			if !set {
				return value{}, []error{fmt.Errorf("Undefined gate `%s`",
					b.identifier.lexeme)}
			} else if seq.sequence == nil {
				return value{}, []error{fmt.Errorf("Invalid operation, expecting `%s` to be a sequence",
					b.identifier.lexeme)}
			} else if len(errs) > 0 {
				return value{}, append(errs,
					fmt.Errorf("Invalid operation, expecting a digit when accessing `%s`",
						b.identifier.lexeme))
			} else if idx.isSequence() || idx.isBoolean() {
				return value{}, []error{fmt.Errorf(
					"Expecting a digic when accessing `%s` gate.",
					b.identifier.lexeme)}
			} else if idx.number >= len(seq.sequence.internal) {
				return value{}, []error{fmt.Errorf(
					"Out of bounds error, max is %d and tried to access %d on `%s` sequence.",
					len(seq.sequence.internal)-1, idx.number, b.identifier.lexeme)}
			} else {
				return seq.sequence.internal[idx.number].eval(env)
			}
		} else {
			if len(gate.args) != len(b.args) {
				return value{}, []error{fmt.Errorf("Arity error, `%s` "+
					"expects %d arguments but got %d instead.",
					b.identifier.lexeme, len(gate.args), len(b.args))}
			}

			gate.env.parent = &env
			subEnv := newEnvironment(gate.env)
			defer func() {
				gate.env.parent = nil
			}()

			for i, arg := range gate.args {
				func(i int) {
					b.args[i].env = &env
					subEnv.setBinding(arg.lexeme, b.args[i])

					defer func() {
						b.args[i].env = nil
					}()
				}(i)
			}

			res, errs := gate.body.eval(subEnv)

			if len(errs) > 0 {
				return value{}, errs
			}

			if res.isSequence() {
				snapshop, errs := res.sequence.freeze(subEnv)
				return value{sequence: &snapshop}, errs
			} else {
				return res, errs
			}
		}
	} else if b.identifier != nil {
		val, set := env.getBinding(b.identifier.lexeme)

		if !set {
			return value{}, []error{fmt.Errorf("Undefined identifier `%s`",
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
			return value{}, []error{fmt.Errorf(
				"Internal error, detected a circular variable reference and "+
					"expected a parent environment for lookup but none found "+
					"for `%s` binding", b.identifier.lexeme)}
		}
	} else if b.literal != nil {
		return value{boolean: &boolean{b.literal.internal}}, nil
	} else if b.sequence != nil {
		return value{sequence: b.sequence}, nil
	} else if b.num != nil {
		num, err := strconv.Atoi(b.num.lexeme)

		if err != nil {
			return value{}, []error{
				fmt.Errorf("Error converting to number: %v", err)}
		} else {
			return value{number: num}, []error{}
		}
	} else {
		return value{}, []error{errors.New("Invalid evaluation path")}
	}
}

func (b boolean) eval(env environment) (value, []error) {
	return value{boolean: &boolean{b.internal}}, nil
}

func (e *environment) getBinding(label string) (expression, bool) {
	val, ok := e.bindings[label]

	if !ok && e.parent != nil {
		return e.parent.getBinding(label)
	} else {
		return val, ok
	}
}

func (e *environment) getMethod(label string) (method, bool) {
	val, ok := e.methods[label]

	if !ok && e.parent != nil {
		return e.parent.getMethod(label)
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
		methods:  getBuiltins(),
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

func (s sequence) freeze(env environment) (sequence, []error) {
	var errs []error
	snapshop := sequence{}

	for _, expr := range s.internal {
		val, err := expr.eval(env)
		errs = append(errs, err...)

		if val.isBoolean() {
			snapshop.internal = append(snapshop.internal, expression{
				literal: &boolean{
					internal: val.boolean.internal,
				},
			})
		} else if val.isSequence() {
			snapshop, err := val.sequence.freeze(env)
			errs = append(errs, err...)
			snapshop.internal = append(snapshop.internal, expression{
				sequence: &snapshop,
			})
		} else if val.isNumber() {
			snapshop.internal = append(snapshop.internal, expression{
				num: &token{
					id:     numTok,
					lexeme: strconv.FormatInt(int64(val.number), 10),
				},
			})
		}
	}

	return snapshop, errs
}

func (v value) isBoolean() bool {
	return v.boolean != nil
}

func (v value) isSequence() bool {
	return v.sequence != nil
}

func (v value) isNumber() bool {
	return !v.isBoolean() && !v.isSequence()
}

func (v value) getTypeId() typeId {
	switch {
	case v.isBoolean():
		return typeBoolean

	case v.isSequence():
		return typeSequence

	case v.isNumber():
		return typeNumber

	default:
		return typeInvalid
	}
}

func (v value) equals(other value, env environment) bool {
	if v.getTypeId() != other.getTypeId() {
		return false
	}

	switch {
	case v.isBoolean():
		return v.boolean.internal == other.boolean.internal

	case v.isNumber():
		return v.number == other.number

	case v.isSequence():
		if len(v.sequence.internal) != len(other.sequence.internal) {
			return false
		}

		for i, _ := range v.sequence.internal {
			e1, err := v.sequence.internal[i].eval(env)

			if err != nil {
				return false
			}

			e2, err := other.sequence.internal[i].eval(env)

			if err != nil {
				return false
			}

			if e1.equals(e2, env) == false {
				return false
			}
		}

		return true
	}

	return false
}

func getBuiltins() map[string]method {
	return map[string]method{
		"and": andBuiltin,
		"eq":  eqBuiltin,
		"ge":  geBuiltin,
		"gt":  gtBuiltin,
		"le":  leBuiltin,
		"lt":  ltBuiltin,
		"mi":  miBuiltin,
		"not": notBuiltin,
		"or":  orBuiltin,
		"xor": xorBuiltin,
	}
}

func andBuiltin(env environment, args ...value) (value, []error) {
	if err := strictArityCheck("and", 2, args...); err != nil {
		return value{}, []error{err}
	}

	if err := strictTypeCheck("and", 1, args[0], typeBoolean); err != nil {
		return value{}, []error{err}
	}

	if err := strictTypeCheck("and", 2, args[1], typeBoolean); err != nil {
		return value{}, []error{err}
	}

	return value{
		boolean: &boolean{
			args[0].boolean.internal && args[1].boolean.internal,
		},
	}, nil
}

func orBuiltin(env environment, args ...value) (value, []error) {
	if err := strictArityCheck("or", 2, args...); err != nil {
		return value{}, []error{err}
	}

	if err := strictTypeCheck("or", 1, args[0], typeBoolean); err != nil {
		return value{}, []error{err}
	}

	if err := strictTypeCheck("or", 2, args[1], typeBoolean); err != nil {
		return value{}, []error{err}
	}

	return value{
		boolean: &boolean{
			args[0].boolean.internal || args[1].boolean.internal,
		},
	}, nil
}

func miBuiltin(env environment, args ...value) (value, []error) {
	if err := strictArityCheck("mi", 2, args...); err != nil {
		return value{}, []error{err}
	}

	if err := strictTypeCheck("mi", 1, args[0], typeBoolean); err != nil {
		return value{}, []error{err}
	}

	if err := strictTypeCheck("mi", 2, args[1], typeBoolean); err != nil {
		return value{}, []error{err}
	}

	return value{
		boolean: &boolean{
			!args[0].boolean.internal || args[1].boolean.internal,
		},
	}, nil
}

func xorBuiltin(env environment, args ...value) (value, []error) {
	if err := strictArityCheck("xor", 2, args...); err != nil {
		return value{}, []error{err}
	}

	if err := strictTypeCheck("xor", 1, args[0], typeBoolean); err != nil {
		return value{}, []error{err}
	}

	if err := strictTypeCheck("xor", 2, args[1], typeBoolean); err != nil {
		return value{}, []error{err}
	}

	return value{
		boolean: &boolean{
			(args[0].boolean.internal || args[1].boolean.internal) &&
				!(args[0].boolean.internal && args[1].boolean.internal),
		},
	}, nil
}

func eqBuiltin(env environment, args ...value) (value, []error) {
	if err := strictArityCheck("eq", 2, args...); err != nil {
		return value{}, []error{err}
	}

	if args[0].getTypeId() != args[1].getTypeId() {
		return value{}, []error{fmt.Errorf("Type error, `eq` expects both arguments to be of the same type but got `%s` and `%s` instead.",
			args[0].getTypeId(), args[1].getTypeId())}
	}

	return value{
		boolean: &boolean{
			args[0].equals(args[1], env),
		},
	}, nil
}

func geBuiltin(env environment, args ...value) (value, []error) {
	if err := strictArityCheck("ge", 2, args...); err != nil {
		return value{}, []error{err}
	}

	if err := strictOneOfTypeCheck("ge", 1, args[0], typeBoolean, typeNumber); err != nil {
		return value{}, []error{err}
	}

	if err := strictOneOfTypeCheck("ge", 2, args[1], typeBoolean, typeNumber); err != nil {
		return value{}, []error{err}
	}

	return value{
		boolean: &boolean{
			numberCast(args[0]).number >= numberCast(args[1]).number,
		},
	}, nil
}

func gtBuiltin(env environment, args ...value) (value, []error) {
	if err := strictArityCheck("gt", 2, args...); err != nil {
		return value{}, []error{err}
	}

	if err := strictOneOfTypeCheck("gt", 1, args[0], typeBoolean, typeNumber); err != nil {
		return value{}, []error{err}
	}

	if err := strictOneOfTypeCheck("gt", 2, args[1], typeBoolean, typeNumber); err != nil {
		return value{}, []error{err}
	}

	return value{
		boolean: &boolean{
			numberCast(args[0]).number > numberCast(args[1]).number,
		},
	}, nil
}

func leBuiltin(env environment, args ...value) (value, []error) {
	if err := strictArityCheck("le", 2, args...); err != nil {
		return value{}, []error{err}
	}

	if err := strictOneOfTypeCheck("le", 1, args[0], typeBoolean, typeNumber); err != nil {
		return value{}, []error{err}
	}

	if err := strictOneOfTypeCheck("le", 2, args[1], typeBoolean, typeNumber); err != nil {
		return value{}, []error{err}
	}

	return value{
		boolean: &boolean{
			numberCast(args[0]).number <= numberCast(args[1]).number,
		},
	}, nil
}

func ltBuiltin(env environment, args ...value) (value, []error) {
	if err := strictArityCheck("lt", 2, args...); err != nil {
		return value{}, []error{err}
	}

	if err := strictOneOfTypeCheck("lt", 1, args[0], typeBoolean, typeNumber); err != nil {
		return value{}, []error{err}
	}

	if err := strictOneOfTypeCheck("lt", 2, args[1], typeBoolean, typeNumber); err != nil {
		return value{}, []error{err}
	}

	return value{
		boolean: &boolean{
			numberCast(args[0]).number < numberCast(args[1]).number,
		},
	}, nil
}

func notBuiltin(env environment, args ...value) (value, []error) {
	if err := strictArityCheck("not", 1, args...); err != nil {
		return value{}, []error{err}
	}

	if err := strictTypeCheck("not", 1, args[0], typeBoolean); err != nil {
		return value{}, []error{err}
	}

	return value{boolean: &boolean{!args[0].boolean.internal}}, []error{}
}

func strictArityCheck(label string, expected int, args ...value) error {
	if expected != len(args) {
		return fmt.Errorf("Arity error, `%s` expects %d arguments but got %d instead.",
			label, expected, len(args))
	}

	return nil
}

func strictTypeCheck(label string, pos int, got value, expId typeId) error {
	gotId := got.getTypeId()

	if expId != gotId {
		return fmt.Errorf("Type error, `%s` expects a `%s` in position %d but got `%s` instead.",
			label, expId, pos, gotId)
	}

	return nil
}

func strictOneOfTypeCheck(label string, pos int, got value, expIds ...typeId) error {
	gotId := got.getTypeId()

	for _, t := range expIds {
		if t == gotId {
			return nil
		}
	}

	return fmt.Errorf("Type error, `%s` expects one of %v in position %d but got `%s` instead.",
		label, expIds, pos, gotId)
}

// Converts booleans into number and keeps everything else the same.
func numberCast(v value) value {
	if v.getTypeId() == typeBoolean {
		num := 0

		if v.boolean.internal {
			num = 1
		}

		return value{
			number: num,
		}
	}

	return v
}
