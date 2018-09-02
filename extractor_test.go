package extractor

import (
	"regexp"
	"testing"
)

func TestA(t *testing.T) {
	var l *Lexer
	//test scan unknow ident
	l = InitLexer("{")
	if A("{")(l) == false {
		t.Fatalf("A(\"{\") should match \"{\"")
	}
	//test scan int
	l = InitLexer("123")
	if A("123")(l) == false {
		t.Fatalf("A(\"123\") should match \"123\"")
	}

	//test scan int
	l = InitLexer("123.123")
	if A("123.123")(l) == false {
		t.Fatalf("A(\"123\") should match \"123\"")
	}

	//test scan string
	l = InitLexer("\"123 123\"")
	if A("\"123 123\"")(l) == false {
		t.Fatalf("A(\"123 123\") should match \"123 123\"")
	}

	l = InitLexer("abc dd")
	if ASkip("abc")(l) && A("dd")(l) == false {
		t.Fatalf("A('abc') && A('dd') should match \"abc dd\"")
	}

}

func TestEx(t *testing.T) {
	var l *Lexer
	l = InitLexer("abc")
	var result string
	if Ex(&result)(l) == false || result != "abc" {
		t.Fatalf("Ex(result) should get abc, but actual: %v", result)
	}

	l = InitLexer("abc dd")
	if ASkip("abc")(l) && Ex(&result)(l) == false || result != "dd" {
		t.Fatalf("Ex(result) should get dd, but actual: %v", result)
	}
}

func TestExLine(t *testing.T) {
	var l *Lexer
	l = InitLexer("abc abcd\n cced\n")
	var result string
	if ExLine(&result)(l) == false || result != "abc abcd" {
		t.Fatalf("ExLine(result) should get abc abcd, but actual: %v", result)
	}
	if ExLine(&result)(l) == false || result != "cced" {
		t.Fatalf("ExLine(result) should get cced, but actual: %v", result)
	}

	l = InitLexer("abc dd")
	if ASkip("abc")(l) && ExLine(&result)(l) == false || result != "dd" {
		t.Fatalf("Ex(result) should get dd, but actual: %v", result)
	}
}

func TestAnd(t *testing.T) {
	var l *Lexer
	var result string
	l = InitLexer("abc dd")
	if And(ASkip("abc"), Ex(&result))(l) == false || result != "dd" {
		t.Fatalf("Ex(result) should get dd, but actual: %v", result)
	}

	//test seq with no skip before extra
	l = InitLexer("abc dd")
	if And(A("abc"), Ex(&result))(l) == false || result != "abc" {
		t.Fatalf("Ex(result) should get abc, but actual: %v", result)
	}
}

func TestOr(t *testing.T) {
	var l *Lexer
	var result string
	l = InitLexer("abc dd")
	if Or(ASkip("abc"), Ex(&result))(l) == false || result != "" {
		t.Fatalf("Ex(result) should get empty string, but actual: %v", result)
	}

	//test seq with no skip before extra
	l = InitLexer("abc dd")
	if Or(ASkip("abcd"), Ex(&result))(l) == false || result != "abc" {
		t.Fatalf("Ex(result) should get abc, but actual: %v", result)
	}
}

func TestMany(t *testing.T) {
	var l *Lexer
	//multiple matching
	l = InitLexer("abc abc abc dd")
	if Many(ASkip("abc"))(l) == false {
		t.Fatalf("Many(ASkip(abc) should return true")
	}
	if A("dd")(l) == false {
		t.Fatalf("We should reach 'dd' after Many(ASkip(abc))")
	}
	//single matching
	l = InitLexer("abc dd")
	if Many(ASkip("abc"))(l) == false {
		t.Fatalf("Many(ASkip(abc) should return true")
	}
	if A("dd")(l) == false {
		t.Fatalf("We should reach 'dd' after Many(ASkip(abc))")
	}

	//no matching
	l = InitLexer("abc abc abc dd")
	if Many(ASkip("abc"))(l) == false {
		t.Fatalf("Many(ASkip(abc) should return true")
	}
	if A("dd")(l) == false {
		t.Fatalf("We should reach 'dd' after Many(ASkip(abc))")
	}
}

func TestManyWithEx(t *testing.T) {
	var l *Lexer
	l = InitLexer("{abc} {bc} {dc} abc dd")
	var results []string
	var exResults = func(l *Lexer) bool {
		var result string
		success := And(ASkip("{"), Ex(&result), ASkip("}"))(l)
		results = append(results, result)
		return success
	}

	if Many(exResults)(l) == false {
		t.Fatalf("Many(...) should return true")
	}

	if results[0] != "abc" || results[1] != "bc" || results[2] != "dc" {
		t.Fatalf("expected [abc bc dc], actual: %v", results)
	}
}

func TestExIs(t *testing.T) {
	var l *Lexer
	l = InitLexer("!")
	var result bool
	if ExIs(&result, ASkip("!"))(l) == false {
		t.Fatalf("ExIs should return false")
	}

	if result == false {
		t.Fatalf("Result of ExIs should be true")
	}

	l = InitLexer("abb! ddd")

	if And(ASkip("abb"), ExIs(&result, ASkip("!")), ASkip("ddd"))(l) == false {
		t.Fatalf("Seq(ASkip(abb), ExIs(&result,!), ASkip(ddd) should return true for 'abb! ddd'")
	}

	if result == false {
		t.Fatalf("Result of ExIs should be true")
	}

	l = InitLexer("abb ddd")

	if And(ASkip("abb"), ExIs(&result, ASkip("!")), ASkip("ddd"))(l) == false {
		t.Fatalf("Seq(ASkip(abb), ExIs(&result,!), ASkip(ddd) should return true for 'abb! ddd'")
	}

	if result == true {
		t.Fatalf("Result of ExIs should be fase")
	}
}

func BenchmarkExtract(b *testing.B) {
	var l *Lexer
	for i := 0; i < b.N; i++ {
		l = InitLexer("abc(dd)")
		l.GetText()
		//var result string
		//_ = ASkip("abc")(l) //&& ASkip("(")(l) && Ex(&result)(l)
	}
}

func BenchmarkRegexp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rgx := regexp.MustCompile(`abc\((.*?)\)`)
		src := `abc(tag)`
		_ = rgx.FindStringSubmatch(src)
	}
}
