package ast

import (
	"fmt"
	"toylanguage/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type GenAST struct {
	*parser.BaseGoParserListener
}

func (gen *GenAST) VisitTerminal(node antlr.TerminalNode) {

}

func (gen *GenAST) VisitErrorNode(node antlr.ErrorNode) {

}

func (gen *GenAST) EnterEveryRule(ctx antlr.ParserRuleContext) {
	fmt.Println(ctx.GetText())
}

func (gen *GenAST) ExitEveryRule(ctx antlr.ParserRuleContext) {

}
