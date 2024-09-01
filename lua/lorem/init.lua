local Path = require "plenary.path"
local Job = require "plenary.job"
local unpack = table.unpack or unpack

local M = {}

-- Default configuration
local config = {
  WordsPerSentence = 10,
  SentencesPerParagraph = 5,
}

-- Setup function to configure the module
M.setup = function(user_config)
  -- Override defaults with user settings
  config = vim.tbl_extend("force", config, user_config or {})

  -- Set environment variables based on the configuration
  vim.env.WORDS_PER_SENTENCE = tostring(config.WordsPerSentence)
  vim.env.SENTENCES_PER_PARAGRAPH = tostring(config.SentencesPerParagraph)
end

--- Filter options based on a prefix.
-- @param options table<string>: The list of options to filter.
-- @param prefix string: The prefix to filter by.
-- @return table<string>: The filtered options.
local filter_opts = function(options, prefix)
  local filtered_opts = {}
  for _, opt in ipairs(options) do
    if opt:find("^" .. prefix) then
      table.insert(filtered_opts, opt)
    end
  end
  return filtered_opts
end

--- Generate lorem ipsum text based on type and amount.
-- @param textType string: The type of text to generate (e.g., "words", "paragraphs").
-- @param amount number: The amount of text to generate.
-- @return nil
M.generate = function(textType, amount)
  amount = tonumber(amount)
  textType = tostring(textType)

  if type(amount) ~= "number" or amount <= 0 then
    print "Invalid amount. Expected a positive number."
    return
  end

  local bin_path = Path:new(vim.fn.stdpath "data", "lorem")
  local command = bin_path:absolute() .. " --" .. textType .. " " .. amount

  Job:new({
    command = "sh",
    args = { "-c", command },
    on_exit = function(j, return_val)
      if return_val == 0 then
        local result = table.concat(j:result(), "\n")
        vim.schedule(function()
          local lines = {}
          for line in result:gmatch "[^\r\n]+" do
            table.insert(lines, line)
          end
          vim.api.nvim_put(lines, "l", true, true)
        end)
      else
        print "Error running lorem command"
      end
    end,
  }):start()
end

--- Provide autocomplete suggestions for the LoremIpsum command.
-- @param arg_lead string: The leading part of the argument being completed.
-- @param cmd_line string: The full command line.
-- @param cursor_pos number: The position of the cursor in the command line.
-- @return table<string>: The list of suggestions.
M.lorem_complete = function(arg_lead, cmd_line, cursor_pos)
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

--- Build the Go binary and place it in the data directory.
-- @return nil
M.build = function()
  local data_path = vim.fn.stdpath "data"

  -- Go build command to compile the entire project and output the binary
  local build_cmd = string.format("go build -o %s/lorem", data_path)

  -- Execute the Go build command using Plenary's Job module
  Job:new({
    command = "sh",
    args = { "-c", build_cmd },
    on_exit = function(j, return_val)
      if return_val == 0 then
        print("Go binary built successfully and placed in " .. data_path)
      else
        print "Failed to build Go binary"
      end
    end,
  }):start()
end

--- Create the "LoremIpsum" user command.
-- @param opts table: The options table provided by Neovim for user commands.
-- @return nil
vim.api.nvim_create_user_command("LoremIpsum", function(opts)
  if #opts.fargs ~= 2 then
    print "Invalid number of arguments. Usage: LoremIpsum <words|paragraphs> <amount>"
    return
  end
  M.generate(unpack(opts.fargs))
end, {
  nargs = "+",
  complete = M.lorem_complete,
})

-- Return the module for use in other scripts.
-- @return table: The module table containing all public functions.
return M
