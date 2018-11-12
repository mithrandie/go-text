package text

import "unicode"

// Calculates string width to be displayed.
func Width(s string, eastAsianEncoding bool, countDiacriticalSign bool) int {
	l := 0

	inEscSeq := false // Ignore ANSI Escape Sequence
	for _, r := range s {
		if inEscSeq {
			if unicode.IsLetter(r) {
				inEscSeq = false
			}
		} else if r == 27 {
			inEscSeq = true
		} else {
			l = l + RuneWidth(r, eastAsianEncoding, countDiacriticalSign)
		}
	}
	return l
}

// Calculates character width to be displayed.
func RuneWidth(r rune, eastAsianEncoding bool, countDiacriticalSign bool) int {
	switch {
	case unicode.IsControl(r):
		return 0
	case !countDiacriticalSign && unicode.In(r, ZeroWidthTable):
		return 0
	case unicode.In(r, FullWidthTable):
		return 2
	case eastAsianEncoding && unicode.In(r, AmbiguousTable):
		return 2
	}
	return 1
}

// Calculates byte size of a character.
func RuneByteSize(r rune, encoding Encoding) int {
	switch encoding {
	case SJIS:
		return sjisRuneByteSize(r)
	default:
		return len(string(r))
	}
}

func sjisRuneByteSize(r rune) int {
	switch {
	case unicode.In(r, SJISSingleByteTable) || unicode.IsControl(r):
		return 1
	}
	return 2
}

// Calculates byte size of a string.
func ByteSize(s string, encoding Encoding) int {
	size := 0
	switch encoding {
	case UTF8:
		size = len(s)
	default:
		for _, c := range s {
			size = size + RuneByteSize(c, encoding)
		}
	}
	return size
}

// Returns if a string is Right-to-Left horizontal writing characters.
func IsRightToLeftLetters(s string) bool {
	return 0 < len(s) && unicode.In([]rune(s)[0], RightToLeftTable)
}
