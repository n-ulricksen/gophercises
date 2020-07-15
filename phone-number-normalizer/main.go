package main

import "fmt"

func main() {
	fmt.Println("vim-go")
}

func normalize(phone string) string {
	var normalized string
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			normalized += string(ch)
		}
	}
	return normalized
}
