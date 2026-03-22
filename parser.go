package main

import ("os";"fmt")

func _(){fmt.Print();os.Exit(0)}

func recurse(code []byte) []Token {
	var p = P{
		idx: -1,
		code: code,
	}
	for p.next() {
		switch p.cur {
			case '#': if p.peek() == '|' { p.comment() ; p.toss() }
			case '(': if !p.esc {
				p.Toks = append(p.Toks, mktok(IGNORE, BOX, nil, nil))
				p.next()
				thing, note := get_note(p.seek_first())
				internal_note, t := p.match_name(thing)
				p.Toks = append(p.Toks, mktok(internal_note, t, thing, note))
			}
			case '\\': p.esc = !p.esc
			case '"': if !p.esc {
				p.Toks = append(p.Toks, mktok(VALUE, STRING, p.collapse_str(), nil))
			}
		  case ')':
				p.Toks = append(p.Toks, mktok(IGNORE, EOX, nil, nil))

			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				p.Toks = append(p.Toks, mktok(VALUE, NUMBER, p.seek_num(), nil))

			//skip to newline
			case ';': p.seek_to('\n')

			//might do something with this
			case '\n', ' ', '\r', '\t':

			default: {
				thing, note := get_note(p.seek_first())
				//thing := append([]byte{p.cur}, p.seek_whitespace()...)
				internal_note, t := p.match_name(thing)
				p.Toks = append(p.Toks, mktok(internal_note, t, thing, note))
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
