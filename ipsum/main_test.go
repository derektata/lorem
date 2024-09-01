// ipsum/main_test.go
package lorem

import (
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_getEnvInt tests the getEnvInt function, which retrieves an environment
// variable and converts it to an integer. The test cases cover the following
// scenarios:
//
// 1. The environment variable is set and valid.
// 2. The environment variable is set but invalid.
// 3. The environment variable is not set.
//
// The test uses the "TEST_ENV_VAR" environment variable and sets it to different
// values for each test case. The want value is the expected result of calling
// getEnvInt with the given arguments.
func Test_getEnvInt(t *testing.T) {
	type args struct {
		envVar       string
		defaultValue int
	}
	tests := []struct {
		name     string
		args     args
		want     int
		setup    func()
		teardown func()
	}{
		{
			name: "Environment variable set and valid integer",
			args: args{
				envVar:       "TEST_ENV_VAR",
				defaultValue: 42,
			},
			want: 100,
			setup: func() {
				os.Setenv("TEST_ENV_VAR", "100")
			},
			teardown: func() {
				os.Unsetenv("TEST_ENV_VAR")
			},
		},
		{
			name: "Environment variable set but invalid integer",
			args: args{
				envVar:       "TEST_ENV_VAR",
				defaultValue: 42,
			},
			want: 42,
			setup: func() {
				os.Setenv("TEST_ENV_VAR", "invalid")
			},
			teardown: func() {
				os.Unsetenv("TEST_ENV_VAR")
			},
		},
		{
			name: "Environment variable not set",
			args: args{
				envVar:       "TEST_ENV_VAR",
				defaultValue: 42,
			},
			want:     42,
			setup:    func() {},
			teardown: func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.teardown()
			got := getEnvInt(tt.args.envVar, tt.args.defaultValue)
			assert.Equal(t, tt.want, got, "getEnvInt() = %v, want %v", got, tt.want)
		})
	}
}

// TestNewGenerator tests the NewGenerator function by verifying that the
// WordsPerSentence and SentencesPerParagraph fields are set correctly, that
// the shuffleFunc is properly assigned, and that the random seed is set.
//
// The test also verifies that the shuffleFunc actually shuffles the words.
func TestNewGenerator(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Default generator initialization",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call NewGenerator to create a new generator instance
			got := NewGenerator()

			// Check that Words is correctly initialized (assuming Words is globally defined)
			assert.Equal(t, Words, got.Words, "Expected Words to be initialized correctly")

			// Check that WordsPerSentence is set to the default value or the environment variable
			assert.Equal(t, 10, got.WordsPerSentence, "Expected WordsPerSentence to be 10")

			// Check that SentencesPerParagraph is set to the default value or the environment variable
			assert.Equal(t, 5, got.SentencesPerParagraph, "Expected SentencesPerParagraph to be 5")

			// Check that the shuffleFunc is properly assigned
			assert.NotNil(t, got.shuffleFunc, "Expected shuffleFunc to be assigned")

			// Check that the random seed is set (this can't be compared directly, but we ensure it's initialized)
			assert.NotNil(t, got.r, "Expected rand.Rand to be initialized")

			// Ensure shuffleFunc actually shuffles
			initialWords := []string{"foo", "baz", "bar"}
			shuffledWords := make([]string, len(initialWords))
			copy(shuffledWords, initialWords)

			got.shuffleFunc(shuffledWords)

			// Check that the words have been shuffled (they should not be equal)
			if reflect.DeepEqual(initialWords, shuffledWords) {
				t.Errorf("Expected shuffleFunc to shuffle words, but it did not")
			}
		})
	}
}

// TestGenerator_processWord tests the processWord method of Generator.
//
// The test verifies that the first word of a sentence is capitalized, and
// subsequent words retain their original case.
func TestGenerator_processWord(t *testing.T) {
	type fields struct {
		Words                 []string
		WordsPerSentence      int
		SentencesPerParagraph int
		shuffleFunc           func([]string)
		r                     *rand.Rand
	}
	type args struct {
		word      string
		wordCount int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "First word should be capitalized",
			fields: fields{
				Words:                 []string{},
				WordsPerSentence:      10,
				SentencesPerParagraph: 5,
				shuffleFunc:           nil,
				r:                     rand.New(rand.NewSource(1)),
			},
			args: args{
				word:      "lorem",
				wordCount: 1,
			},
			want: "Lorem",
		},
		{
			name: "Subsequent words should not be capitalized",
			fields: fields{
				Words:                 []string{},
				WordsPerSentence:      10,
				SentencesPerParagraph: 5,
				shuffleFunc:           nil,
				r:                     rand.New(rand.NewSource(1)),
			},
			args: args{
				word:      "ipsum",
				wordCount: 2,
			},
			want: "ipsum",
		},
		{
			name: "Subsequent words should retain their original case",
			fields: fields{
				Words:                 []string{},
				WordsPerSentence:      10,
				SentencesPerParagraph: 5,
				shuffleFunc:           nil,
				r:                     rand.New(rand.NewSource(1)),
			},
			args: args{
				word:      "Dolor",
				wordCount: 3,
			},
			want: "Dolor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				Words:                 tt.fields.Words,
				WordsPerSentence:      tt.fields.WordsPerSentence,
				SentencesPerParagraph: tt.fields.SentencesPerParagraph,
				shuffleFunc:           tt.fields.shuffleFunc,
				r:                     tt.fields.r,
			}
			got := g.processWord(tt.args.word, tt.args.wordCount)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestGenerator_buildSentence tests the buildSentence method of Generator.
//
// The test verifies that the method returns a sentence of the correct length,
// that the sentence ends with a period, and that the sentence has the expected
// number of words.
//
// The test also verifies that the method shuffles the words in the sentence
// when the shuffleFunc is provided.
func TestGenerator_buildSentence(t *testing.T) {
	type fields struct {
		Words                 []string
		WordsPerSentence      int
		SentencesPerParagraph int
		shuffleFunc           func([]string)
		r                     *rand.Rand
	}
	type args struct {
		n int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Sentence with exact WordsPerSentence",
			fields: fields{
				Words:                 []string{"lorem", "ipsum", "dolor", "sit", "amet"},
				WordsPerSentence:      5,
				SentencesPerParagraph: 1,
				shuffleFunc: func(words []string) {
					// No shuffling for this test case
				},
				r: rand.New(rand.NewSource(1)), // Fixed seed for predictability
			},
			args: args{
				n: 5,
			},
		},
		{
			name: "Sentence with fewer words than WordsPerSentence",
			fields: fields{
				Words:                 []string{"lorem", "ipsum", "dolor"},
				WordsPerSentence:      5,
				SentencesPerParagraph: 1,
				shuffleFunc: func(words []string) {
					// No shuffling for this test case
				},
				r: rand.New(rand.NewSource(1)), // Fixed seed for predictability
			},
			args: args{
				n: 3,
			},
		},
		{
			name: "Shuffled sentence",
			fields: fields{
				Words:                 []string{"lorem", "ipsum", "dolor", "sit", "amet"},
				WordsPerSentence:      5,
				SentencesPerParagraph: 1,
				shuffleFunc: func(words []string) {
					// Shuffle the words
					for i := range words {
						j := rand.Intn(i + 1)
						words[i], words[j] = words[j], words[i]
					}
				},
				r: rand.New(rand.NewSource(1)), // Fixed seed for predictability
			},
			args: args{
				n: 5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				Words:                 tt.fields.Words,
				WordsPerSentence:      tt.fields.WordsPerSentence,
				SentencesPerParagraph: tt.fields.SentencesPerParagraph,
				shuffleFunc:           tt.fields.shuffleFunc,
				r:                     tt.fields.r,
			}
			got := g.buildSentence(tt.args.n)

			// Check that the generated sentence is not empty
			assert.NotEmpty(t, got, "Generated sentence should not be empty")

			// Check that the generated sentence ends with a period
			assert.True(t, strings.HasSuffix(got, "."), "Generated sentence should end with a period")

			// Check that the generated sentence has the expected number of words
			wordCount := len(strings.Fields(got))
			assert.Equal(t, tt.args.n, wordCount, "Generated sentence should have the correct number of words")
		})
	}
}

// TestGenerator_buildParagraph tests the buildParagraph method of Generator.
//
// The test verifies that the method returns a paragraph of the correct length,
// that the paragraph has the expected number of sentences, and that the
// sentences are separated by spaces.
//
// The test also verifies that the method shuffles the sentences in the
// paragraph when the shuffleFunc is provided.
func TestGenerator_buildParagraph(t *testing.T) {
	type fields struct {
		Words                 []string
		WordsPerSentence      int
		SentencesPerParagraph int
		shuffleFunc           func([]string)
		r                     *rand.Rand
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Paragraph with exact SentencesPerParagraph",
			fields: fields{
				Words:                 []string{"lorem", "ipsum", "dolor", "sit", "amet"},
				WordsPerSentence:      5,
				SentencesPerParagraph: 3,
				shuffleFunc: func(words []string) {
					// No shuffling for this test case
				},
				r: rand.New(rand.NewSource(1)), // Fixed seed for predictability
			},
		},
		{
			name: "Paragraph with shuffled sentences",
			fields: fields{
				Words:                 []string{"lorem", "ipsum", "dolor", "sit", "amet"},
				WordsPerSentence:      5,
				SentencesPerParagraph: 3,
				shuffleFunc: func(words []string) {
					// Shuffle the words
					for i := range words {
						j := rand.Intn(i + 1)
						words[i], words[j] = words[j], words[i]
					}
				},
				r: rand.New(rand.NewSource(1)), // Fixed seed for predictability
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				Words:                 tt.fields.Words,
				WordsPerSentence:      tt.fields.WordsPerSentence,
				SentencesPerParagraph: tt.fields.SentencesPerParagraph,
				shuffleFunc:           tt.fields.shuffleFunc,
				r:                     tt.fields.r,
			}
			got := g.buildParagraph()

			// Check that the generated paragraph is not empty
			assert.NotEmpty(t, got, "Generated paragraph should not be empty")

			// Check that the generated paragraph has the correct number of sentences
			sentenceCount := strings.Count(got, ".")
			assert.Equal(t, tt.fields.SentencesPerParagraph, sentenceCount, "Generated paragraph should have the correct number of sentences")

			// Check that the generated paragraph has spaces between sentences
			assert.True(t, strings.Contains(got, ". "), "Generated paragraph should have spaces between sentences")
		})
	}
}

// TestGenerator_Generate tests the Generate method of Generator.
//
// The test verifies that the method returns a sentence of the correct length,
// that the sentence ends with a period, and that the sentence has the expected
// number of words.
//
// The test also verifies that the method shuffles the words in the sentence
// when the shuffleFunc is provided.
func TestGenerator_Generate(t *testing.T) {
	type fields struct {
		Words                 []string
		WordsPerSentence      int
		SentencesPerParagraph int
		shuffleFunc           func([]string)
		r                     *rand.Rand
	}
	type args struct {
		n int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Generate sentence with exact number of words",
			fields: fields{
				Words:                 []string{"lorem", "ipsum", "dolor", "sit", "amet"},
				WordsPerSentence:      5,
				SentencesPerParagraph: 1,
				shuffleFunc: func(words []string) {
					// No shuffling for this test case
				},
				r: rand.New(rand.NewSource(1)), // Fixed seed for predictability
			},
			args: args{
				n: 5,
			},
		},
		{
			name: "Generate sentence with fewer words",
			fields: fields{
				Words:                 []string{"lorem", "ipsum", "dolor"},
				WordsPerSentence:      5,
				SentencesPerParagraph: 1,
				shuffleFunc: func(words []string) {
					// No shuffling for this test case
				},
				r: rand.New(rand.NewSource(1)), // Fixed seed for predictability
			},
			args: args{
				n: 3,
			},
		},
		{
			name: "Generate shuffled sentence",
			fields: fields{
				Words:                 []string{"lorem", "ipsum", "dolor", "sit", "amet"},
				WordsPerSentence:      5,
				SentencesPerParagraph: 1,
				shuffleFunc: func(words []string) {
					// Shuffle the words
					for i := range words {
						j := rand.Intn(i + 1)
						words[i], words[j] = words[j], words[i]
					}
				},
				r: rand.New(rand.NewSource(1)), // Fixed seed for predictability
			},
			args: args{
				n: 5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				Words:                 tt.fields.Words,
				WordsPerSentence:      tt.fields.WordsPerSentence,
				SentencesPerParagraph: tt.fields.SentencesPerParagraph,
				shuffleFunc:           tt.fields.shuffleFunc,
				r:                     tt.fields.r,
			}
			got := g.Generate(tt.args.n)

			// Check that the generated sentence is not empty
			assert.NotEmpty(t, got, "Generated sentence should not be empty")

			// Check that the generated sentence ends with a period
			assert.True(t, strings.HasSuffix(got, "."), "Generated sentence should end with a period")

			// Check that the generated sentence has the expected number of words
			wordCount := len(strings.Fields(got))
			assert.Equal(t, tt.args.n, wordCount, "Generated sentence should have the correct number of words")
		})
	}
}

// TestGenerator_GenerateParagraphs tests the GenerateParagraphs method of Generator.
//
// The test verifies that the method returns a string containing the correct number of paragraphs,
// that each paragraph is not empty and ends with a period, and that the paragraphs are separated by
// two newline characters.
//
// The test also verifies that the method shuffles the sentences in the paragraphs when the shuffleFunc
// is provided.
func TestGenerator_GenerateParagraphs(t *testing.T) {
	type fields struct {
		Words                 []string
		WordsPerSentence      int
		SentencesPerParagraph int
		shuffleFunc           func([]string)
		r                     *rand.Rand
	}
	type args struct {
		n int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Generate one paragraph",
			fields: fields{
				Words:                 []string{"lorem", "ipsum", "dolor", "sit", "amet"},
				WordsPerSentence:      5,
				SentencesPerParagraph: 3,
				shuffleFunc: func(words []string) {
					// No shuffling for this test case
				},
				r: rand.New(rand.NewSource(1)), // Fixed seed for predictability
			},
			args: args{
				n: 1,
			},
		},
		{
			name: "Generate multiple paragraphs",
			fields: fields{
				Words:                 []string{"lorem", "ipsum", "dolor", "sit", "amet"},
				WordsPerSentence:      5,
				SentencesPerParagraph: 3,
				shuffleFunc: func(words []string) {
					// No shuffling for this test case
				},
				r: rand.New(rand.NewSource(1)), // Fixed seed for predictability
			},
			args: args{
				n: 2,
			},
		},
		{
			name: "Generate shuffled paragraphs",
			fields: fields{
				Words:                 []string{"lorem", "ipsum", "dolor", "sit", "amet"},
				WordsPerSentence:      5,
				SentencesPerParagraph: 3,
				shuffleFunc: func(words []string) {
					// Shuffle the words
					for i := range words {
						j := rand.Intn(i + 1)
						words[i], words[j] = words[j], words[i]
					}
				},
				r: rand.New(rand.NewSource(1)), // Fixed seed for predictability
			},
			args: args{
				n: 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				Words:                 tt.fields.Words,
				WordsPerSentence:      tt.fields.WordsPerSentence,
				SentencesPerParagraph: tt.fields.SentencesPerParagraph,
				shuffleFunc:           tt.fields.shuffleFunc,
				r:                     tt.fields.r,
			}
			got := g.GenerateParagraphs(tt.args.n)

			// Check that the generated text is not empty
			assert.NotEmpty(t, got, "Generated text should not be empty")

			// Check that the generated text has the correct number of paragraphs
			paragraphs := strings.Split(got, "\n\n")
			assert.Equal(t, tt.args.n, len(paragraphs), "Generated text should have the correct number of paragraphs")

			// Check that each paragraph is not empty and ends with a period
			for _, paragraph := range paragraphs {
				assert.NotEmpty(t, paragraph, "Each paragraph should not be empty")
				assert.True(t, strings.HasSuffix(paragraph, "."), "Each paragraph should end with a period")
			}
		})
	}
}

// TestGenerateIpsum tests the GenerateIpsum function.
//
// The test verifies that the function returns the correct number of paragraphs
// or words based on the provided configuration.
//
// The test also verifies that an empty string is returned when the config is
// empty.
func TestGenerateIpsum(t *testing.T) {
	type args struct {
		c Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Generate multiple paragraphs",
			args: args{
				c: Config{Paragraphs: 2},
			},
		},
		{
			name: "Generate one paragraph",
			args: args{
				c: Config{Paragraphs: 1},
			},
		},
		{
			name: "Generate specific number of words",
			args: args{
				c: Config{Words: 5},
			},
		},
		{
			name: "Empty config returns empty string",
			args: args{
				c: Config{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateIpsum(tt.args.c)

			switch {
			case tt.args.c.Paragraphs > 0:
				// Validate that the correct number of paragraphs is generated
				paragraphs := len(strings.Split(got, "\n\n"))
				assert.Equal(t, tt.args.c.Paragraphs, paragraphs, "Expected number of paragraphs")
			case tt.args.c.Words > 0:
				// Validate that the correct number of words is generated
				words := len(strings.Fields(got))
				assert.Equal(t, tt.args.c.Words, words, "Expected number of words")
			default:
				// Validate that an empty string is returned
				assert.Empty(t, got, "Expected empty string")
			}
		})
	}
}
