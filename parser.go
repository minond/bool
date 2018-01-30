package main

type ast struct{}

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
expression   = binaryOp ;
binaryOp     = unaryOp { BIN_OPERATOR unaryOp } ;
unaryOp      = UNI_OPERATOR unaryOp
             | primary ;

primary      = BOOLEAN
             | "(" expression ")"

*/
func parse(toks []token) ast {
	return ast{}
}
