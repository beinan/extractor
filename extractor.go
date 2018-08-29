package extractor

type Op func(*Lexer) bool

func A(word string) func(*Lexer) bool {
	return func(l *Lexer) bool {
		//fmt.Println("a word expected", word, " actual:", l.GetText())
		if l.GetText() != word {
			return false
		}
		return true
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

func Seq(ops ...Op) Op {
	return func(l *Lexer) bool {
		for _, op := range ops {
			if !op(l) {
				return false
			}
		}
		return true
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
