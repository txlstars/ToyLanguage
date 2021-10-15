package ast

import (
	"bytes"
	"fmt"
	"strings"
	"toylanguage/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type DumpASTTool struct {
	parser.BaseGoParserVisitor
	buf    *strings.Builder
	indent []byte
}

func NewDumpASTTool() *DumpASTTool {
	return &DumpASTTool{
		buf: new(strings.Builder),
	}
}

func (d *DumpASTTool) Dump() {
	fmt.Println(d.buf.String())
}

func (d *DumpASTTool) writeString(texts ...string) {
	var text bytes.Buffer
	text.Write(d.indent)
	for i, v := range texts {
		text.WriteString(v)
		if i != len(texts)-1 {
			text.WriteByte(' ')
		}
	}
	text.WriteByte('\n')
	d.buf.WriteString(text.String())
}

func (d *DumpASTTool) pushIndent() {
	d.indent = append(d.indent, ' ', ' ')
}

func (d *DumpASTTool) popIndent() {
	d.indent = d.indent[:len(d.indent)-2]
}

func (v *DumpASTTool) Visit(tree antlr.ParseTree) interface{} {
	if tree != nil {
		return tree.Accept(v)
	}
	return nil
}

func (d *DumpASTTool) VisitChildren(node antlr.RuleNode) interface{} {
	for _, v := range node.GetChildren() {
		v.(antlr.ParseTree).Accept(d)
	}
	return nil
}

func (d *DumpASTTool) VisitSourceFile(ctx *parser.SourceFileContext) interface{} {
	d.writeString("SourceFile")
	d.pushIndent()
	defer d.popIndent()
	d.VisitChildren(ctx)
	return nil
}

func (d *DumpASTTool) VisitPackageClause(ctx *parser.PackageClauseContext) interface{} {
	d.writeString("PackageClause")
	d.pushIndent()
	defer d.popIndent()
	d.writeString(ctx.GetPackageName().GetText())
	return nil
}

func (d *DumpASTTool) VisitImportDecl(ctx *parser.ImportDeclContext) interface{} {
	d.writeString("ImportDecl")
	d.pushIndent()
	defer d.popIndent()
	d.VisitChildren(ctx)
	return nil
}

func (d *DumpASTTool) VisitImportSpec(ctx *parser.ImportSpecContext) interface{} {
	t := []string{}
	if alias := ctx.GetAlias(); alias != nil {
		t = append(t, ctx.GetAlias().GetText())
	}
	t = append(t, d.VisitImportPath(ctx.ImportPath().(*parser.ImportPathContext)).(string))
	d.writeString(t...)
	return nil
}

func (d *DumpASTTool) VisitImportPath(ctx *parser.ImportPathContext) interface{} {
	return ctx.String_().GetText()
}

func (d *DumpASTTool) VisitFunctionDecl(ctx *parser.FunctionDeclContext) interface{} {
	d.writeString("FunctionDecl")
	return nil
}

func (d *DumpASTTool) VisitMethodDecl(ctx *parser.MethodDeclContext) interface{} {
	d.writeString("MethodDecl")
	return nil
}

func (d *DumpASTTool) VisitDeclaration(ctx *parser.DeclarationContext) interface{} {
	d.VisitChildren(ctx)
	return nil
}

func (d *DumpASTTool) VisitConstDecl(ctx *parser.ConstDeclContext) interface{} {
	d.writeString("ConstDecl")
	return nil
}

func (d *DumpASTTool) VisitTypeDecl(ctx *parser.TypeDeclContext) interface{} {
	d.writeString("TypeDecl")
	return nil
}

func (d *DumpASTTool) VisitVarDecl(ctx *parser.VarDeclContext) interface{} {
	d.writeString("VarDecl")
	return nil
}
