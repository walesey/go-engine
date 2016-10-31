package parser

type Token int

const (
	_ Token = iota

	EOF
	WHITESPACE
	IDENTIFIER

	HASH    //#
	UNKNOWN // eg. '{', '}', '+', ...

	// keywords
	INCLUDE
)
