package shaderBuilder

import (
	"io"
)

type Parser struct {
	in io.ByteReader
	
	token   Token // The token
	literal string      // The literal of the token, if any
	chr rune
}

// next moves pointer to next non-whitespace token
func (self *Parser) next() {
	for {
		self.token, self.literal = self.scan()
		if self.token != WHITESPACE {
			break
		}
	}
}