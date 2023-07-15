package lorem

import (
	"math/rand"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	WordsPerSentence      = 10
	SentencesPerParagraph = 5
	CommaAddChance        = 2
)

type Generator struct {
	Words                 []string
	WordsPerSentence      int
	SentencesPerParagraph int
	CommaAddChance        int
	r                     *rand.Rand
}

// NewGenerator creates a new instance of the Generator struct.
//
// No parameters are required.
//
// It returns a pointer to a Generator.
func NewGenerator() *Generator {
	return &Generator{
		Words:                 Words,
		WordsPerSentence:      WordsPerSentence,
		SentencesPerParagraph: SentencesPerParagraph,
		CommaAddChance:        CommaAddChance,
		r:                     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// shuffleWords shuffles the words in the Generator.
//
// No parameters.
//
// No return values.
func (g *Generator) shuffleWords() {
	g.r.Shuffle(len(g.Words), func(i, j int) {
		g.Words[i], g.Words[j] = g.Words[j], g.Words[i]
	})
}

// capitalizeWord capitalizes the given word.
func (g *Generator) capitalizeWord(word string) string {
	word = cases.Title(language.English).String(word)
	return word
}

// addCommaToWord adds a comma to a word based on certain conditions.
//
// Parameters:
//
// - word: the word to which the comma is to be added (string)
//
// - wordCount: the count of words in the sentence (int)
//
// Returns:
//
// - a string containing the word with a comma added.
func (g *Generator) addCommaToWord(word string, wordCount int) string {
	if wordCount%g.WordsPerSentence != 0 && wordCount%g.WordsPerSentence < 5 && g.r.Intn(g.CommaAddChance) == 0 {
		word += ","
	}

	return word
}

// buildSentence generates a sentence by shuffling words and adding punctuation.
//
// Parameter:
//
// - n : is the length of the sentence in words.
//
// Return:
//
// - a string containing the generated sentence.
func (g *Generator) buildSentence(n int) string {
	g.shuffleWords()
	var sb strings.Builder

	for i := 0; i < n; i += g.WordsPerSentence {
		for j := 0; j < g.WordsPerSentence && i+j < n; j++ {
			word := g.Words[(i+j)%len(g.Words)]

			// Capitalize the first word of the sentence
			if j == 0 {
				word = g.capitalizeWord(word)
			}

			// Add a comma to the word
			word = g.addCommaToWord(word, j+1)

			sb.WriteString(" ")
			sb.WriteString(word)
		}

		if i < n-g.WordsPerSentence {
			sb.WriteString(". ")
		} else {
			sb.WriteString(".")
		}
	}

	return strings.TrimSpace(sb.String())
}

// buildParagraph generates a paragraph of text.
//
// It builds the paragraph by repeatedly calling the buildSentence function to
// generate sentences. The number of sentences per paragraph is determined by
// the SentencesPerParagraph field of the Generator struct that the function is
// called on. The sentences are then concatenated using a strings.Builder and
// returned as a single string.
//
// No parameters are required for this function.
//
// The function returns a string, which represents the generated paragraph.
func (g *Generator) buildParagraph() string {
	var sb strings.Builder

	for i := 0; i < g.SentencesPerParagraph; i++ {
		sentence := g.buildSentence(g.WordsPerSentence)

		sb.WriteString(" ")
		sb.WriteString(sentence)
	}

	return strings.TrimSpace(sb.String())
}

// Generate generates a sentence using the Generator struct.
//
// Parameter:
//
// - n: the number of words in the sentence.
//
// Return:
//
// - a string containing the generated sentence.
func (g *Generator) Generate(n int) string {
	// Build the sentence
	sentence := g.buildSentence(n)

	return sentence
}

// GenerateParagraphs generates n paragraphs and returns them as a string.
//
// Parameter:
//
// - n: the number of paragraphs to generate.
//
// Return:
//
// - a string containing the generated paragraphs.
func (g *Generator) GenerateParagraphs(n int) string {
	var sb strings.Builder

	for i := 0; i < n; i++ {
		// Build the paragraph
		paragraph := g.buildParagraph()

		// Add a space between paragraphs
		if i > 0 {
			sb.WriteString("\n\n")
		}

		sb.WriteString(paragraph)
	}

	return sb.String()
}
