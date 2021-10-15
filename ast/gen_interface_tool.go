package ast

import (
	"fmt"
	"go/format"
	"strings"
	"toylanguage/parser"
)

type GenInterfaceTool struct {
	parser.BaseGoParserListener
	structTypeInfo map[string][]*parser.MethodDeclContext
}

func NewGenInterfaceTool() *GenInterfaceTool {
	return &GenInterfaceTool{
		structTypeInfo: make(map[string][]*parser.MethodDeclContext),
	}
}

func (g *GenInterfaceTool) EnterMethodDecl(ctx *parser.MethodDeclContext) {
	receiver := ctx.Receiver().(*parser.ReceiverContext)
	parameters := receiver.Parameters().(*parser.ParametersContext)
	parameterDecl := parameters.ParameterDecl(0).(*parser.ParameterDeclContext)
	typ := parameterDecl.Type_().(*parser.Type_Context)
	typName := genReceiverTypeName(typ)
	g.structTypeInfo[typName] = append(g.structTypeInfo[typName], ctx)
}

func (s *GenInterfaceTool) CodeGen() {
	for k, v := range s.structTypeInfo {
		genInterface(k, v)
	}
}

func genReceiverTypeName(ctx *parser.Type_Context) string {
	if typName := ctx.TypeName(); typName != nil {
		return typName.(*parser.TypeNameContext).IDENTIFIER().GetText()
	}

	if typ := ctx.Type_(); typ != nil {
		return genTypeName(typ.(*parser.Type_Context))
	}
	typLit := ctx.TypeLit().(*parser.TypeLitContext)
	pointerType := typLit.PointerType().(*parser.PointerTypeContext)
	return genTypeName(pointerType.Type_().(*parser.Type_Context))
}

func genTypeName(ctx *parser.Type_Context) string {
	if typName := ctx.TypeName(); typName != nil {
		return typName.(*parser.TypeNameContext).IDENTIFIER().GetText()
	}

	if typ := ctx.Type_(); typ != nil {
		return genTypeName(typ.(*parser.Type_Context))
	}

	typLit := ctx.TypeLit().(*parser.TypeLitContext)
	if typ := typLit.ArrayType(); typ != nil {
		arrayType := typ.(*parser.ArrayTypeContext)
		len := arrayType.ArrayLength().GetText()
		elementType := arrayType.ElementType().(*parser.ElementTypeContext).Type_().(*parser.Type_Context)
		return fmt.Sprintf("[%s]%s", len, genTypeName(elementType))
	} else if typ := typLit.PointerType(); typ != nil {
		pointerType := typ.(*parser.PointerTypeContext)
		return fmt.Sprintf("*%s", genTypeName(pointerType.Type_().(*parser.Type_Context)))
	} else if typ := typLit.FunctionType(); typ != nil {
		functionType := typ.(*parser.FunctionTypeContext)
		return fmt.Sprintf("func%s", genSignature(functionType.Signature().(*parser.SignatureContext)))
	} else if typ := typLit.InterfaceType(); typ != nil {
		// xxx
	} else if typ := typLit.SliceType(); typ != nil {
		sliceType := typ.(*parser.SliceTypeContext)
		elementType := sliceType.ElementType().(*parser.ElementTypeContext).Type_().(*parser.Type_Context)
		return fmt.Sprintf("[]%s", genTypeName(elementType))
	} else if typ := typLit.MapType(); typ != nil {
		mapType := typ.(*parser.MapTypeContext)
		kType := mapType.Type_().(*parser.Type_Context)
		vType := mapType.ElementType().(*parser.ElementTypeContext).Type_().(*parser.Type_Context)
		return fmt.Sprintf("[%s]%s", genTypeName(kType), genTypeName(vType))
	} else if typ := typLit.ChannelType(); typ != nil {
		// xxx
	}

	return ""
}

func genInterface(interfaceName string, methods []*parser.MethodDeclContext) {
	buf := new(strings.Builder)
	buf.WriteString(fmt.Sprintf("type %sInterface interface {\n", interfaceName))

	for _, method := range methods {
		genMethod(buf, method)
	}

	buf.WriteString("}\n")

	if r, err := format.Source([]byte(buf.String())); err != nil {
		panic(err)
	} else {
		fmt.Println(string(r))
	}
}

func genMethod(buf *strings.Builder, method *parser.MethodDeclContext) {
	methodName := method.IDENTIFIER().GetText()
	firstRune := []rune(methodName)[0]
	if strings.ToUpper(string(firstRune)) != string(firstRune) {
		return
	}
	buf.WriteString(fmt.Sprintf("%s%s\n", methodName,
		genSignature(method.Signature().(*parser.SignatureContext))))
}

func genSignature(signature *parser.SignatureContext) string {
	buf := new(strings.Builder)
	buf.WriteString("(")
	paramterDecls := signature.Parameters().(*parser.ParametersContext).AllParameterDecl()
	for _, v := range paramterDecls {
		buf.WriteString(genTypeName(v.(*parser.ParameterDeclContext).Type_().(*parser.Type_Context)))
		buf.WriteString(",")
	}

	if ctx := signature.Result(); ctx != nil {
		result := ctx.(*parser.ResultContext)
		buf.WriteString(") (")
		if typ := result.Type_(); typ != nil {
			buf.WriteString(genTypeName(typ.(*parser.Type_Context)))
			buf.WriteString(",")
		} else {
			paramterDecls := result.Parameters().(*parser.ParametersContext).AllParameterDecl()
			for _, v := range paramterDecls {
				buf.WriteString(
					genTypeName(v.(*parser.ParameterDeclContext).Type_().(*parser.Type_Context)))
				buf.WriteString(",")
			}
		}
	}
	buf.WriteString(")\n")
	return buf.String()
}
