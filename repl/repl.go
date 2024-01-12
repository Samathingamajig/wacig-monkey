package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Samathingamajig/wacig-monkey/compiler"
	"github.com/Samathingamajig/wacig-monkey/lexer"
	"github.com/Samathingamajig/wacig-monkey/parser"
	"github.com/Samathingamajig/wacig-monkey/vm"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Whoops! Compilation failed:\n %s\n", err)
			continue
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Whoops! Executing bytecode failed:\n %s\n", err)
			continue
		}

		stackTop := machine.StackTop()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Uh... we ran into some monkey business\n")
	io.WriteString(out, "Here are the parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t")
		io.WriteString(out, msg)
		io.WriteString(out, "\n")
	}
}
