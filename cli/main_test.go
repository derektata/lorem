package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateLorem(t *testing.T) {
	tests := []struct {
		name        string
		config      config
		expectEmpty bool
	}{
		{name: "Generate 10 Words", config: config{10, 0}, expectEmpty: false},
		{name: "Generate 1 Paragraph", config: config{0, 1}, expectEmpty: false},
		{name: "Zero Words and Paragraphs", config: config{0, 0}, expectEmpty: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := generateLorem(test.config)
			if test.expectEmpty {
				assert.Empty(t, result)
			} else {
				assert.NotEmpty(t, result)
			}
		})
	}
}

func TestRun(t *testing.T) {
	*words = 10
	*paragraphs = 2
	Run() // assert no panic during the running of the Run function
}
