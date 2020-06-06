package parser

// isWhitespace determines if a character is a whitespace character.
func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
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
