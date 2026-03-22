package main

import (
	"os"
	"fmt"
	"strconv"
	keeper "github.com/Supraboy981322/keeper/golang"
)

func eval(input []Token) []Token {
	//ignore empty input
	if len(input) < 1 { return []Token{} }
	if debug { fmt.Println("eval") }

	shiftN := func(n int) {
		loop: {
			keeper.Shift(&input)
			n--
			if n > 0 { goto loop }
		}
	}

	//local helper to seek to EOX 
	mk_args := func(in *[]Token) []Token {
		toks := seek_toks(in)
		shiftN(len(toks))
		keeper.Shift(&toks)
		return toks
	}

	//local helper to string-together the raw strings
	//  of each arg from seek to EOX 
	string_args := func() []byte {
		var str []byte
		for _ , t := range mk_args(&input) {
			keeper.Add(&str, append(t.Raw, ' ')...)
		}
		return str
	}

	drain_args := func() []Token {
		args := mk_args(&input)
		shiftN(len(args))
		return args
	}

	var mem []Token
	//defer func() {
	//	for _, t := range mem { t.print() }
	//}()

	call := func(f func([]Token, Token)[]Token, caller Token) {
		keeper.DrainInto(&mem, keeper.PtrOf(f(drain_args(), caller))) 
	}

	//who needs a 'for' loop anyways?
	loop: {
		//fmt.Println(len(input))
		//for _, t := range input { t.print() }
		thing := input[0]

		if thing.Note == TokTypeNote(VALUE) && thing.Type != TokType(VOID) {
			keeper.DrainInto(&mem, keeper.PtrOf(drain_args()))
			goto skip
		}

		switch thing.Type {
			case TokType(VOID): if len(input) <= 1 { return mem } else { goto skip }

			case TokType(QUIT): return append(mem, thing)

			//builtin functions
			case TokType(PRINT): { call(builtin.Print, thing) }
			case TokType(RUN):   {
				keeper.DrainInto(&mem, keeper.PtrOf(recurse_eval(recurse(string_args()), thing)))
			}

			//comparison
			case TokType(GREATER_THAN), TokType(LESS_THAN), TokType(EQL_TO): {
				call(compare, thing)
			}

			//return on EOX
			case TokType(EOX): return mem
			//recurse on BOX
			case TokType(BOX): { call(recurse_eval, void) }

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
							unmatch_token(thing), builtin.Un_Escape(thing.Raw), len(thing.Raw),
						),
					)
				}

				// TODO: functions
				fmt.Printf(
					"\n\x1b[1;31munknown fn name\x1b[0m\n" +
					"\x1b[35m(debug: %#v) (unmatched: %s) (type note: %v)\x1b[0m:\n" +
					"\t\x1b[33m->\x1b[0m  |%s|\n",
					thing, unmatch_token(thing), thing.Note, builtin.Un_Escape(thing.Raw),
				)
				os.Exit(1)
			}
		}
		//keeper.Shift(&input)
		skip: {
			//filter-out void results from mem
			keeper.FilterFunc(
				func(tok Token) bool {
					return tok.Type != VOID
				},
				&mem,
			)
			if 0 < len(input) { goto loop }
		}
	}
	if debug { fmt.Println("eval done") }
	return mem
}

func recurse_eval(input []Token, _ Token) []Token {
	if len(input) < 1 { return void_return() }
	if input[0].Type == BOX { keeper.Shift(&input) }
	return eval(input)
}

func compare(args []Token, how Token) []Token {
	nums := []int{}
	for _, a := range args {
		if a.Type != NUMBER { builtin.Err_Out("NaN: " + string(a.Raw)) }
		n, _ := strconv.Atoi(string(a.Raw))
		keeper.Add(&nums, n)
	}

	//I WISH I HAD A TERNARY
	if len(nums) > 2 {
		builtin.Err_Out(
			fmt.Sprintf("too many args for comparison: %d (%v)", len(nums), nums),
		)
	} else if len(nums) < 2 {
		builtin.Err_Out(
			fmt.Sprintf("not enough args for comparison: %d (%v)", len(nums), nums),
		)
	}

	//these could've been reduced to ternaries
	switch how.Type {
		case GREATER_THAN: if nums[0] > nums[1] { return []Token{True} } else { return []Token{False} }
		case LESS_THAN: if nums[0] < nums[1] { return []Token{True} } else { return []Token{False} }
		case EQL_TO: if nums[0] == nums[1] { return []Token{True} } else { return []Token{False} }
	}

	return []Token{void}
}
