package main

// Returns the number of words in given camelcase string
func camelcase(s string) int32 {
	var count int32 = 1
	for _, ch := range s {
		if int(ch) >= 65 && int(ch) <= 90 {
			count++
		}
	}
	return count
}
