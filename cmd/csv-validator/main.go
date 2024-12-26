package main

import (
	"bufio"
	"csv-validator/internal/parser"
	"log"
	"os"
)

func main() {
	f, err := os.Open(os.Args[1])

	if err != nil {
		log.Fatalf("Error opening the file: %s", err.Error())
	}

	defer f.Close()

	err = parser.Validate(bufio.NewReader(f))

	if err != nil {
		switch err.(type) {
		case *parser.ParserValidationError:
			log.Println(err)
			os.Exit(2)
		default:
			log.Fatal(err)
		}
	}

	log.Println("ok")
}
