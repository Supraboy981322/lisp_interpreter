package main

import ("os";"fmt")


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
