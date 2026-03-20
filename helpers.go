package main

import (
	keeper "github.com/Supraboy981322/keeper/golang"
)

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
	for (*input)[0].Type != EOX {
		keeper.Add(&output, (*input)[0])
		keeper.Shift(input)
		if len(*input) < 1 { return output }
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

// TODO: more universal seek (ie: not having separate logic for string)
func (p *P) seek_to(c byte) []byte {
	var mem []byte
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
