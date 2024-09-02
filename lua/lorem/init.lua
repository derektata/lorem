local Path = require "lorem.path"

--- Default configuration.
-- @type table<string, number>
local config = {
  WordsPerSentence = 10,
  SentencesPerParagraph = 5,
}

--- Filter options based on a prefix.
-- @param options table<string>: The list of options to filter.
-- @param prefix string: The prefix to filter by.
-- @return table<string>: The filtered options.
local function filter_opts(options, prefix)
  local filtered_opts = {}
  for _, opt in ipairs(options) do
    if opt:find("^" .. prefix) then
      table.insert(filtered_opts, opt)
    end
  end
  return filtered_opts
end

-- @module lorem
local M = {}

--- Build the Go binary and place it in the data directory.
-- @return nil
function M.build()
  local path = Path.bin_path()
  -- Go build command to compile the entire project and output the binary
  local build_cmd = string.format("go build -o %s/lorem", path)

  local handle = io.popen(build_cmd)

  if handle == nil then
    print "Failed to build lorem binary."
    return
  end

  handle:close()

  print "Successfully built lorem binary."
end

--- Setup function to configure the module.
-- @param user_config table: A table containing user configuration to override defaults.
-- @return nil
function M.setup(user_config)
  config = vim.tbl_extend("force", config, user_config or {})
  vim.env.WORDS_PER_SENTENCE = tostring(config.WordsPerSentence)
  vim.env.SENTENCES_PER_PARAGRAPH = tostring(config.SentencesPerParagraph)
end

--- Execute a command and return the output.
-- @param cmd string: The command to execute.
-- @return table: A list of lines from the command output.
function M.execute_command(cmd)
  local handle, err = io.popen(cmd, "r")
  if not handle then
    print("Failed to execute command: " .. err)
    return nil
  end

  local result = handle:read "*a"
  handle:close()

  local lines = {}
  for line in result:gmatch "[^\r\n]+" do
    table.insert(lines, line)
  end

  return lines
end

--- Generate lorem ipsum text based on type and amount.
-- @param textType string: The type of text to generate (e.g., "words", "paragraphs").
-- @param amount number: The amount of text to generate.
-- @return nil
function M.generate(textType, amount)
  textType = tostring(textType)
  amount = tonumber(amount)

  if type(amount) ~= "number" or amount <= 0 then
    print "Invalid amount. Expected a positive number."
    return
  end

  local path = Path.bin_path()
  local cmd = path .. " --" .. textType .. " " .. amount
  local lines = M.execute_command(cmd)

  if lines == nil then
    print "No output from command."
    return
  end

  vim.api.nvim_put(lines, "l", true, true)
end

--- Provide autocomplete suggestions for the LoremIpsum command.
-- @param arg_lead string: The leading part of the current argument.
-- @param cmd_line string: The entire command line.
-- @param cursor_pos number: The position of the cursor in the command line.
-- @return table<string>: A list of suggestions based on the input.
function M.lorem_complete(arg_lead, cmd_line, cursor_pos)
  local args = vim.split(cmd_line, "%s+")
  if #args == 2 then
    return filter_opts({ "words", "paragraphs" }, arg_lead)
  elseif #args == 3 then
    local textType = args[2]
    local options = {
      words = { "10", "20", "50", "100" },
      paragraphs = { "1", "2", "3", "5" },
    }
    return filter_opts(options[textType] or {}, arg_lead)
  end
end

--- Create the "LoremIpsum" user command.
-- @param opts table: The options table provided by Neovim for user commands.
-- @return nil
vim.api.nvim_create_user_command("LoremIpsum", function(opts)
  if #opts.fargs ~= 2 then
    print "Invalid number of arguments. Usage: LoremIpsum <words|paragraphs> <amount>"
    return
  end
  M.generate(opts.fargs[1], tonumber(opts.fargs[2]))
end, {
  nargs = "+",
  complete = M.lorem_complete,
})

return M

