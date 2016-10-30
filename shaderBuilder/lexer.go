package shaderBuilder

func isLineWhiteSpace(chr rune) bool {
	switch chr {
	case '\u0009', '\u000b', '\u000c', '\u0020', '\u00a0', '\ufeff':
		return true
	case '\u000a', '\u000d', '\u2028', '\u2029':
		return false
	case '\u0085':
		return false
	}
	return false
}

func (self *Parser) scan() (tkn Token, literal string) {
	
	return
}

func (self *Parser) read() {
	n, err := self.in.ReadByte()
	
}