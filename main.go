package main

import (
	"flag"
	"fmt"
	"os"
	"toylanguage/lexer"
	"toylanguage/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

var dumpTokens = flag.Bool("dump-tokens", false, "print tokens")

var dumpAST = flag.Bool("dump-ast", false, "print ast")

func init() {
	if len(os.Args) <= 1 {
		os.Args = append(os.Args, "-h")
	}
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	fmt.Println("usage: toylanguage [flags] [fliename]")
	fmt.Println("------------------------------------")
	fmt.Println("flags:")
	flag.CommandLine.PrintDefaults()
}

func DumpTokens(goLexer *lexer.GoLexer) {
	for {
		token := goLexer.NextToken()
		fmt.Println(token)
		if token.GetTokenType() == antlr.TokenEOF {
			return
		}
	}
}

func main() {
	if len(flag.Args()) <= 0 {
		usage()
		return
	}
	fmt.Println(flag.Args())
	input, err := antlr.NewFileStream(flag.Args()[0])
	if err != nil {
		panic(err)
	}

	goLexer := lexer.NewGoLexer(input)
	if *dumpTokens {
		DumpTokens(goLexer)
		return
	}

	goParser := parser.NewGoParser(antlr.NewCommonTokenStream(goLexer, 0))
	goAST := goParser.SourceFile()

	if *dumpAST {
		text := goAST.ToStringTree(goLexer.RuleNames, goAST.GetParser())
		fmt.Println(text)
		return
	}
}
