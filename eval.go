package main

import (
	"os"
	"fmt"
	keeper "github.com/Supraboy981322/keeper/golang"
)

func eval(input []Token) {
	//ignore empty input
	if len(input) < 1 { return }

	//local helper to seek to EOX 
	mk_args := func() []Token {
		keeper.Shift(&input)
		return seek_toks(&input)
	}

	//local helper to string-together the raw strings
	//  of each arg from seek to EOX 
	string_args := func() []byte {
		var str []byte
		for _ , t := range mk_args() {
			keeper.Add(&str, append(t.Raw, ' ')...)
		}
		return str
	}

	//who needs a 'for' loop anyways?
	loop: {
		thing := input[0]
		switch thing.Type {

			//builtin functions
			case TokType(STDOUT): { builtin.Stdout(mk_args()) }
			case TokType(STDERR): { builtin.Stderr(mk_args()) }
			case TokType(RUN):    { eval(recurse(string_args())) }

			//ignore EOX and BOX
			case TokType(EOX), TokType(BOX):

			//err on invalid tokens  TODO: there's probably a better way to handle this
			case TokType(INVALID): builtin.Err_Out(
				"invalid token as fn call: |" + string(thing.Raw) + "|",
			)

			default: {
				//err if not a function
				if thing.Note != FN {
					builtin.Err_Out(
						fmt.Sprintf(
							"expected executable, but got a %s |%s| (len: %d)",
							unmatch_token(thing), string(thing.Raw), len(thing.Raw),
						),
					)
				}

				// TODO: functions
				fmt.Printf(
					"\n\x1b[1;31munknown fn name\x1b[0m\n" +
					"\x1b[35m(debug: %#v) (unmatched: %s) (type note: %v)\x1b[0m:\n" +
					"\t\x1b[33m->\x1b[0m  |%s|\n",
					thing, unmatch_token(thing), thing.Note, thing.Raw,
				)
				os.Exit(1)
			}
		}
		keeper.Shift(&input)
		if 0 < len(input) { goto loop }
	}
}
