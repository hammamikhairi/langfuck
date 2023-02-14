package lexer

import (
	. "LanguageFuck/Types"
	. "LanguageFuck/Utils"
)

type Lexer struct {
	Content      string
	Content_len  int
	Cursor       int
	Line         int
	KeywordsTree *map[string]uint8
}

func LexerInit(content string, tree *map[string]uint8) *Lexer {
	return &Lexer{content, len(content), 0, 0, tree}
}

func (l *Lexer) ChopChar(len int) {
	for i := 0; i < len; i++ {
		current := string(l.Content[l.Cursor])
		l.Cursor++
		if current == "\n" {
			l.Line++
		}
	}
}

func (l *Lexer) Trim() {
	for l.Cursor < l.Content_len-1 && (IsSpace(string(l.Content[l.Cursor])) || l.getCharAt(l.Cursor) == "\n") {
		l.ChopChar(1)
	}
}

func (l *Lexer) getCharAt(pos int) string {
	return string(l.Content[pos])
}

func (l *Lexer) startsWith(prefix string) bool {

	if len(prefix) == 1 {
		return l.getCharAt(l.Cursor) == prefix
	}

	for i := 0; i < len(prefix); i++ {
		if prefix[i] != l.Content[l.Cursor+1] {
			return false
		}
	}
	return true
}

func (l *Lexer) NextToken() *Token {
	l.Trim()
	l.getCharAt(l.Cursor)

	token := &Token{}
	token.Addr = Vec2i{X: l.Cursor, Line: l.Line}

	st := 0

	if l.Cursor >= l.Content_len {
		token.Kind = TOKEN_END
		token.Len = 1
		return token
	}

	if l.startsWith("\"") {
		token.Kind = TOKEN_STRING
		l.ChopChar(1)
		for l.Cursor < l.Content_len-1 {

			if l.getCharAt(l.Cursor) == "\"" && l.getCharAt(l.Cursor-1) != "\\" {
				break
			}

			l.ChopChar(1)
			st++
		}
		st += 2

		l.ChopChar(1)
		token.Len = st
		return token
	}

	if l.startsWith("\t") {
		token.Kind = TOKEN_TAB
		l.ChopChar(1)
		st++
		for l.getCharAt(l.Cursor) == "\t" {
			l.ChopChar(1)
			st++
		}
		token.Len = st
		return token
	}

	if l.startsWith("/") {
		l.ChopChar(1)
		if !l.startsWith("/") {
			l.Cursor--
			l.ChopChar(1)
			token.Kind = TOKEN_SYMBOL
			token.Len = 1
			return token
		}
		token.Kind = TOKEN_COMMENT
		for l.Cursor < l.Content_len && l.getCharAt(l.Cursor) != "\n" {
			st++
			l.ChopChar(1)
		}
		if l.Cursor < l.Content_len {
			st++
			l.ChopChar(1)
		}
		token.Len = st
		return token
	}

	if IsAlpha(l.getCharAt(l.Cursor)) {
		token.Kind = TOKEN_SYMBOL
		for l.Cursor < l.Content_len && IsSymbolChar(l.getCharAt(l.Cursor)) {
			l.ChopChar(1)
			st++
		}

		lastToken := l.Content[token.Addr.X : token.Addr.X+st]
		// PREPROC
		if lastToken == "package" {
			token.Kind = TOKEN_PREPROC
			for l.Cursor < l.Content_len && l.getCharAt(l.Cursor) != "\n" {
				st++
				l.ChopChar(1)
			}
			token.Len = st
			return token
		}

		// PREPROC
		if lastToken == "import" {
			token.Kind = TOKEN_PREPROC
			l.Trim()
			end := "\n"
			if l.startsWith("(") {
				end = ")"
			}

			for l.Cursor < l.Content_len && l.getCharAt(l.Cursor) != end {
				st++
				l.ChopChar(1)
			}

			if end == ")" {
				st += 1
				l.ChopChar(1)
			}

			st++
			l.ChopChar(1)

			token.Len = st
			return token
		}

		// KEYWORDS
		if val, ok := (*l.KeywordsTree)[lastToken]; ok {
			switch val {
			case 0:
				token.Kind = TOKEN_KEYWORD
			case 1:
				token.Kind = TOKEN_TYPE
			case 2:
				token.Kind = TOKEN_LIB
			}
		}

		if token.Kind == TOKEN_LIB {
			for l.Cursor < l.Content_len && (IsSymbolChar(l.getCharAt(l.Cursor)) || l.getCharAt(l.Cursor) == ".") {
				l.ChopChar(1)
				st++
			}
			token.Len = st
			return token
		}

		token.Len = st
		return token
	}

	l.ChopChar(1)
	token.Kind = TOKEN_INVALID
	token.Len = 1
	return token
}

func (l *Lexer) GetTokens() *[]*Token {
	tokens := []*Token{}
	for l.Cursor < l.Content_len {
		next := l.NextToken()
		tokens = append(tokens, next)
	}
	return &tokens
}

func (l *Lexer) GetTokenContent(token *Token) string {
	Assert(token.Addr.X+token.Len <= l.Content_len, "TOKEN OUT OF RANGE")
	return l.Content[token.Addr.X : token.Addr.X+token.Len]
}

func GetTokenName(tk TokenKind) string {
	switch tk {
	case TOKEN_INVALID:
		return "invalid token"
	case TOKEN_PREPROC:
		return "preprocessor directive"
	case TOKEN_SYMBOL:
		return "symbol"
	case TOKEN_KEYWORD:
		return "keyword"
	case TOKEN_TYPE:
		return "type"
	case TOKEN_LIB:
		return "lib"
	case TOKEN_COMMENT:
		return "comment"
	case TOKEN_STRING:
		return "string"
	case TOKEN_TAB:
		return "tabulation"
	case TOKEN_END:
		return "EOF"
	}
	return "UNREACHABLE"
}

// to lex multiple files
func (l *Lexer) ResetContent(content string) {
	l.Content = content
	l.Content_len = len(content)
	l.Cursor = 0
	l.Line = 0
}