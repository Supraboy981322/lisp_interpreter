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
	var output, mem []Token
	var recursing bool
	for len(*input) > 0 {
		thing := (*input)[0]
		keeper.Shift(input)
		switch thing.Type {
			case BOX: { recursing = true }

			case EOX: if recursing {
		 	 keeper.DrainInto(&output, keeper.PtrOf(recurse_eval(mem, void)))
				mem = []Token{}
				recursing = false
			} else {
				return output
			}

			default: if recursing {
				keeper.Add(&mem, thing)
			} else {
				keeper.Add(&output, thing)
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
		case "TRUE":   return TokTypeNote(VALUE),    TokType(TRUE)
		case "FALSE":  return TokTypeNote(VALUE),    TokType(FALSE)
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
		case TokType(VOID):         return "[VOID]"
		case TokType(TRUE):         return "[TRUE]"
		case TokType(FALSE):        return "[FALSE]"
		case TokType(QUIT):         return "[QUIT]"
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

func (p *P) seek_first() []byte {
	var mem []byte
	var esc bool
	defer p.back()
	_ = esc
	for p.next() {
		switch p.cur {
			case ' ', '\n', '\t', '\r', ')': return mem
			default: mem = append(mem, p.cur)
		}
	}

	return nil
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
	defer p.back()
	for p.next() {
		if p.cur >= '0' && p.cur <= '9' {
			keeper.Add(&str, p.cur)
		} else {
			return str
		}
	}
	return nil
}

func void_return() []Token {
	return []Token {void}
}

func (t Token) note() string {
	switch t.Note {
		case NONE:     return "[NONE]"
		case FN:       return "[FN]"
		case OPERATOR: return "[OPERATOR]"
		case VALUE:    return "[VALUE]"
		case COMPARE:  return "[COMPARE]"
		case IGNORE:   return "[IGNORE]"
		default: panic("MISSING NOTE TYPE in Token.note()")
	}
}

func (t Token) print() {
	fmt.Printf(
		"token {\n" + 
		"  \x1b[38;2;0;245;240mRaw\x1b[0m:  %#v, \x1b[3;38;2;125;125;125m//%s\x1b[0m\n" + 
		"  \x1b[38;2;0;245;240mType\x1b[0m: %s, \x1b[3;38;2;125;125;125m//%d\x1b[0m\n" +
		"  \x1b[38;2;0;245;240mNote\x1b[0m: %s, \x1b[3;38;2;125;125;125m//%d\x1b[0m\n" +
		"}\n",
		t.Raw, builtin.Un_Escape(t.Raw),
		unmatch_token(t), t.Type,
		t.note(), t.Note,
	)
}
