package parser

import "bytes"

const eof = rune(0)

func (p *Parser) scan() (tkn Token, literal string) {
	p.read()

	if isWhitespace(p.chr) {
		tkn, literal = p.scanWhitespace()
		return
	} else if isLetter(p.chr) {
		tkn, literal = p.scanIdentifier()
		return
	}

	switch p.chr {
	case eof:
		tkn, literal = EOF, ""
	case '#':
		tkn, literal = HASH, string(p.chr)
	case '"':
		tkn, literal = STRING, p.scanString()
	default:
		tkn, literal = UNKNOWN, string(p.chr)
	}

	return
}

func (p *Parser) read() {
	chr, _, err := p.in.ReadRune()
	if err != nil {
		p.chr = eof
	} else {
		p.chr = chr
	}
}

func (p *Parser) unread() {
	_ = p.in.UnreadRune()
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

func (p *Parser) scanWhitespace() (tkn Token, literal string) {
	var buf bytes.Buffer
	buf.WriteRune(p.chr)

	for {
		if p.read(); p.chr == eof {
			break
		} else if !isWhitespace(p.chr) {
			p.unread()
			break
		} else {
			buf.WriteRune(p.chr)
		}
	}

	tkn, literal = WHITESPACE, buf.String()
	return
}

func (p *Parser) scanIdentifier() (tkn Token, literal string) {
	var buf bytes.Buffer
	buf.WriteRune(p.chr)

	for {
		if p.read(); p.chr == eof {
			break
		} else if !isIdentifierRune(p.chr) {
			p.unread()
			break
		} else {
			buf.WriteRune(p.chr)
		}
	}

	tkn, literal = IDENTIFIER, buf.String()

	// is it a keyword
	tkn = checkKeyword(literal)

	return
}

func (p *Parser) scanString() string {
	var buf bytes.Buffer

	for {
		if p.read(); p.chr == eof || p.chr == '"' {
			break
		} else {
			buf.WriteRune(p.chr)
		}
	}

	return buf.String()
}
