package SkyLine_Configuration_Engine_Backend_Source

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func Start(Filename string) {

	f, x := os.Open(Filename)
	if x != nil {
		log.Fatal("Engine Fail (File): Could not open file due to -> ", x)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var line []string
	for scanner.Scan() {
		line = append(line, scanner.Text())
	}
	if line == nil {
		log.Fatal("Engine Fail (File): The engine did not FAIL but the engine did not find any lines of code to parse in this file, so the engine sent a PKILL signal")
	}
	data, x := ioutil.ReadFile(Filename)
	if x != nil {
		log.Fatal("Engine Fail (File): Could not open or read all of the file due to -> ", x)
	}
	parser := NewParser(New(string(data)))
	program := parser.ParseProgram()
	if len(parser.Errors()) > 0 {
		log.Fatal("Engine erorr: There was an error -> ", parser.Errors()[0])
	}

	Env := Start_Engine_Environment_Create()
	result := Eval(program, Env)
	if _, ok := result.(*ObjectNULL); ok {
		return
	}
	_, x = io.WriteString(os.Stdout, result.ObjectInspectFunc()+"\n")
	if x != nil {
		log.Fatal(x)
	}
}
