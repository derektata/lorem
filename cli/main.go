package cli

import (
	"fmt"

	lorem "github.com/derektata/lorem/ipsum"
	"github.com/spf13/pflag"
)

var (
	words      = pflag.IntP("words", "w", 0, "number of words to generate")
	paragraphs = pflag.IntP("paragraphs", "p", 0, "number of paragraphs to generate")
)

type config struct {
	words      int
	paragraphs int
}

// generateLorem generates Lorem Ipsum text based on the given configuration.
//
// The function takes a config struct as a parameter.
// The config struct contains the following fields:
//
// - paragraphs: the number of paragraphs to generate (if greater than 0).
//
// - words: the number of words to generate (if greater than 0).
//
// The function returns a string containing the generated Lorem Ipsum text.
func generateLorem(c config) string {
	g := lorem.NewGenerator()
	g.WordsPerSentence = 10
	g.SentencesPerParagraph = 5
	g.CommaAddChance = 3

	if c.paragraphs > 0 {
		return g.GenerateParagraphs(c.paragraphs)
	} else if c.words > 0 {
		return g.Generate(c.words)
	}
	return ""
}

// Run generates Lorem Ipsum text based on the given configuration.
func Run() {
	pflag.Parse()

	c := config{*words, *paragraphs}
	result := generateLorem(c)
	fmt.Println(result)
}
