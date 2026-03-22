package main

import (
	"os"
	"fmt"
	"strconv"
	"math/big"
)

func _(){fmt.Print()}

type Builtin struct{}
var builtin Builtin

func (Builtin) Err_Out(str string) {
	os.Stderr.WriteString("\n" + str + "\n")
	os.Exit(1)
}
func (b Builtin) Err_OutF(format string, args ...any) {
	str := fmt.Sprintf(format, args...)
	b.Err_Out(str)
}

func (Builtin) Print(input []Token, caller Token) []Token {
	var f func([]byte)(int,error)
	switch string(caller.User_Note) {
		case "", "stdout", "out", "1": f = os.Stdout.Write
		case "stderr", "err", "2": f = os.Stderr.Write
	}
	for _, t := range input { f(t.Raw) }
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

//this is perfect for Zig, I hate that I have to cast the return type
func (Builtin) ToNum(str []byte, int_type IntType) (*big.Int, int64) {
	var e error
	var res int64
	switch int_type {
		case BIG: {
			res := new(big.Int)
			res, ok := res.SetString(string(str), 0)
			if !ok {
				builtin.Err_OutF("NaN: %s", string(str))
			}
			return res, 0
		}
		case U8: {
			var n uint64
			n, e = strconv.ParseUint(string(str), 0, 8)
			res = int64(n)
		}
		case I8: {
			res, e = strconv.ParseInt(string(str), 0, 8)
		}
		case U16: {
			var n uint64
			n, e = strconv.ParseUint(string(str), 0, 8)
			res = int64(n)
		}
		case I16: {
			res, e = strconv.ParseInt(string(str), 0, 16)
		}
		case U32: {
			var n uint64
			n, e = strconv.ParseUint(string(str), 0, 32)
			res = int64(n)
		}
		case I32: {
			res, e = strconv.ParseInt(string(str), 0, 32)
		}
		case I64: {
			res, e = strconv.ParseInt(string(str), 0, 64)
		}
	}
	if e != nil {
		builtin.Err_OutF("NaN: %s", string(str))
	}
	return nil, res
}
