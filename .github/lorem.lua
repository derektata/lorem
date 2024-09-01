-- plugins/lorem.lua
return {
	"derektata/lorem",
	branch = { "plugin" },
	dependencies = { "nvim-lua/plenary.nvim" }, -- Ensure Plenary is a dependency
	build = require("lorem").build(),
	config = function()
		require("lorem").setup({
			WordsPerSentence = 8, 		-- Default is 10
			SentencesPerParagraph = 4, 	-- Default is 5
		})

		-- Optionally, bind a key to generate lorem ipsum text
		vim.api.nvim_set_keymap("n", "<leader>li", ":LoremIpsum words 100<CR>", { noremap = true, silent = true })
	end,
	lazy = true,    	 -- Load the plugin lazily
	cmd = "LoremIpsum",  -- Load the plugin when the command is used
}
