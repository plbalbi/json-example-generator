package main

import(
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	// "github.com/plbalbi/json-example-generator/model"
	"github.com/plbalbi/json-example-generator/parser"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fmt.Println("Hello Cats!")
	flag.Parse()
	var data []byte
	var err error
	switch flag.NArg() {
	case 0:
		fmt.Println("Parsing from stdin...")
		data, err = ioutil.ReadAll(os.Stdin)
		check(err)

		stringToParse := string(data)
		fmt.Println("String to be parsed:")
		fmt.Println(stringToParse)
		
		
		lexResult, lexError := parser.Parse(stringToParse)
		if lexError != nil {
			fmt.Printf("Errors: %v\n", lexError)	
		}
		
		fmt.Printf("Datatypes counted: %v\n", lexResult.StructsCount())
		fmt.Printf("First datatype seen: %v\n", lexResult.FirstDataTypeSeen())
		generatedStruct := lexResult.GenerateDataType()
		fmt.Println(generatedStruct)

		
		break
	case 1:
		data, err = ioutil.ReadFile(flag.Arg(0))
		check(err)
		fmt.Printf("file data: %v\n", string(data))
		break
	default:
		fmt.Printf("input must be from stdin or file\n")
		os.Exit(1)
	}
	
}
