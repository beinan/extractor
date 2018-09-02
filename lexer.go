package extractor

import (
	"strings"
	"text/scanner"
)

type Lexer struct {
	src  string
	s    *scanner.Scanner
	tok  rune   //last rune
	text string //the token text scanner just scanned
}

type Position struct {
	Offset int
	Line   int
	Column int
}

func InitLexer(text string) *Lexer {
	s := &scanner.Scanner{
		Mode: scanner.ScanIdents | scanner.ScanFloats | scanner.ScanInts | scanner.ScanStrings,
	}
	s.Init(strings.NewReader(text))
	l := &Lexer{
		src: text,
		s:   s,
	}
	l.Next()
	return l
}

func (l *Lexer) Next() bool {
	l.tok = l.s.Scan()
	if l.tok == scanner.EOF {
		return false
	}
	l.text = l.s.TokenText()
	return true
}

func (l *Lexer) GetText() string {
	return l.text
}

func (l *Lexer) GetTok() rune {
	return l.tok
}

func (l *Lexer) IsIdent() bool {
	return l.tok == scanner.Ident
}

func (l *Lexer) LineStr() string {
	var line strings.Builder
	line.WriteString(l.GetText())
	for {
		next := l.s.Next()
		if next == '\n' || next == scanner.EOF {
			//l.Next() //go to next line
			break
		}
		line.WriteRune(next)
	}
	return line.String()
}

func (l *Lexer) Pos() Position {
	return Position{
		Offset: l.s.Offset,
		Line:   l.s.Line,
		Column: l.s.Column,
	}
}
