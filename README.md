Bool is a reactive domain-specific language and environment for Boolean
Algebra. Download and install the `bool` binary using `go get
github.com/minond/bool`. While a language, right now bool mostly lives in a
repl and doesn't support external source. But that may be added in the future
(see TODO's section in this readme). Here's an example of it's reactive nature:

```
$ bool
> x is y ^ ¬z
< ok

> x
< error: Cannot evaluate expression due to errors:
< error: Undefined identifier `y`
< error: Undefined identifier `z`

> where y is true
< ok

> and z is false
< ok

> x
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

## Language grammar

```ebnf
grammar      = { statement };
statement    = [ "where" | "and" ] binding | expression ;

binding      = identifier "is" expression ;
expression   = unary { BIN_OPERATOR unary } ;
unary        = [ UNI_OPERATOR ] unary
             | primary ;

primary      = BOOLEAN
             | identifier
             | "(" expression ")" ;

identifier   = LETTER , { LETTER | DIGIT | "_" } ;

BIN_OPERATOR = "^" | "∧" | "=" | "≡" | ">" | "→" | "v" | "∨" | "*" | "⊕" ;
UNI_OPERATOR = "¬" | "!" | "not" ;
LETTER       = "a" | .. | "z" ;
DIGIT        = "0" | .. | "9" ;
BOOLEAN      = "true" | "false" | "1" | "0" ;
```

## TODO

- Read files and stdin as source.
- Add feature to print truth tables.
- Print prettier ASTs.
