package parser

import (
	"bufio"
	"io"
)

type Parser struct {
	in *bufio.Reader

	fragOut, vertOut, geoOut io.Writer

	token   Token  // The token
	literal string // The literal of the token, if any
	chr     rune
}

func New(src io.Reader, frag, vert, geo io.Writer) *Parser {
	return &Parser{
		in:      bufio.NewReader(src),
		fragOut: frag,
		vertOut: vert,
		geoOut:  geo,
	}
}

func (self *Parser) Parse() {
	for {
		self.next()
		self.writeAll([]byte(self.literal))
		if self.token == EOF {
			break
		}
	}
}

func (self *Parser) next() {
	self.token, self.literal = self.scan()
}

func (self *Parser) writeAll(data []byte) {
	self.writeFrag(data)
	self.writeVert(data)
	self.writeGeo(data)
}

func (self *Parser) writeFrag(data []byte) {
	if self.fragOut != nil {
		self.fragOut.Write(data)
	}
}

func (self *Parser) writeVert(data []byte) {
	if self.vertOut != nil {
		self.vertOut.Write(data)
	}
}

func (self *Parser) writeGeo(data []byte) {
	if self.geoOut != nil {
		self.geoOut.Write(data)
	}
}
