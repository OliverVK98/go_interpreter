package parser

import (
	"go-interpreter/ast"
	"go-interpreter/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `let x = 5;
let y = 10;
let foobar = 8388383;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParsesErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program doesn't contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not `let`. Got %q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not `%s`. got %s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral not `%s`. got %s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `return 5;
return 10;
return 8388383;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParsesErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program doesn't contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnstmt, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%v", stmt)
			continue
		}

		if returnstmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not `return`. got %v", returnstmt.TokenLiteral())
		}
	}
}

func checkParsesErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parses had %v errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %v", msg)
	}

	t.FailNow()

}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParsesErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has incorrect number of statements. got = %v", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got = %v", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got = %v", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %v. got %v", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %v. got %v", "foobar", ident.TokenLiteral())
	}
}
