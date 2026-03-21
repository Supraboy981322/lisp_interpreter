package main

import (
	"os"
	"fmt"
	keeper "github.com/Supraboy981322/keeper/golang"
)

func eval(input []Token) {
	if len(input) < 1 { return }
	mk_args := func() []Token {
		return nil
	}
	_ = mk_args
	loop: {
		thing := input[0]
		switch thing.Type {
			case TokType(STDOUT): {
				keeper.Shift(&input)
				builtin.Stdout(seek_toks(&input));
			}
			case TokType(STDERR): {
				keeper.Shift(&input)
				builtin.Stderr(seek_toks(&input));
			}
			case TokType(RUN): {
				keeper.Shift(&input)
				toks := seek_toks(&input)
				var str []byte
				for _, t := range toks {
					keeper.Add(&str, append(t.Raw, ' ')...)
				}
				eval(recurse(str))
			}
			case TokType(EOX), TokType(BOX):
			case TokType(INVALID): builtin.Err_Out(
				"invalid token as fn call: |" + string(thing.Raw) + "|",
			)
			default: {
				if thing.Note != FN {
					builtin.Err_Out(
						fmt.Sprintf(
							"expected executable, but got a %s |%s| (len: %d)",
							unmatch_token(thing), string(thing.Raw), len(thing.Raw),
						),
					)
				}
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
