package main

import (
	"fmt"
	keeper "github.com/Supraboy981322/keeper/golang"
)

func _(){fmt.Print()}

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


func (p *P) peekN(n int) []byte {
	if len(p.code) <= p.idx+n { return nil }
	return p.code[p.idx:][:n]
}

func seek_toks(input *[]Token) []Token {
	var output []Token
	var depth int
	var started bool
	for {
		if started && depth < 1 { break }
		keeper.Add(&output, (*input)[0])
		keeper.Shift(input)
		switch (*input)[0].Type {
			case BOX: { depth++ }
			case EOX: {
				if depth-1 < 1 { return output } else { depth-- }
			}
		}
	}
	return output
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
		case "stdout": return TokTypeNote(FN),       TokType(STDOUT)
		case "stderr": return TokTypeNote(FN),       TokType(STDERR)
		case "run":    return TokTypeNote(FN),       TokType(RUN)
		case "?":      return TokTypeNote(OPERATOR), TokType(IF)
		case "?!":     return TokTypeNote(OPERATOR), TokType(ELSE)
		case "&":      return TokTypeNote(OPERATOR), TokType(AND)
		case "|":      return TokTypeNote(OPERATOR), TokType(OR)
		case ";":      return TokTypeNote(NONE),     TokType(COMMENT)
		case "<":      return TokTypeNote(COMPARE),  TokType(LESS_THAN)
		case ">":      return TokTypeNote(COMPARE),  TokType(GREATER_THAN)
		case "=":      return TokTypeNote(COMPARE),  TokType(EQL_TO)
		default:       return TokTypeNote(IGNORE),   TokType(INVALID)
	}
}

func unmatch_token(tok Token) string {
	switch tok.Type {
		case TokType(STDOUT):       return "[STDOUT]"
		case TokType(STDERR):       return "[STDERR]"
		case TokType(RUN):          return "[RUN]"
		case TokType(IF):           return "[IF]"
		case TokType(ELSE):         return "[ELSE]"
		case TokType(AND):          return "[AND]"
		case TokType(OR):           return "[OR]"
		case TokType(COMMENT):      return "[COMMENT]"
		case TokType(INVALID):      return "[INVALID]"
		case TokType(STRING):       return "[STRING]"
		case TokType(EOX):          return "[EOX]"
		case TokType(BOX):          return "[BOX]"
		case TokType(GREATER_THAN): return "[GREATER_THAN]"
		case TokType(LESS_THAN):    return "[LESS_THAN]"
		case TokType(EQL_TO):       return "[EQL_TO]"
		case TokType(NUMBER):       return "[NUMBER]"
		default:
			panic("UNKNOWN TOKEN: |" + string(tok.Raw) + "|")
	}
}

func (p *P) collapse_str() []byte {
	var mem []byte
	var esc bool
	for p.next() {
		if p.cur == '"' && !esc { return mem }
		if esc {
			esc = false
			if p.cur == '"'  || p.cur == '\\' {
				keeper.Add(&mem, p.cur)
			} else {
				keeper.Add(&mem, builtin.Get_Esc(p.cur))
			}
		} else {
			if p.cur == '\\' {
				esc = true
			} else {
				keeper.Add(&mem, p.cur)
			}
		}
	}
	return mem
}

// TODO: track depth to return full list 
func (p *P) seek_to(c byte) []byte {
	var mem []byte
	var esc bool
	if c == '"' {
		return p.collapse_str()
	}
	for p.next() {
		switch p.cur {
			case '\\': if esc {
				mem = append(mem, p.cur)
			}; esc = !esc

			case '"': if esc {
				mem = append(mem, p.cur)
			} else {
				mem = append(mem, p.collapse_str()...)
				p.toss()
			}
			
			case c: if esc {
				mem = append(mem, p.cur)
			} else {
				return mem
			}

			default: mem = append(mem, p.cur)
		}
	}
	return nil
}

func (p *P) previous() byte {
	return p.code[p.idx-1]
}

func (p *P) back() byte {
	p.cur = p.previous()
	return p.cur
}

func (p *P) seek_whitespace() []byte {
	var mem []byte
	var esc bool
	defer p.back()
	_ = esc
	for p.next() {
		switch p.cur {
			case ' ', '\n', '\t', '\r': return mem
			default: mem = append(mem, p.cur)
		}
	}

	return nil
}

func (p *P) seek_num() []byte {
	str := []byte{p.cur}
	for p.next() {
		if p.cur >= '0' && p.cur <= '9' {
			keeper.Add(&str, p.cur)
		} else {
			return str
		}
	}
	return nil
}
