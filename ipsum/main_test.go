package lorem

import (
	"strings"
	"testing"
)

// Test for NewGenerator
func TestNewGenerator(t *testing.T) {
	generator := NewGenerator()
	if generator == nil {
		t.Error("NewGenerator() should not return nil")
	}
}

// Test for shuffleWords
func TestShuffleWords(t *testing.T) {
	generator := NewGenerator()
	wordsBeforeShuffle := append([]string(nil), generator.Words...)
	generator.shuffleWords()
	if strings.Join(generator.Words, "") == strings.Join(wordsBeforeShuffle, "") {
		t.Error("shuffleWords() should shuffle the words")
	}
}

// Test for capitalizeWord
func TestCapitalizeWord(t *testing.T) {
	generator := NewGenerator()
	word := "test"
	capitalizedWord := generator.capitalizeWord(word)
	if capitalizedWord != "Test" {
		t.Errorf("capitalizeWord() should capitalize the first letter, got: %s, want: %s", capitalizedWord, "Test")
	}
}

// Test for addCommaToWord
func TestAddCommaToWord(t *testing.T) {
	generator := NewGenerator()
	word := "test"
	wordWithComma := generator.addCommaToWord(word, 3)
	if !strings.HasSuffix(wordWithComma, ",") && !strings.HasSuffix(wordWithComma, word) {
		t.Errorf("addCommaToWord() should potentially add comma at the end of the word, got: %s", wordWithComma)
	}
}

// Test for buildSentence
func TestBuildSentence(t *testing.T) {
	generator := NewGenerator()
	sentence := generator.buildSentence(10)
	words := strings.Split(sentence, " ")
	if len(words) != 10 {
		t.Errorf("buildSentence() should return a sentence of specified length, got: %d, want: %d", len(words), 10)
	}
}

// Test for buildParagraph
func TestBuildParagraph(t *testing.T) {
	generator := NewGenerator()
	paragraph := generator.buildParagraph()
	sentences := strings.Split(paragraph, ". ")
	if len(sentences) != SentencesPerParagraph {
		t.Errorf("buildParagraph() should return a paragraph with correct sentence count, got: %d, want: %d", len(sentences), SentencesPerParagraph)
	}
}

// Test for Generate
func TestGenerate(t *testing.T) {
	generator := NewGenerator()
	sentence := generator.Generate(10)
	words := strings.Split(sentence, " ")
	if len(words) != 10 {
		t.Errorf("Generate() should return a sentence of specified length, got: %d, want: %d", len(words), 10)
	}
}

// Test for GenerateParagraphs
func TestGenerateParagraphs(t *testing.T) {
	generator := NewGenerator()
	paragraphs := generator.GenerateParagraphs(2)
	paragraphList := strings.Split(paragraphs, "\n\n")
	if len(paragraphList) != 2 {
		t.Errorf("GenerateParagraphs() should return specified number of paragraphs, got: %d, want: %d", len(paragraphList), 2)
	}
}
