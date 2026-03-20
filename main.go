package main

type TokType int
const (
	INVALID TokType = iota
	PRINT
	RUN
	IF
	ELSE
	OR
	AND
	STRING
	EOX
	BOX
	COMMENT
)

type TokTypeNote int
const (
	NONE TokTypeNote = iota
	FN
	OPERATOR
	VALUE
	IGNORE
)

type P struct{
	idx int
	cur byte
	esc bool
	code []byte
	Toks []Token
}
type Token struct {
	Raw []byte
	Type TokType
	Note TokTypeNote
}
func main() {
	code := []byte(`
(print "\e[33mcar\bt\e[0m =\t=\nfoo\rb\vfoo\a\n\"")
; (run "(print \"bar\")")
(print "\nfoo")
`)
	eval(recurse(code))
}
