Bool is a domain-specific reactive language and environment for Boolean Algebra
and Logic Gate Programming. Download and install the `bool` binary using `go
get github.com/minond/bool`. While a language, right now bool mostly lives in a
repl and doesn't support external source. But that may be added in the future
(see TODO's section in this readme). Here's an example of it's reactive nature:

```
$ bool
> x is y ∧ ¬z
> x
< error: Cannot evaluate expression due to errors:
< error: Undefined identifier `y`
< error: Undefined identifier `z`

> y is true
> z is false
> x
= true

> gate Xor (x, y) = (x ∨ y) ∧ ¬(x ∧ y)
> Xor(true, true)
= false

> Xor(true, false)
= true

> .quit
< Goodbye
```

## Usage

The repl supports command, all of which start with a period followed by the
command name. Here's the output of running `.help`:

```
$ bool
> .help
< .mode: display or change evaluation mode to scan, parse, or eval.
< .keyboard: print a keyboard with valid operations and their ascii representation.
< .help: view this help text.
< .quit: exit program.
```

The repl also supports different modes, "eval" being on by default. The other
modes are "parse" which displays an AST, and "scan" which displays all tokens
before. These are mostly useful for debugging the parser and interpreter but
are cool nonetheless.

## Language

The language is pretty straightforward with some minor exceptions: functions
are called gates, variable binding is done with an "is" keyword sice "=" is
used for equality, and gate bodies are limited to a single expression _but_ you
can extend the scope values of a gate right after declaring it.

```text
> x is true
> y is ¬x
> y
= false

> x is false
> y
= true
```

Here we bind `x` to `true` and `y` to the inverse of `x`, `false`. We can later
update `x` to be `false` and re-evaluate `y` to see the correct answer of
`true`.

```text
gate Mux (a, b, x) = oa ∨ ob
  where nx is ¬x
    and oa is a ∧ nx
    and ob is b ∧ x
```

Here we create a gate called `Mux` which is technically just the evaluation of
`oa ∨ ob`. `oa` and `ob` are expressions that we bind using the "where" and
"and" keywords, or binding continuations, thus making them private to `Mux`.
Binding continuations outside of gate declarations result in an error.

```text
gate Adder (a, b, c) = [sum, carry]
  where s_ab is a ⊕ b
    and c_ab is a ∧ b
    and c_ac is a ∧ c
    and c_bc is b ∧ c
    and carry is c_ab ∨ c_ac ∨ c_bc
    and sum is c ⊕ s_ab
```

Arrays are called Sequences in bool and work similarly to how they do in most
other languages.

```ebnf
program        = { statement };
statement      = binding
               | gate-decl
               | expression ;

gate-decl      = "gate" identifier "(" [ gate-decl-args ] ")" "=" expression ;
gate-decl-args = identifier { "," identifier } ;
gate-call      = identifier "(" [ gate-call-args ] ")" ;
gate-call-args = expression { "," expression } ;

seq-decl       = "[" [ expression { "," expression } ] "]" ;
seq-grab       = ( seq-decl | identifier | gate-call ) "[" DIGIT "]" ;

binding        = [ "where" | "and" ] identifier "is" expression ;
expression     = unary { BIN_OPERATOR unary } ;
unary          = [ UNI_OPERATOR ] unary
               | primary ;

primary        = BOOLEAN
               | identifier
               | number
               | gate-call
               | seq-decl
               | seq-grab
               | "(" expression ")" ;
               | "[" [ expression { "," expression } ] "]" ;

number         = { DIGIT } ;
identifier     = LETTER , { LETTER | DIGIT | "_" } ;

BIN_OPERATOR   = "^" | "∧" | "=" | "≡" | ">" | "→" | "v" | "∨" | "*" | "⊕" ;
UNI_OPERATOR   = "¬" | "!" | "not" ;
LETTER         = "a" | .. | "z" ;
DIGIT          = "0" | .. | "9" ;
BOOLEAN        = "true" | "false" | "1" | "0" ;
```

## TODO

- Read files and stdin as source.
- Add feature to print truth tables.
- Print prettier ASTs.
- Move logic from main into a proper runtime.
- Arrays.
- Ummm, tests.
