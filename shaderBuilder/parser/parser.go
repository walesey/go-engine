package parser

import (
	"bufio"
	"io"
)

type Parser struct {
	in  *bufio.Reader
	out io.Writer

	token   Token  // The token
	literal string // The literal of the token, if any
	chr     rune
}

func New(src io.Reader, out io.Writer) *Parser {
	return &Parser{
		in:  bufio.NewReader(src),
		out: out,
	}
}

func (self *Parser) Parse() {
	for {
		self.next()
		self.out.Write([]byte(self.literal))
		if self.token == EOF {
			break
		}
	}
}

func (self *Parser) next() {
	self.token, self.literal = self.scan()
}
