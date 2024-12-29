package parser

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type ParserValidationError struct {
	msg    string
	line   int
	column int
}

func (err *ParserValidationError) Error() string {
	return fmt.Sprintf("line %d column %d: %s", err.line, err.column, err.msg)
}

const (
	comma int = iota
	field
	lineBreak
	endOfText
)

const (
	asciiStartOfText = 2
	asciiEndOfText   = 3
)

type parser struct {
	reader           *bufio.Reader
	char             rune
	line             int
	column           int
	tokenKind        int
	tokenStartLine   int
	tokenStartColumn int
}

func (parser *parser) readNextChar() error {
	c, _, err := parser.reader.ReadRune()

	if err != nil {
		if err != io.EOF {
			return err
		}

		parser.char = asciiEndOfText
	} else {
		if c == '\n' {
			parser.line++
			parser.column = 1
		} else {
			parser.column++
		}

		parser.char = c
	}

	return nil
}

func (parser *parser) readNextToken() error {
	parser.tokenStartLine, parser.tokenStartColumn = parser.line, parser.column

	switch parser.char {
	case asciiEndOfText:
		parser.tokenKind = endOfText
		return nil
	case '\n':
		parser.tokenKind = lineBreak
		return parser.readNextChar()
	case ',':
		parser.tokenKind = comma
		return parser.readNextChar()
	case '"':
		for {
			parser.readNextChar()

			switch parser.char {
			case asciiEndOfText:
				return &ParserValidationError{
					msg:    "unclosed quoted field",
					line:   parser.line,
					column: parser.column,
				}
			case '"':
				err := parser.readNextChar()

				if err != nil {
					return err
				}

				if parser.char != '"' {
					parser.tokenKind = field
					return nil
				}
			}
		}
	default:
		for {
			err := parser.readNextChar()

			if err != nil {
				return err
			}

			if parser.char == ',' || parser.char == '\n' || parser.char == asciiEndOfText {
				parser.tokenKind = field
				return nil
			}
		}
	}
}

func (parser *parser) parseRecord() (int, error) {
	fieldCount := 0

	for {
		if parser.tokenKind == field {
			err := parser.readNextToken()

			if err != nil {
				return 0, err
			}
		}

		fieldCount++

		switch parser.tokenKind {
		case comma:
			err := parser.readNextToken()

			if err != nil {
				return 0, err
			}
		case lineBreak, endOfText:
			err := parser.readNextToken()

			if err != nil {
				return 0, err
			}

			return fieldCount, nil
		case field:
			return 0, &ParserValidationError{
				msg:    "expected comma or line break or end of text",
				line:   parser.tokenStartLine,
				column: parser.tokenStartColumn,
			}
		}
	}
}

func Validate(file *os.File) error {
	parser := &parser{
		reader: bufio.NewReader(file),
		char:   asciiStartOfText,
		line:   1,
		column: 0,
	}

	err := parser.readNextChar()

	if err != nil {
		return err
	}

	err = parser.readNextToken()

	if err != nil {
		return err
	}

	if parser.tokenKind == endOfText {
		return &ParserValidationError{
			msg:    fmt.Sprintf("expected at least 1 column"),
			line:   parser.tokenStartLine,
			column: parser.tokenStartColumn,
		}
	}

	fieldCountFirstLine, err := parser.parseRecord()

	if err != nil {
		return err
	}

	for {
		if parser.tokenKind == endOfText {
			return nil
		}

		err = parser.readNextToken()

		if err != nil {
			return err
		}

		fieldCount, err := parser.parseRecord()

		if err != nil {
			return err
		}

		if fieldCountFirstLine != fieldCount {
			return &ParserValidationError{
				msg:    fmt.Sprintf("expected %d fields, found %d", fieldCountFirstLine, fieldCount),
				line:   parser.tokenStartLine,
				column: parser.tokenStartColumn,
			}
		}
	}
}
