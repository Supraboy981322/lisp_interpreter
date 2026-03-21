package main

import (
	"os"
	"fmt"
)

func _(){fmt.Print()}

type Builtin struct{}
var builtin Builtin

func (Builtin) Err_Out(str string) {
	os.Stderr.WriteString("\n" + str + "\n")
	os.Exit(1)
}
func (b Builtin) Err_OutF(format string, args ...any) {
	str := fmt.Sprintf(format, args)
	b.Err_Out(str)
}

func (Builtin) Stderr(input []Token, _ Token) []Token {
	for _, t := range input {
		os.Stderr.Write(t.Raw)
	}
	return void_return()
}
func (Builtin) Stdout(input []Token, _ Token) []Token {
	for _, t := range input {
		os.Stdout.Write(t.Raw)
	}
	return void_return()
}

func (Builtin) Get_Esc(b byte) byte {
	switch b {

		// NOTE: escape and (mostly) control characters
		case 'n': return '\n'   //newline
		case 'b': return '\b'   //backspace
		case 'r': return '\r'   //return
		case 'a': return '\a'   //bell character (never used a terminal that does anything with this)
		case 't': return '\t'   //tab
		case 'v': return '\v'   //vertical tab (I keep forgetting this one exists)
		case 'f': return '\f'   //form-feed (what's the purpose of this in modern terminals?)
		case 'e': return '\x1b' //much simpler to just have a one-char ansi escape
		case '0': return 0      //null character

		//just return the input character if unknown
		default: { return b }
	}
}

func (Builtin) Un_Escape(str []byte) []byte {
	if len(str) < 1 { return nil }
	var res []byte
	var i int; var b byte; loop: {
		b = str[i]
		switch b {
			case '\n':    res = append(res, []byte("\\n")...)
			case '\b':    res = append(res, []byte("\\b")...) 
			case '\r':    res = append(res, []byte("\\r")...)
			case '\a':    res = append(res, []byte("\\a")...)
			case '\t':    res = append(res, []byte("\\t")...)
			case '\v':    res = append(res, []byte("\\v")...)
			case '\f':    res = append(res, []byte("\\f")...)
			case '\x1b':  res = append(res, []byte("\\e")...)
			case 0:       res = append(res, []byte("\\0")...)
			default: 			res = append(res, b)
		}
		i++
		if i < len(str) { goto loop }
	}
	return res
}
