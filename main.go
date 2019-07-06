package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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
		fmt.Println("Parsing from stdin...")
		data, err = ioutil.ReadAll(os.Stdin)
		check(err)
		break
	case 1:
		fmt.Println("Parsing from file...")
		data, err = ioutil.ReadFile(flag.Arg(0))
		check(err)
		break
	default:
		fmt.Printf("input must be from stdin or file\n")
		os.Exit(1)
	}

	stringToParse := string(data)
	fmt.Println("String to be parsed:")
	fmt.Println(stringToParse)

	lexResult, lexError := parser.Parse(stringToParse)
	if lexError != nil {
		log.Printf("Errors: %v\n", lexError)

	} else {
		fmt.Printf("Datatypes counted: %v\n", lexResult.StructsCount())
		fmt.Printf("First datatype seen: %v\n", lexResult.FirstDataTypeSeen())
		generatedStruct := lexResult.GenerateDataType()
		generatedStruct, _ = jsonPrettyPrint(generatedStruct)
		fmt.Printf("%s", generatedStruct)
	}

}

func jsonPrettyPrint(jsonStr string) (string, error) {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(jsonStr), "", "  ")
	return out.String(), err
}
