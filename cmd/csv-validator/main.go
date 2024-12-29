package main

import (
	"csv-validator/internal/parser"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	err := parser.Validate(os.Stdin)

	if err != nil {
		switch err.(type) {
		case *parser.ParserValidationError:
			log.Println(err)
			os.Exit(2)
		default:
			log.Fatal(err)
		}
	}

	log.Println("valid")
}
