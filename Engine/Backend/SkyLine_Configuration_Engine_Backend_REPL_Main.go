package SkyLine_Configuration_Engine_Backend_Source

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func Start(Filename string) {

	f, x := os.Open(Filename)
	if x != nil {
		Message := CallErrorStr(
			fmt.Sprint(SLC_FileSystem_ErrorWhenOpeningOrLoadingFile),
			"Could not open file due to -> "+fmt.Sprint(x),
			Filename,
		)
		println(Message)
		os.Exit(0)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var line []string
	for scanner.Scan() {
		line = append(line, scanner.Text())
	}
	if line == nil {
		Message := CallErrorStr(
			fmt.Sprint(SLC_FileSystem_NULL_FIELDS),
			"Engine refused to parse file due to there being nothing there",
			Filename+" -> NULL",
		)
		println(Message)
		os.Exit(0)
	}
	data, x := ioutil.ReadFile(Filename)
	if x != nil {
		Message := CallErrorStr(
			fmt.Sprint(SLC_FileSystem_ErrorWhenOpeningOrLoadingFile),
			"Could not open file due to -> "+fmt.Sprint(x),
			Filename,
		)
		println(Message)
		os.Exit(0)
	}
	parser := NewParser(New(string(data)))
	program := parser.ParseProgram()
	if len(parser.Errors()) > 0 {
		log.Fatal(parser.Errors()[0])
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
