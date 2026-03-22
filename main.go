package main

import (
	"os"
	"fmt"
	"bufio"
	"slices"
	keeper "github.com/Supraboy981322/keeper/golang"
)

type TokType int
const (
	INVALID TokType = iota
	PRINT   //'print'
	RUN     //'run'
	QUIT    //'quit'
	IF      //'?'
	ELSE    //'?!'
	TRUE    //'true'
	FALSE   //'false'
	OR      //'|'
	AND     //'&'
	BOX     //'('
	EOX     //')'
	COMMENT //';'
	STRING
	NUMBER
	GREATER_THAN //'>'
	LESS_THAN    //'<'
	EQL_TO       //'='
	WHITESPACE
	VOID
)

type TokTypeNote int
const (
	NONE TokTypeNote = iota
	FN
	OPERATOR
	VALUE
	COMPARE
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
	User_Note []byte
}

func _(){fmt.Print()}

var debug_mode, repl bool

func main() {
	var code []byte
	if len(os.Args) < 2 { builtin.Err_Out("not enough args, need filename") }
	var taken []int
	next_arg := func(i int, a string) []byte {
		keeper.DrainInto(&taken, &[]int{i+1, i})
		if len(os.Args) <= i+1 { builtin.Err_OutF("provided %s arg, but no value provided", a) }
		return []byte(os.Args[i+1])
	}
	loop: for i, a := range os.Args[1:] {
		if slices.Contains(taken, i) { continue loop }
		switch a {
			case "eval":  { code = []byte(next_arg(i+1, a)) }
			case "repl":  { repl = true }
			case "debug": { debug_mode = true }
			default: {
				var e error
				code, e = os.ReadFile(os.Args[1])
				if e != nil { builtin.Err_OutF("couldn't read file: %v", e) }
			}
		}
	}
	
	if repl {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() && repl {
			code = []byte(scanner.Text())
			vague_toks := []TokType{}
			toks := run(code) 
			for _, t := range toks { keeper.Add(&vague_toks, t.Type) }
			if slices.Contains(vague_toks, QUIT) { repl = false }
		}
	} else { run(code) }
}

func run(code []byte) []Token {
	tokens := recurse(code)
	if debug_mode {
		for _, t := range tokens {
			fmt.Printf(
				"(%v) %s\n",
				unmatch_token(t),
				string(builtin.Un_Escape(t.Raw)),
			)
		}
		fmt.Println("\n\n========\n\n")
	}
	return recurse_eval(tokens, void)
}
