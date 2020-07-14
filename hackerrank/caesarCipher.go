package main

func caesarCipher(s string, k int32) string {
	var cipher string
	// Iterate through s
	for _, ch := range s {
		shifted := shiftIfLetter(ch, k)
		cipher += shifted
	}
	// If character is a letter, shift the letter k times in place

	return cipher
}

func shiftIfLetter(ch rune, k int32) string {
	var shifted string
	switch {
	case ch >= 97 && ch <= 122:
		shifted = string(((ch-97)+k)%26 + 97)
	case ch >= 65 && ch <= 90:
		shifted = string(((ch-65)+k)%26 + 65)
	default:
		shifted = string(ch)
	}

	return shifted
}
