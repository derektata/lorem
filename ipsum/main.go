// ipsum/main.go
package lorem

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	WordsPerSentence      = 10
	SentencesPerParagraph = 5
)

type Generator struct {
	Words                 []string
	WordsPerSentence      int
	SentencesPerParagraph int
	shuffleFunc           func([]string)
	r                     *rand.Rand
}

type Config struct {
	Words      int
	Paragraphs int
}

// getEnvInt returns the integer value of the environment variable
// specified by envVar. If the environment variable is not set, or its
// value is not a valid integer, it returns defaultValue.
func getEnvInt(envVar string, defaultValue int) int {
	if val, ok := os.LookupEnv(envVar); ok {
		if intValue, err := strconv.Atoi(val); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// NewGenerator returns a new Generator with default values set.
//
// The generator can be customized by setting the following environment variables:
//
// - WORDS_PER_SENTENCE: The number of words per sentence.
// - SENTENCES_PER_PARAGRAPH: The number of sentences per paragraph.
//
// If the environment variables are not set, the generator will use the default values.
func NewGenerator() *Generator {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Generator{
		Words:                 Words,
		WordsPerSentence:      getEnvInt("WORDS_PER_SENTENCE", 10),
		SentencesPerParagraph: getEnvInt("SENTENCES_PER_PARAGRAPH", 5),
		r:                     r,
		shuffleFunc: func(words []string) {
			r.Shuffle(len(words), func(i, j int) { words[i], words[j] = words[j], words[i] })
		},
	}
}

// processWord capitalizes the first word in the sentence.
//
// The function takes two parameters:
//
// - word: the word to process
// - wordCount: the position of the word in the sentence
//
// Return:
//
// - the processed word
func (g *Generator) processWord(word string, wordCount int) string {
	// Capitalize the first word in the sentence
	if wordCount == 1 {
		return cases.Title(language.English).String(word)
	}
	return word
}

// buildSentence generates a sentence by shuffling the words and adding
// punctuation.
//
// The function shuffles the words in the sentence before constructing the
// sentence. Then it builds the sentence by repeatedly calling processWord to add
// the words to the sentence. The function returns a string containing the
// generated sentence.
func (g *Generator) buildSentence(n int) string {
	// Shuffle words before constructing the sentence
	g.shuffleFunc(g.Words)

	var sb strings.Builder

	for i := 0; i < n; i++ {
		word := g.Words[i%len(g.Words)]
		word = g.processWord(word, i+1)

		if i > 0 {
			sb.WriteString(" ")
		}

		sb.WriteString(word)

		// Add a sporadic comma after the word if it's not the last word in the sentence
		if i < n-1 && g.r.Float64() < 0.3 { // 30% chance of adding a comma
			sb.WriteString(",")
		}
	}

	// Add a period at the end of the sentence
	sb.WriteString(".")
	return sb.String()
}

// buildParagraph generates a paragraph by repeatedly calling the buildSentence
// function to generate sentences. The sentences are then concatenated using a
// strings.Builder and returned as a single string.
//
// The function builds the paragraphs by repeatedly calling the buildSentence
// function to generate sentences. The sentences are then concatenated using a
// strings.Builder and returned as a single string.
func (g *Generator) buildParagraph() string {
	var sb strings.Builder
	for i := 0; i < g.SentencesPerParagraph; i++ {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(g.buildSentence(g.WordsPerSentence))
	}
	return sb.String()
}

// Generate generates n words of Lorem Ipsum text.
//
// The function builds the text by calling buildSentence with the provided number of words.
// The function returns a string containing the generated Lorem Ipsum text.
func (g *Generator) Generate(n int) string {
	return g.buildSentence(n)
}

// GenerateParagraphs generates n paragraphs of Lorem Ipsum text.
//
// The function builds the paragraphs by repeatedly calling the buildParagraph
// function to generate paragraphs. The paragraphs are then concatenated using a
// strings.Builder and returned as a single string.
func (g *Generator) GenerateParagraphs(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteString("\n\n")
		}
		sb.WriteString(g.buildParagraph())
	}
	return sb.String()
}

// GenerateIpsum generates Lorem Ipsum text based on the given configuration.
//
// The function takes a config struct as a parameter.
// The config struct contains the following fields:
//
// - paragraphs: the number of paragraphs to generate (if greater than 0).
//
// - words: the number of words to generate (if greater than 0).
//
// The function returns a string containing the generated Lorem Ipsum text.
func GenerateIpsum(c Config) string {
	g := NewGenerator()
	if c.Paragraphs > 0 {
		return g.GenerateParagraphs(c.Paragraphs)
	} else if c.Words > 0 {
		return g.Generate(c.Words)
	}
	return ""
}
