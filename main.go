package main

import (
	"os"
	"fmt"
)

type TokType int
const (
	INVALID TokType = iota
	STDOUT  //'stdout'
	STDERR  //'stderr'
	RUN     //'run'
	IF      //'?'
	ELSE    //'?!'
	TRUE    //'true'
	FALSE   //'false'
	OR      //'|'
	AND     //'&'
	STRING
	BOX     //'('
	EOX     //')'
	COMMENT //';'
	WHITESPACE
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

func _(){fmt.Print()}

func main() {
	if len(os.Args) < 2 { builtin.Err_Out("not enough args, need filename") }
	code, e := os.ReadFile(os.Args[1])
	if e != nil { panic(e) }

	//for _, t := range recurse(code) {
	//	fmt.Println(string(t.Raw))
	//}
	eval(recurse(code))
}
