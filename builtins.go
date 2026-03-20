package main

import (
	"os"
)

type Builtin struct{}
var builtin Builtin

func (Builtin) Err_Out(str string) {
	os.Stderr.WriteString(str)
	os.Exit(1)
}

func (Builtin) Print(str []byte) {
	var res []byte
	var i int; var b byte; loop: {
		b = str[i]
		switch b {
			case '\\': {
				if len(str) <= i+1 { builtin.Err_Out("unexpected end of string") }
				res = append(res, builtin.get_control_char(str[i+1]))
				i++
			}
			default: { res = append(res, b) }
		}
		i++
		if i < len(str) {	goto loop }
	}
	if res == nil { return }
	os.Stdout.Write(res)
}

func (Builtin) get_control_char(b byte) byte {
	switch b {

		case 'n': return '\n'   //newline
		case 'b': return '\b'   //backspace
		case 'r': return '\r'   //return
		case 'a': return '\a'   //bell character (never used a terminal that does anything with this)
		case 't': return '\t'   //tab
		case 'v': return '\v'   //vertical tab (I keep forgetting this one exists)
		case 'f': return '\f'   //form-feed (what's the purpose of this in modern terminals?)
		case 'e': return '\x1b' //much simpler to just have a one-char ansi escape
		case '0': return 0      //null character

		//anything else isn't accepted as valid
		default: builtin.Err_Out("invalid escape")
	}

	return 0 //shouldn't happen, default errs out
}
