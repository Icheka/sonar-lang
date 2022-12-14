package ast

import (
	"bytes"
	"strings"

	"github.com/icheka/sonar-lang/sonar-lang/token"
)

// The base Node interface
type Node interface {
	TokenLiteral() string
	String() string
}

// All statement nodes implement this
type Statement interface {
	Node
	statementNode()
}

// All expression nodes implement this
type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Statements
type LetStatement struct {
	Token     token.Token // the token.LET token
	Name      *Identifier
	Value     Expression
	TokenInfo interface{}
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
	TokenInfo   interface{}
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
	TokenInfo  interface{}
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
	TokenInfo  interface{}
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type WhileStatement struct {
	Token       token.Token // the 'while' token
	Condition   Expression
	Consequence *BlockStatement
	TokenInfo   interface{}
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ws.TokenLiteral())
	out.WriteString(" (")
	out.WriteString(ws.Condition.String())
	out.WriteString(" ) {")
	out.WriteString(ws.Consequence.String())
	out.WriteString("}")

	return out.String()
}

type ForStatement struct {
	Token       token.Token // the 'for' token
	Counter     Node        // the 'i' in 'for (i, v in [0, 1])'
	Value       Node        // the 'v' part in 'for (i, v in [0, 1])'
	Operator    token.Token // the infix operator used. For now, and maybe forever, it will always be 'in'
	Iterable    Expression
	Consequence *BlockStatement
	TokenInfo   interface{}
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	var out bytes.Buffer

	out.WriteString(fs.TokenLiteral())
	out.WriteString(" (")
	out.WriteString(fs.Counter.String())
	out.WriteString(", ")
	out.WriteString(fs.Value.String())
	out.WriteString(" in ")
	out.WriteString(fs.Iterable.String())
	out.WriteString(" ) {")
	out.WriteString(fs.Consequence.String())
	out.WriteString("}")

	return out.String()
}

type BreakStatement struct {
	Token     token.Token // the 'break' token
	TokenInfo interface{}
}

func (bs *BreakStatement) statementNode()       {}
func (bs *BreakStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BreakStatement) String() string       { return bs.TokenLiteral() }

type ContinueStatement struct {
	Token     token.Token // the 'continue' token
	TokenInfo interface{}
}

func (cs *ContinueStatement) statementNode()       {}
func (cs *ContinueStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ContinueStatement) String() string       { return cs.TokenLiteral() }

// *****************
//  * Expressions *
// *****************
type Identifier struct {
	Token     token.Token // the token.IDENT token
	Value     string
	TokenInfo interface{}
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type Boolean struct {
	Token     token.Token
	Value     bool
	TokenInfo interface{}
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type IntegerLiteral struct {
	Token     token.Token
	Value     int64
	TokenInfo interface{}
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type FloatLiteral struct {
	Token     token.Token
	Value     float64
	TokenInfo interface{}
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

type PrefixExpression struct {
	Token     token.Token // The prefix token, e.g. !
	Operator  string
	Right     Expression
	TokenInfo interface{}
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token     token.Token // The operator token, e.g. +
	Left      Expression
	Operator  string
	Right     Expression
	TokenInfo interface{}
}

func (oe *InfixExpression) expressionNode()      {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")

	return out.String()
}

type PostfixExpression struct {
	Token     token.Token
	Operator  string
	TokenInfo interface{}
}

func (oe *PostfixExpression) expressionNode()      {}
func (oe *PostfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *PostfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(oe.Token.Literal)
	out.WriteString(oe.Operator)
	out.WriteString(")")

	return out.String()
}

type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
	TokenInfo   interface{}
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
	TokenInfo  interface{}
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
	TokenInfo interface{}
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type StringLiteral struct {
	Token     token.Token
	Value     string
	TokenInfo interface{}
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type ArrayLiteral struct {
	Token     token.Token // the '[' token
	Elements  []Expression
	TokenInfo interface{}
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type IndexExpression struct {
	Token     token.Token // The [ token
	Left      Expression
	Index     Expression
	TokenInfo interface{}
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

type HashLiteral struct {
	Token     token.Token // the '{' token
	Pairs     map[Expression]Expression
	TokenInfo interface{}
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type AssignmentExpression struct {
	Token      token.Token // the identifier token
	Identifier *Identifier
	Value      Expression
	Operator   string
	TokenInfo  interface{}
}

func (as *AssignmentExpression) expressionNode()      {}
func (as *AssignmentExpression) TokenLiteral() string { return as.Token.Literal }
func (as *AssignmentExpression) String() string {
	var out *bytes.Buffer

	out.WriteString(as.Identifier.String())
	out.WriteString(" " + as.Operator + " ")
	out.WriteString(as.Value.String())

	return out.String()
}

type SquareBracketAssignment struct {
	Token     token.Token // the [ token
	Value     Expression
	Key       Expression
	Left      Expression
	TokenInfo interface{}
}

func (as *SquareBracketAssignment) expressionNode()      {}
func (as *SquareBracketAssignment) TokenLiteral() string { return as.Token.Literal }
func (as *SquareBracketAssignment) String() string {
	var out *bytes.Buffer

	out.WriteString(as.Left.String())
	out.WriteString("[")
	out.WriteString(as.Key.String())
	out.WriteString("] = ")
	out.WriteString(as.Value.String())

	return out.String()
}

type NullValueExpression struct {
	TokenInfo interface{}
}

func (nv *NullValueExpression) expressionNode()      {}
func (nv *NullValueExpression) TokenLiteral() string { return "null" }
func (nv *NullValueExpression) String() string {
	return nv.TokenLiteral()
}
