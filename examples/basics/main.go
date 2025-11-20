package main

import (
	"fmt"
	"log"

	htmltomarkdown "github.com/lintstar/html-to-markdown/v2"
)

func main() {
	input := `<strong>Bold Text</strong>`

	markdown, err := htmltomarkdown.ConvertString(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(markdown)
	// Output: **Bold Text**
}
