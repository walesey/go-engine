package parser

type Token int

const (
	_ Token = iota

	EOF
	WHITESPACE
	IDENTIFIER
	STRING
	NUMBER

	LEFT_PARENTHESIS  // (
	RIGHT_PARENTHESIS // )

	HASH      // #
	PLUS      // +
	MINUS     // -
	MULTIPLY  // *
	SLASH     // /
	REMAINDER // %
	UNKNOWN

	// keywords
	INCLUDE
	VERT
	FRAG
	GEO
	ENDVERT
	ENDFRAG
	ENDGEO
	LOOKUP
)

func checkKeyword(literal string) Token {
	switch literal {
	case "include":
		return INCLUDE
	case "vert":
		return VERT
	case "frag":
		return FRAG
	case "geo":
		return GEO
	case "endvert":
		return ENDVERT
	case "endfrag":
		return ENDFRAG
	case "endgeo":
		return ENDGEO
	case "lookup":
		return LOOKUP
	}
	return IDENTIFIER
}
