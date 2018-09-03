package extractor

import (
	"fmt"
	"text/scanner"
)

type Op func(*Lexer) bool

type ExtractError struct {
	Msg string
	Pos Position
}

func (e *ExtractError) Error() string {
	return fmt.Sprintf("%v  (line:%v, column:%v, offset:%v",
		e.Msg, e.Pos.Line, e.Pos.Column, e.Pos.Offset)
}

func A(word string) func(*Lexer) bool {
	return func(l *Lexer) bool {
		fmt.Println("a word expected", word, " actual:", l.GetText())
		if l.GetText() != word {
			return false
		}
		return true
	}
}

func Must(word string) Op {
	return func(l *Lexer) bool {
		if A(word)(l) {
			Skip(l)
			return true
		}
		return ThrowError(fmt.Sprintf("Expected: %v, Actual: %v", word, l.GetText()))(l)
	}
}

func ThrowError(msg string) Op {
	return func(l *Lexer) bool {
		panic(&ExtractError{
			Msg: msg,
			Pos: l.Pos(),
		})
	}
}

func Skip(l *Lexer) bool {
	l.Next()
	return true
}

func ASkip(word string) func(*Lexer) bool {
	return func(l *Lexer) bool {
		return A(word)(l) && Skip(l)
	}
}

func Ex(result *string) func(*Lexer) bool {
	return func(l *Lexer) bool {
		//fmt.Println("extract:", l.GetText())
		*result = l.GetText()
		l.Next()
		return true
	}
}

func ExId(result *string) func(*Lexer) bool {
	return func(l *Lexer) bool {
		//fmt.Println("extract:", l.GetText())
		if l.GetTok() == scanner.Ident {
			*result = l.GetText()
			l.Next()
			return true
		}
		return false
	}
}

func ExLine(result *string) func(*Lexer) bool {
	return func(l *Lexer) bool {
		*result = l.LineStr()
		l.Next()
		return true
	}
}

func And(ops ...Op) Op {
	return func(l *Lexer) bool {
		for _, op := range ops {
			if !op(l) {
				return false
			}
		}
		return true
	}
}

func Or(ops ...Op) Op {
	return func(l *Lexer) bool {
		for _, op := range ops {
			if op(l) {
				return true
			}
		}
		return false
	}
}

func Many(op Op) Op {
	return func(l *Lexer) bool {
		for op(l) {
		}
		return true
	}
}

func Option(op Op) Op {
	return func(l *Lexer) bool {
		op(l)
		return true
	}
}

func ExIs(result *bool, op Op) Op {
	return func(l *Lexer) bool {
		*result = op(l)
		return true
	}
}
