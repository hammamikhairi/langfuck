package types

// TEMPORARILY FOR GO ONLY
type TokenKind uint8

const (
	TOKEN_INVALID TokenKind = iota
	TOKEN_PREPROC
	TOKEN_SYMBOL
	TOKEN_KEYWORD
	TOKEN_TYPE
	TOKEN_LIB
	TOKEN_COMMENT
	TOKEN_STRING
	TOKEN_TAB
	TOKEN_END
)