# lorem

lorem is a Go utility library for generating Lorem Ipsum text.

## Todo List

- [X] - Build library code
  - [X] - Customize how many words per sentence
  - [X] - Customize how many sentences per paragraph
  - [X] - Customize the chance of a comma being added to a sentence
- [X] - Build simple CLI
- [ ] - Build Neovim Plugin
  - [ ] - Bootstrap binary for Neovim Plugin
  - [X] - Create ':LoremIpsum' command w/ args
  - [X] - Tab completion between words and paragraphs
  - [ ] - Integrate with completion engine
- [ ] - Write out unit tests for library
  - [X] - Write out tests for Generator
  - [X] - Write out tests for CLI
  - [ ] - Write out tests for Neovim Plugin

## Installation

To install the library, run the following command:

```
go get -u github.com/derektata/lorem
```

If you would like the CLI, grab the latest release from the [Latest Releases][Latest Release].

## Usage as a library

```go
package main

import (
    lorem "github.com/derektata/lorem/ipsum"
)

func main() {
    // Create a new generator
    g := lorem.NewGenerator()
    g.WordsPerSentence = 10 // Customize how many words per sentence
    g.SentencesPerParagraph = 5 // Customize how many sentences per paragraph

    // Generate words
    words := g.Generate(100)
    fmt.Println(words)

    // Generate paragraphs
    paragraphs := g.GenerateParagraphs(3)
    fmt.Println(paragraphs)
}
```

## Usage as a CLI

```bash
Usage of lorem:
  -p, --paragraphs int   number of paragraphs to generate
  -w, --words int        number of words to generate (default 100)

Examples:
  $ lorem -p 5
  $ lorem -w 500
```

## Usage as a Neovim Plugin

Make sure you have the binary from the [Latest Releases][Latest Release] in your PATH, then:

Installing the plugin via Lazy:
```lua
-- plugins/lorem.lua
return {
	"derektata/lorem",
	dependencies = { "nvim-lua/plenary.nvim" }, -- Ensure Plenary is a dependency
	build = require("lorem").build(),
	config = function()
		require("lorem").setup({
			WordsPerSentence = 8, 		  -- Default is 10
			SentencesPerParagraph = 4, 	-- Default is 5
		})

		-- Optionally, bind a key to generate lorem ipsum text
		vim.api.nvim_set_keymap("n", "<leader>li", ":LoremIpsum words 100<CR>", { noremap = true, silent = true })
	end,
	lazy = true, -- Load the plugin lazily
	cmd = "LoremIpsum", -- Load the plugin when the LoremIpsum command is used
}
```

Usage in Editor:
```
:LoremIpsum <tab>
```

<img src=".github/complete.png">

<img src=".github/lorem-plugin.gif">

## License

MIT License

[Latest Release]: https://github.com/derektata/lorem/releases