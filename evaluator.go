package main

type environment struct {
	bindings map[string]boolean
	parent   *environment
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
	err        error
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
	if b.literal != nil {
		return *b.literal
	} else if b.identifier != nil {
		val, _ := env.get(b.identifier.lexeme)
		return val
	}

	return boolean{}
}

func (b boolean) eval(env environment) boolean {
	return b
}

func (e *environment) get(label string) (boolean, bool) {
	val, ok := e.bindings[label]
	return val, ok
}

func (e *environment) set(label string, value boolean) boolean {
	e.bindings[label] = value
	return value
}

func newEnvironment(parent *environment) environment {
	return environment{
		bindings: make(map[string]boolean),
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
