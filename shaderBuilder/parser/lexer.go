package parser

import "bytes"

const eof = rune(0)

func (self *Parser) scan() (tkn Token, literal string) {
	self.read()

	if isWhitespace(self.chr) {
		tkn, literal = self.scanWhitespace()
		return
	} else if isLetter(self.chr) {
		tkn, literal = self.scanIdentifier()
		return
	}

	switch self.chr {
	case eof:
		tkn, literal = EOF, ""
	case '#':
		tkn, literal = HASH, string(self.chr)
	default:
		tkn, literal = UNKNOWN, string(self.chr)
	}

	return
}

func (self *Parser) read() {
	chr, _, err := self.in.ReadRune()
	if err != nil {
		self.chr = eof
	} else {
		self.chr = chr
	}
}

func (self *Parser) unread() {
	_ = self.in.UnreadRune()
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isIdentifierRune(ch rune) bool {
	return isLetter(ch) || isDigit(ch) || ch == '_'
}

func (self *Parser) scanWhitespace() (tkn Token, literal string) {
	var buf bytes.Buffer
	buf.WriteRune(self.chr)

	for {
		if self.read(); self.chr == eof {
			break
		} else if !isWhitespace(self.chr) {
			self.unread()
			break
		} else {
			buf.WriteRune(self.chr)
		}
	}

	tkn, literal = WHITESPACE, buf.String()
	return
}

func (self *Parser) scanIdentifier() (tkn Token, literal string) {
	var buf bytes.Buffer
	buf.WriteRune(self.chr)

	for {
		if self.read(); self.chr == eof {
			break
		} else if !isIdentifierRune(self.chr) {
			self.unread()
			break
		} else {
			buf.WriteRune(self.chr)
		}
	}

	tkn, literal = IDENTIFIER, buf.String()

	// is it a keyword
	switch literal {
	case "INCLUDE":
		tkn = INCLUDE
	}

	return
}
