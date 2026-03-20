package main

import ("os";"fmt")

type TokType int
const (
	INVALID TokType = iota
	PRINT
	IF
	ELSE
	OR
	AND
	STRING
)

type Builtin struct{}
var builtin Builtin

type P struct{
	idx int
	cur byte
	code []byte
	Toks []Token
}
type Token struct {
	Raw []byte
	Type TokType
}
func main() {
	eval(recurse([]byte(os.Args[1])))
}

func recurse(code []byte) []Token {
	var p = P{
		code: code,
	}
	for p.next() {
		switch p.cur {
			case '(': p.Toks = append(p.Toks, recurse(p.seek_to(')'))...)
			case '"': {
				p.Toks = append(p.Toks, mktok(STRING, p.seek_to('"')))
			}
		  case ' ': panic("foo")
			case ')':
			default: {
				c := p.cur
				thing := append([]byte{c}, p.seek_to(' ')...)
				new_tok := Token {
					Raw: thing,
					Type: p.match_name(thing),
				}
				p.Toks = append(p.Toks, new_tok)
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

func mktok(t TokType, raw []byte) Token {
	return Token {
		Type: t,
		Raw: raw,
	}
}

func (P) match_name(name []byte) TokType {
	switch string(name) {
		case "print": return TokType(PRINT)
		case "?": return TokType(IF)
		case "?!": return TokType(ELSE)
		case "&": return TokType(AND)
		case "|": return TokType(OR)
		default: return TokType(INVALID)
	}
}

func eval(input []Token) {
	if len(input) < 1 { return }
	for _, t := range input {
		fmt.Println("|" + string(t.Raw) + "|")
	}
	switch input[0].Type {
		case TokType(PRINT): { builtin.Print(input[1].Raw) }
		default: {
			fmt.Printf("undeclared fn name: %s", input[0].Raw)
			os.Exit(1)
		}
	}
}

func (p *P) next() bool {
	p.cur = p.seek()
	if p.cur == 0 { return false }
	return true
}

func (p *P) peak() byte {
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
	b := p.peak()  ;  p.toss()
	return b
}

func (p *P) seekN(n int) []byte {
	s := p.peakN(n)  ;  p.tossN(n)
	return s
}

func (p *P) seek_to(c byte) []byte {
	var mem []byte
	for p.next() {
		if p.cur == c { return mem }
		mem = append(mem, p.cur)
	}
	return nil
}

func (p *P) peakN(n int) []byte {
	if len(p.code) <= p.idx+n { return nil }
	return p.code[p.idx:][:n]
}

func (Builtin) Print(str []byte) {
	// TODO: parse string for escapes
	fmt.Println(string(str))
}
