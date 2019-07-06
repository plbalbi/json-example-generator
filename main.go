package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/plbalbi/json-example-generator/parser"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	flag.Parse()
	var data []byte
	var err error
	switch flag.NArg() {
	case 0:
		data, err = ioutil.ReadAll(os.Stdin)
		check(err)
		break
	case 1:
		data, err = ioutil.ReadFile(flag.Arg(0))
		check(err)
		break
	default:
		fmt.Fprintf(os.Stderr, "input must be from stdin or file\n")
		os.Exit(1)
	}

	stringToParse := string(data)

	lexResult, lexError := parser.Parse(stringToParse)
	if lexError != nil {
		fmt.Fprintf(os.Stderr, "Parsing failed: %s\n", lexError.Error())
		os.Exit(1)
	} else {
		generatedStruct := lexResult.GenerateDataType()
		generatedStruct, _ = jsonPrettyPrint(generatedStruct)
		fmt.Printf("%s\n", generatedStruct)
		os.Exit(0)
	}
}

func jsonPrettyPrint(jsonStr string) (string, error) {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(jsonStr), "", "  ")
	return out.String(), err
}
