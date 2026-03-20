package main

import ("os"
	"fmt"
	keeper "github.com/Supraboy981322/keeper/golang"
)

func _(){keeper.Add(&[]rune{}, 0)}

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

func recurse(code []byte) []Token {
	var p = P{
		idx: -1,
		code: code,
	}
	var comment bool
	loop: for p.next() {
		if comment && p.cur != '\n' { continue loop } else if comment { comment = false }
		switch p.cur {
			case '(': if !p.esc {
				p.Toks = append(p.Toks, mktok(IGNORE, BOX, nil))
				p.Toks = append(p.Toks, recurse(p.seek_to(')'))...)
				p.idx--
			}
			case '\\': p.esc = !p.esc
			case '"': if !p.esc {
				p.Toks = append(p.Toks, mktok(VALUE, STRING, p.seek_to('"')))
			}
		  case ')': p.Toks = append(p.Toks, mktok(IGNORE, EOX, nil))
			case ';': if !p.esc { comment = true } 
			case '\n', ' ', '\r', '\t':
			default: {
				c := p.cur
				thing := append([]byte{c}, p.seek_to(' ')...)
				fmt.Println("=" + string(thing) + "=")
				go_is_dumb := func(foo TokTypeNote, bar TokType) (TokTypeNote, TokType, []byte) {
					return foo, bar, thing
				}
				p.Toks = append(p.Toks, mktok(go_is_dumb(p.match_name(thing))))
			}
		}
	}
	var invalid []byte
	for _, tok := range p.Toks { if tok.Type == TokType(INVALID) { invalid = tok.Raw } }
	if invalid != nil { 
		fmt.Printf("unexpected token: %s\n", string(invalid))
		os.Exit(1)
	}
	return p.Toks
}

func mktok(note TokTypeNote, t TokType, raw []byte) Token {
	return Token {
		Type: t,
		Raw: raw,
		Note: note,
	}
}

func (P) match_name(name []byte) (TokTypeNote, TokType) {
	switch string(name) {
		case "print": return TokTypeNote(FN), TokType(PRINT)
		case "run": return TokTypeNote(FN),  TokType(RUN)
		case "?": return TokTypeNote(OPERATOR),  TokType(IF)
		case "?!": return TokTypeNote(OPERATOR), TokType(ELSE)
		case "&": return TokTypeNote(OPERATOR), TokType(AND)
		case "|": return TokTypeNote(OPERATOR), TokType(OR)
		case ";": return TokTypeNote(NONE), TokType(COMMENT)
		default: return TokTypeNote(IGNORE), TokType(INVALID)
	}
}

func seek_toks(input *[]Token) []Token {
	var output []Token
	for (*input)[0].Type != EOX {
		keeper.Add(&output, (*input)[0])
		keeper.Shift(input)
		if len(*input) < 1 { return output }
	}
	return output
}

func eval(input []Token) {
	if len(input) < 1 { return }
	loop: {
		thing := input[0]
		switch thing.Type {
			case TokType(PRINT): {
				keeper.Shift(&input)
				builtin.Print(seek_toks(&input));
			}
			case TokType(EOX), TokType(BOX):
			case TokType(INVALID): builtin.Err_Out(
				string(append([]byte("invalid token as fn call: "), thing.Raw...)),
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

func (p *P) next() bool {
	p.cur = p.seek()
	if p.cur == 0 { return false }
	return true
}

func (p *P) peek() byte {
	if len(p.code) <= p.idx+1 { return 0 }
	return p.code[p.idx+1]
}

func (p *P) toss() { p.idx++ } 

func (p *P) tossN(n int) {
	for n > 0 || p.idx < len(p.code) {
		p.idx++  ;  n--
	}
}

func (p *P) seek() byte {
	b := p.peek()  ;  p.toss()
	return b
}

func (p *P) seekN(n int) []byte {
	s := p.peekN(n)  ;  p.tossN(n)
	return s
}

func (p *P) seek_to(c byte) []byte {
	var mem []byte
	defer fmt.Println(string(mem))
	var esc bool
	if c == '"' {
		for p.next() {
			if esc { continue }
			switch p.cur {
				case '"': if !esc { return mem } else { mem = append(mem, p.cur) }
				case '\\': esc = !esc
				default: if esc {
					mem = append(mem, builtin.Get_Esc(p.cur))
				} else {
					mem = append(mem, p.cur)
				}
			}
		}
		return mem
	}
	loop: for p.next() {
		if p.cur == '\\' {
			esc = true
			continue loop
		}
		if p.cur == c && !esc { return mem }
		//ternary wouldn've been nice here
		if esc {
			mem = append(mem, builtin.Get_Esc(p.cur))
		} else {
			mem = append(mem, p.cur)
		}
		esc = false
	}
	return nil
}

func (p *P) peekN(n int) []byte {
	if len(p.code) <= p.idx+n { return nil }
	return p.code[p.idx:][:n]
}

func unmatch_token(tok Token) string {
	switch tok.Type {
		case TokType(PRINT):   return "[PRINT]"
		case TokType(RUN):     return "[RUN]"
		case TokType(IF):      return "[IF]"
		case TokType(ELSE):    return "[ELSE]"
		case TokType(AND):     return "[AND]"
		case TokType(OR):      return "[OR]"
		case TokType(COMMENT): return "[COMMENT]"
		case TokType(INVALID): return "[INVALID]"
		case TokType(STRING):  return "[STRING]"
		case TokType(EOX):     return "[EOX]"
		case TokType(BOX):     return "[BOX]"
		default: panic("UNKNOWN TOKEN: |" + string(tok.Raw) + "|")
	}
}
