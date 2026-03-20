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
		case 'n': return '\n'
		case 'b': return '\b'
		case 'r': return '\r'
		case 'a': return '\a'
		case 't': return '\t'
		case 'v': return '\v'
		case 'f': return '\f'
		case 'e': return '\x1b'
		default: builtin.Err_Out("invalid escape")
	}
	return 0
}
