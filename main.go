// main.go
package main

import (
	"fmt"

	lorem "github.com/derektata/lorem/ipsum"
	flag "github.com/spf13/pflag"
)

var (
	words      = flag.IntP("words", "w", 0, "number of words to generate")
	paragraphs = flag.IntP("paragraphs", "p", 0, "number of paragraphs to generate")
)

func main() {
	flag.Parse()

	c := lorem.Config{Words: *words, Paragraphs: *paragraphs}
	result := lorem.GenerateIpsum(c)
	fmt.Println(result)
}
