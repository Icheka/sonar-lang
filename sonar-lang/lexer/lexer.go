package lexer

import (
	"strings"

	"github.com/icheka/sonar-lang/sonar-lang/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	if string(l.ch) == token.SLASH && string(l.peekChar()) == token.SLASH {
		l.skipSingleLineComment()
		return l.NextToken()
	}

	if l.ch == '/' && l.peekChar() == '*' {
		l.skipMultiLineComment()
	}

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		switch l.peekChar() {
		case '+':
			l.readChar()
			tok = newToken(token.POST_INCR, token.POST_INCR)
		case '=':
			l.readChar()
			tok = newToken(token.PLUS_ASSIGN, token.PLUS_ASSIGN)
		default:
			tok = newToken(token.PLUS, l.ch)
		}
	case '-':
		switch l.peekChar() {
		case '=':
			l.readChar()
			tok = newToken(token.MINUS_ASSIGN, token.MINUS_ASSIGN)
		case '-':
			l.readChar()
			tok = newToken(token.POST_DECR, token.POST_DECR)
		default:
			tok = newToken(token.MINUS, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		switch l.peekChar() {
		case '=':
			l.readChar()
			tok = newToken(token.SLASH_ASSIGN, token.SLASH_ASSIGN)
		default:
			tok = newToken(token.SLASH, l.ch)
		}
	case '*':
		switch l.peekChar() {
		case '=':
			l.readChar()
			tok = newToken(token.ASTERISK_ASSIGN, token.ASTERISK_ASSIGN)
		default:
			tok = newToken(token.ASTERISK, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			return l.readNumber()
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipSingleLineComment() {
	for l.ch != '\n' && l.ch != byte(0) {
		l.readChar()
	}
	l.skipWhitespace()
}

func (l *Lexer) skipMultiLineComment() {
	scanning := true
	for scanning {
		switch l.ch {
		case 0:
			scanning = false
		case '*':
			if l.peekChar() == '/' {
				scanning = false
				l.readChar() // advance to '/'
			}
		}
		l.readChar()
	}
	l.skipWhitespace()
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || (isDigit(l.ch) && isLetter(l.input[l.position-1])) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() token.Token {
	position := l.position

	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}
	text := l.input[position:l.position]
	dotCount := strings.Count(text, ".")
	tok := &token.Token{}

	if dotCount > 1 {
		return newToken(token.ILLEGAL, text)
	}

	switch dotCount {
	case 0:
		tok.Type = token.INT
	case 1:
		tok.Type = token.FLOAT
	}
	tok.Literal = text
	return *tok
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken[T byte | string](tokenType token.TokenType, ch T) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
