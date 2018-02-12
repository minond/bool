Bool is a domain-specific reactive language and environment for Boolean Algebra
and Logic Gate Programming. Download and install the `bool` binary using `go
get github.com/minond/bool`. While a language, right now Bool mostly lives in a
repl and doesn't support external source. But that may be added in the future
(see TODO's section in this readme). Here's an example of it's reactive nature:

```text
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

```text
$ bool
> .help
< .reset: reset current environment.
< .mode: display or change evaluation mode to scan, parse, or eval.
< .keyboard: print a keyboard with valid operations and their ascii representation.
< .paste: toggle paste mode.
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
$ bool
> .paste
< paste mode: on

gate Adder (a, b, c) = [sum, carry]
  where s_ab is a ⊕ b
    and c_ab is a ∧ b
    and c_ac is a ∧ c
    and c_bc is b ∧ c
    and carry is c_ab ∨ c_ac ∨ c_bc
    and sum is c ⊕ s_ab

gate Add8 (x, y) = sum
  where b07 is Adder(x(7), y(7), 0)
    and b06 is Adder(x(6), y(6), b07(1))
    and b05 is Adder(x(5), y(5), b06(1))
    and b04 is Adder(x(4), y(4), b05(1))
    and b03 is Adder(x(3), y(3), b04(1))
    and b02 is Adder(x(2), y(2), b03(1))
    and b01 is Adder(x(1), y(1), b02(1))
    and b00 is Adder(x(0), y(0), b01(1))
    and sum is [b00(0), b01(0), b02(0), b03(0), b04(0), b05(0), b06(0), b07(0)]

.paste
< paste mode: off
> Add8([0, 0, 0, 0, 0, 0, 0, 1], [0, 0, 0, 1, 0, 0, 0, 1])
= Seq[8]{0, 0, 0, 1, 0, 0, 1, 0}

> Add8([0, 0, 0, 0, 0, 0, 1, 1], [0, 0, 0, 1, 0, 0, 0, 1])
= Seq[8]{0, 0, 0, 1, 0, 1, 0, 0}

> Add8([0, 0, 0, 0, 0, 0, 1, 1], [0, 0, 0, 1, 1, 0, 0, 1])
= Seq[8]{0, 0, 0, 1, 1, 1, 0, 0}
```

Arrays are called Sequences in Bool and work similarly to how they do in most
other languages. Accessing specific items in a sequence is done using
parentheses and is zero based, where zero is the most significant bit.

```ebnf
program        = { statement };
statement      = binding
               | gate-decl
               | expression ;

gate-decl      = "gate" identifier "(" [ gate-decl-args ] ")" "=" expression ;
gate-decl-args = identifier { "," identifier } ;
gate-call      = identifier "(" [ gate-call-args ] ")" ;
gate-call-args = expression { "," expression } ;

binding        = [ "where" | "and" ] identifier "is" expression ;
expression     = unary { BIN_OPERATOR unary } ;
unary          = [ UNI_OPERATOR ] unary
               | primary ;

primary        = BOOLEAN
               | identifier
               | number
               | gate-call
               | "(" expression ")"
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
- Ummm, tests.
