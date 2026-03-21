package main

import (_"os";"fmt")

func _(){fmt.Print()}

func recurse(code []byte) []Token {
	var p = P{
		idx: -1,
		code: code,
	}
	for p.next() {
		switch p.cur {
			case '#': if p.peek() == '|' { p.comment() ; p.toss() }
			case '(': if !p.esc {
				p.Toks = append(p.Toks, mktok(IGNORE, BOX, nil))
				thing := p.seek_whitespace()
				note, t := p.match_name(thing)
				p.Toks = append(p.Toks, mktok(note, t, thing))
			}
			case '\\': p.esc = !p.esc
			case '"': if !p.esc {
				p.Toks = append(p.Toks, mktok(VALUE, STRING, p.collapse_str()))
			}
		  case ')':
				p.Toks = append(p.Toks, mktok(IGNORE, EOX, nil))
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.Toks = append(p.Toks, mktok(VALUE, NUMBER, p.seek_num()))

			//skip to newline
			case ';': p.seek_to('\n')

			//might do something with this
			case '\n', ' ', '\r', '\t':

			default: {
				thing := append([]byte{p.cur}, p.seek_whitespace()...)
				note, t := p.match_name(thing)
				p.Toks = append(p.Toks, mktok(note, t, thing))
			}
		}
	}

	return p.Toks
}

func (p *P) comment() {
	p.toss()
	var depth int ; loop: for p.next() {
		if p.cur == '#' && p.peek() == '|' { depth++ ; p.idx++ ; continue loop }
		if p.cur == '|' && p.peek() == '#' {
			if depth == 0 { p.toss() ; return }
			depth-- ; continue loop
		}
	}
}
