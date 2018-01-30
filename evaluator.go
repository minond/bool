package main

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

type evaluates interface {
	eval(env environment) boolean
}

func (b binding) eval(env environment) boolean {
	env.set(b.label.lexeme, b.value)
	return b.value.eval(env)
}

func (b expression) eval(env environment) boolean {
	if b.literal != nil {
		return *b.literal
	} else if b.identifier != nil {
		val, _ := env.get(b.identifier.lexeme)
		return val.eval(env)
	}

	return boolean{}
}

func (b boolean) eval(env environment) boolean {
	return b
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
