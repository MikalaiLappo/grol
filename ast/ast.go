// Abstract Syntax Tree for the GROL language.
// Everything is Node. Has a Token() and can be PrettyPrint'ed back to source
// that would parse to the same AST.
package ast

import (
	"io"
	"strconv"
	"strings"

	"fortio.org/log"
	"grol.io/grol/token"
)

type PrintState struct {
	Out             io.Writer
	IndentLevel     int
	ExpressionLevel int
	IndentationDone bool // already put N number of tabs, reset on each new line
	Compact         bool // don't indent at all (compact mode), no newlines, fewer spaces, no comments
	prev            Node
	last            string
}

func DebugString(n Node) string {
	ps := NewPrintState()
	ps.Compact = true
	n.PrettyPrint(ps)
	return ps.String()
}

func NewPrintState() *PrintState {
	return &PrintState{Out: &strings.Builder{}}
}

func (ps *PrintState) String() string {
	return ps.Out.(*strings.Builder).String()
}

// Will print indented to current level. with a newline at the end.
// Only a single indentation per line.
func (ps *PrintState) Println(str ...string) *PrintState {
	ps.Print(str...)
	if !ps.Compact {
		_, _ = ps.Out.Write([]byte{'\n'})
	}
	ps.IndentationDone = false
	return ps
}

func (ps *PrintState) Print(str ...string) *PrintState {
	if len(str) == 0 {
		return ps // So for instance Println() doesn't print \t\n.
	}
	if !ps.Compact && !ps.IndentationDone && ps.IndentLevel > 1 {
		_, _ = ps.Out.Write([]byte(strings.Repeat("\t", ps.IndentLevel-1)))
		ps.IndentationDone = true
	}
	for _, s := range str {
		_, _ = ps.Out.Write([]byte(s))
		ps.last = s
	}
	return ps
}

// --- AST nodes

// Everything in the tree is a Node.
type Node interface {
	Value() *token.Token
	PrettyPrint(ps *PrintState) *PrintState
}

// Common to all nodes that have a token and avoids repeating the same TokenLiteral() methods.
type Base struct {
	*token.Token
}

func (b Base) Value() *token.Token {
	return b.Token
}

func (b Base) PrettyPrint(ps *PrintState) *PrintState {
	// In theory should only be called for literals.
	// log.Debugf("PrettyPrint on base called for %T", b.Value())
	return ps.Print(b.Literal())
}

type ReturnStatement struct {
	Base
	ReturnValue Node
}

func (rs ReturnStatement) PrettyPrint(ps *PrintState) *PrintState {
	ps.Print(rs.Literal())
	if rs.ReturnValue != nil {
		ps.Print(" ")
		rs.ReturnValue.PrettyPrint(ps)
	}
	return ps
}

type Statements struct {
	Base
	Statements []Node
}

func keepSameLineAsPrevious(node Node) bool {
	switch n := node.(type) { //nolint:exahustive // we may add more later
	case *Comment:
		return n.SameLineAsPrevious
	default:
		return false
	}
}

func needNewLineAfter(node Node) bool {
	switch n := node.(type) { //nolint:exahustive // we may add more later
	case *Comment:
		return !n.SameLineAsNext
	default:
		return true
	}
}

func isComment(node Node) bool {
	_, ok := node.(*Comment)
	return ok
}

// Compact mode: Skip comments and decide if we need a space separator or not.
func prettyPrintCompact(ps *PrintState, s Node, i int) bool {
	if isComment(s) {
		return true
	}
	_, prevIsExpr := ps.prev.(*InfixExpression)
	// _, curIsID := s.(*Identifier)
	if prevIsExpr && ps.last != "}" && ps.last != "]" {
		if i > 0 {
			_, _ = ps.Out.Write([]byte{' '})
		}
	}
	return false
}

// Normal/long form print: Decide if using new line or space as separator.
func prettyPrintLongForm(ps *PrintState, s Node, i int) {
	if i > 0 || ps.IndentLevel > 1 {
		if keepSameLineAsPrevious(s) || !needNewLineAfter(ps.prev) {
			log.Debugf("=> PrettyPrint adding just a space")
			_, _ = ps.Out.Write([]byte{' '})
			ps.IndentationDone = true
		} else {
			log.Debugf("=> PrettyPrint adding newline")
			ps.Println()
		}
	}
}

func (p Statements) PrettyPrint(ps *PrintState) *PrintState {
	oldExpressionLevel := ps.ExpressionLevel
	if ps.IndentLevel > 0 {
		ps.Print("{") // first statement might be a comment on same line.
	}
	ps.IndentLevel++
	ps.ExpressionLevel = 0
	var i int
	for _, s := range p.Statements {
		if ps.Compact {
			if prettyPrintCompact(ps, s, i) {
				continue // skip comments entirely.
			}
		} else {
			prettyPrintLongForm(ps, s, i)
		}
		s.PrettyPrint(ps)
		ps.prev = s
		i++
	}
	ps.Println()
	ps.IndentLevel--
	ps.ExpressionLevel = oldExpressionLevel
	if ps.IndentLevel > 0 {
		ps.Print("}")
	}
	return ps
}

type Identifier struct {
	Base
}

func (i Identifier) PrettyPrint(out *PrintState) *PrintState {
	out.Print(i.Literal())
	return out
}

type Comment struct {
	Base
	SameLineAsPrevious bool
	SameLineAsNext     bool
}

func (c Comment) PrettyPrint(out *PrintState) *PrintState {
	out.Print(c.Literal())
	return out
}

type IntegerLiteral struct {
	Base
	Val int64
}

type FloatLiteral struct {
	Base
	Val float64
}

type StringLiteral struct {
	Base
	// Val string // Base.Token.Literal is enough to store the string value.
}

func (s StringLiteral) PrettyPrint(ps *PrintState) *PrintState {
	ps.Print(strconv.Quote(s.Literal()))
	return ps
}

type PrefixExpression struct {
	Base
	Right Node
}

func (p PrefixExpression) PrettyPrint(out *PrintState) *PrintState {
	if out.ExpressionLevel > 0 {
		out.Print("(")
	}
	out.Print(p.Literal())
	out.ExpressionLevel++ // comment out for !(-a) to normalize to !-a
	p.Right.PrettyPrint(out)
	out.ExpressionLevel--
	if out.ExpressionLevel > 0 {
		out.Print(")")
	}
	return out
}

type PostfixExpression struct {
	Base
	Prev *token.Token
}

func (p PostfixExpression) PrettyPrint(out *PrintState) *PrintState {
	if out.ExpressionLevel > 0 {
		out.Print("(")
	}
	out.Print(p.Prev.Literal())
	out.Print(p.Literal())
	if out.ExpressionLevel > 0 {
		out.Print(")")
	}
	return out
}

type InfixExpression struct {
	Base
	Left  Node
	Right Node
}

func (i InfixExpression) PrettyPrint(out *PrintState) *PrintState {
	if out.ExpressionLevel > 0 { // TODO only add parens if precedence requires it.
		out.Print("(")
	}
	isAssign := (i.Token.Type() == token.ASSIGN)
	if !isAssign {
		out.ExpressionLevel++
	}
	i.Left.PrettyPrint(out)
	if out.Compact {
		out.Print(i.Literal())
	} else {
		out.Print(" ", i.Literal(), " ")
	}
	if i.Right == nil {
		out.Print("nil")
	} else {
		i.Right.PrettyPrint(out)
	}
	if !isAssign {
		out.ExpressionLevel--
	}
	if out.ExpressionLevel > 0 {
		out.Print(")")
	}
	return out
}

type Boolean struct {
	Base
	Val bool
}

type IfExpression struct {
	Base
	Condition   Node
	Consequence *Statements
	Alternative *Statements
}

func (ie IfExpression) printElse(out *PrintState) {
	if out.Compact {
		out.Print("else")
	} else {
		out.Print(" else ")
	}
	if len(ie.Alternative.Statements) == 1 && ie.Alternative.Statements[0].Value().Type() == token.IF {
		// else if
		if out.Compact {
			out.Print(" ")
		}
		ie.Alternative.Statements[0].PrettyPrint(out)
		return
	}
	ie.Alternative.PrettyPrint(out)
}

func (ie IfExpression) PrettyPrint(out *PrintState) *PrintState {
	out.Print("if ")
	ie.Condition.PrettyPrint(out)
	if !out.Compact {
		out.Print(" ")
	}
	ie.Consequence.PrettyPrint(out)
	if ie.Alternative != nil {
		ie.printElse(out)
	}
	return out
}

func PrintList(out *PrintState, list []Node, sep string) {
	for i, p := range list {
		if i > 0 {
			out.Print(sep)
		}
		p.PrettyPrint(out)
	}
}

// Similar to CallExpression.
type Builtin struct {
	Base       // The 'len' or 'first' or... core builtin token
	Parameters []Node
}

func (b Builtin) PrettyPrint(out *PrintState) *PrintState {
	out.Print(b.Literal())
	out.Print("(")
	out.ComaList(b.Parameters)
	out.Print(")")
	return out
}

type FunctionLiteral struct {
	Base       // The 'func' token
	Name       *Identifier
	Parameters []Node // last one might be `..` for variadic.
	Body       *Statements
	Variadic   bool
}

func (fl FunctionLiteral) PrettyPrint(out *PrintState) *PrintState {
	out.Print(fl.Literal())
	if fl.Name != nil {
		out.Print(" ")
		out.Print(fl.Name.Literal())
	}
	out.Print("(")
	out.ComaList(fl.Parameters)
	if out.Compact {
		out.Print(")")
	} else {
		out.Print(") ")
	}
	fl.Body.PrettyPrint(out)
	return out
}

type CallExpression struct {
	Base           // The '(' token
	Function  Node // Identifier or FunctionLiteral
	Arguments []Node
}

func (ce CallExpression) PrettyPrint(out *PrintState) *PrintState {
	ce.Function.PrettyPrint(out)
	out.Print("(")
	oldExpressionLevel := out.ExpressionLevel
	out.ExpressionLevel = 0
	out.ComaList(ce.Arguments)
	out.ExpressionLevel = oldExpressionLevel
	out.Print(")")
	return out
}

type ArrayLiteral struct {
	Base     // The [ token
	Elements []Node
}

func (al ArrayLiteral) PrettyPrint(out *PrintState) *PrintState {
	out.Print("[")
	out.ComaList(al.Elements)
	out.Print("]")
	return out
}

type IndexExpression struct {
	Base
	Left  Node
	Index Node
}

func (ie IndexExpression) PrettyPrint(out *PrintState) *PrintState {
	if out.ExpressionLevel > 0 { // TODO only add parens if precedence requires it.
		out.Print("(")
	}
	ie.Left.PrettyPrint(out)
	out.Print(ie.Literal())
	ie.Index.PrettyPrint(out)
	if ie.Token.Type() == token.LBRACKET {
		out.Print("]")
	}
	if out.ExpressionLevel > 0 {
		out.Print(")")
	}
	return out
}

type MapLiteral struct {
	Base  // the '{' token
	Pairs map[Node]Node
	Order []Node // for pretty printing in same order as input
}

func (hl MapLiteral) PrettyPrint(out *PrintState) *PrintState {
	out.Print("{")
	sep := ", "
	if out.Compact {
		sep = ","
	}
	for i, key := range hl.Order {
		if i > 0 {
			out.Print(sep)
		}
		key.PrettyPrint(out)
		out.Print(":")
		hl.Pairs[key].PrettyPrint(out)
	}
	out.Print("}")
	return out
}

type MacroLiteral struct {
	Base
	Parameters []Node
	Body       *Statements
}

func (ml MacroLiteral) PrettyPrint(out *PrintState) *PrintState {
	out.Print(ml.Literal())
	out.Print("(")
	out.ComaList(ml.Parameters)
	if out.Compact {
		out.Print(")")
	} else {
		out.Print(") ")
	}
	ml.Body.PrettyPrint(out)
	return out
}

func (ps *PrintState) ComaList(list []Node) {
	sep := ", "
	if ps.Compact {
		sep = ","
	}
	PrintList(ps, list, sep)
}
