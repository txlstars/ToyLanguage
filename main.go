package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"toylanguage/ast"
	"toylanguage/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

var dumpTokens = flag.Bool("dump-tokens", false, "print tokens")

var dumpAST = flag.Bool("dump-ast", false, "print ast")

var dumpASTIndent = flag.Bool("dump-ast-indent", false, "print ast")

var genInterface = flag.Bool("gen-interface", false, "gen interface")

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

func DumpTokens(goLexer *parser.GoLexer) {
	for {
		token := goLexer.NextToken()
		fmt.Println(token)
		if token.GetTokenType() == antlr.TokenEOF {
			return
		}
	}
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	if len(flag.Args()) <= 0 {
		usage()
		return
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	input, err := antlr.NewFileStream(flag.Args()[0])
	// input, err := ioutil.ReadFile(flag.Args()[0])
	if err != nil {
		panic(err)
	}

	goLexer := parser.NewGoLexer(input)
	if *dumpTokens {
		DumpTokens(goLexer)
		return
	}

	goParser := parser.NewGoParser(antlr.NewCommonTokenStream(goLexer, 0))
	goAST := goParser.SourceFile()
	if *dumpAST {
		fmt.Println(antlr.TreesStringTree(goAST, []string{}, goParser))
		return
	}

	if *genInterface {
		walker := new(antlr.ParseTreeWalker)
		tool := ast.NewGenInterfaceTool()
		walker.Walk(tool, goAST)
		tool.CodeGen()
		return
	}

	if *dumpASTIndent {
		tool := ast.NewDumpASTTool()
		tool.Visit(goAST)
		tool.Dump()
		return
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

}
