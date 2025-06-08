package lexer

import "interpreter/token"

type Lexer struct {
	input           string
	currentPosition int // current position in input (points to current char)
	nextPosition    int // next position in input (after current char)
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.next()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	cur := l.current()
	switch cur {
	case '=': // EQ or ASSIGN
		if l.previewNext() == '=' {
			literal := string(cur) + string(l.next())
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, cur)
		}
	case '+':
		tok = newToken(token.PLUS, cur)
	case '-':
		tok = newToken(token.MINUS, cur)
	case '!': // NOT_EQ or BANG
		if l.previewNext() == '=' {
			literal := string(cur) + string(l.next())
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, cur)
		}
	case '/':
		tok = newToken(token.SLASH, cur)
	case '*':
		tok = newToken(token.ASTERISK, cur)
	case '<':
		tok = newToken(token.LT, cur)
	case '>':
		tok = newToken(token.GT, cur)
	case ';':
		tok = newToken(token.SEMICOLON, cur)
	case ',':
		tok = newToken(token.COMMA, cur)
	case '{':
		tok = newToken(token.LBRACE, cur)
	case '}':
		tok = newToken(token.RBRACE, cur)
	case '(':
		tok = newToken(token.LPAREN, cur)
	case ')':
		tok = newToken(token.RPAREN, cur)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(cur) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(cur) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, cur)
		}
	}

	l.next()
	return tok
}

func (l *Lexer) skipWhitespace() {
	cur := l.current()
	for cur == ' ' || cur == '\t' || cur == '\n' || cur == '\r' {
		cur = l.next()
	}
}

func (l *Lexer) current() byte {
	if l.currentPosition >= len(l.input) {
		return 0
	}
	return l.input[l.currentPosition]
}

func (l *Lexer) next() byte {
	l.currentPosition = l.nextPosition
	l.nextPosition += 1
	return l.current()
}

func (l *Lexer) previewNext() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	}
	return l.input[l.nextPosition]
}

func (l *Lexer) readIdentifier() string {
	currentPosition := l.currentPosition
	for isLetter(l.current()) {
		l.next()
	}
	return l.input[currentPosition:l.currentPosition]
}

func (l *Lexer) readNumber() string {
	currentPosition := l.currentPosition
	for isDigit(l.current()) {
		l.next()
	}
	return l.input[currentPosition:l.currentPosition]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
