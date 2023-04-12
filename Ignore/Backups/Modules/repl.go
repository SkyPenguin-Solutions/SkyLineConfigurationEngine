package Rem

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := Start_Engine_Environment_Create()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := `
		ENGINE(true){
			INIT true {
				system("errors") -> ["debug=true"];
				set x  = 10;
				constant y = 20;
				puts("hello world")
			}
		}
		`
		var parse bool
		split := strings.Split(line, "\n")
		for _, k := range split {
			if strings.TrimSpace(k) == "ENGINE(true){" {
				parse = true
				break
			} else {
				parse = false
			}
		}
		if parse {
			l := New(line)
			p := NewParser(l)

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				continue
			}

			evaluated := Eval(program, env)
			if evaluated != nil {
				io.WriteString(out, evaluated.ObjectInspectFunc())
				io.WriteString(out, "\n")
			}
		} else {
			fmt.Println("[Engine] Warning: NULL fields, requires=engine{init{}}")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
