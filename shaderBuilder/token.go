package shaderBuilder

type Token int

const (
	_ Token = iota

	ILLEGAL
	EOF
	WHITESPACE
)