package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENTIFIER" // add, foobar, x, y, ...
	INT    = "INT"        // 1343456
	FLOAT  = "FLOAT"      // 3.142
	STRING = "STRING"     // "foobar"

	// Operators
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	ASSIGN          = "="
	PLUS_ASSIGN     = "+="
	MINUS_ASSIGN    = "-="
	ASTERISK_ASSIGN = "*="
	SLASH_ASSIGN    = "/="

	LT  = "<"
	GT  = ">"
	LTE = "<="
	GTE = ">="

	EQ        = "=="
	NOT_EQ    = "!="
	POST_INCR = "++"
	POST_DECR = "--"

	AND = "and"
	OR  = "or"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	FULLSTOP  = "."

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	WHILE    = "WHILE"
	FOR      = "FOR"
	IN       = "IN"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"func":     FUNCTION,
	"let":      LET,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"and":      AND,
	"or":       OR,
	"while":    WHILE,
	"for":      FOR,
	"in":       IN,
	"break":    BREAK,
	"continue": CONTINUE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
