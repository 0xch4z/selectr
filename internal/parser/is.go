package parser

// isWhitespace determines if a character is a whitespace character.
func isWhitespace(ch rune) bool {
	switch ch {
	case '\t', '\n', '\v', '\f', '\r', ' ', '\u0085', '\u00A0':
		return true
	}
	return false
}

// isLetter determines if a character is a letter character.
func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// isDigit determines if a character is a digit character.
func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isQoute(ch rune) bool {
	return ch == '\'' || ch == '"'
}
