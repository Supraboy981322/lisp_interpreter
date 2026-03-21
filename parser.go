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

			//skip to newline  TODO: multi-line comments
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
