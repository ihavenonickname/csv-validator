package main

import (
	"bufio"
	"csv-validator/internal/parser"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	err := parser.Validate(bufio.NewReader(os.Stdin))

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
