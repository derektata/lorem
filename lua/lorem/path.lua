--- Module to handle paths related to Neovim's data directory.
-- @module lorem.path

local M = {}

--- Get the standard data directory path for Neovim.
-- @return string: The path to Neovim's standard data directory.
function M.data_path()
  return vim.fn.stdpath "data"
end

--- Get the full path to the "lorem" binary in the data directory.
-- @return string: The path to the "lorem" binary within the data directory.
function M.bin_path()
  return M.data_path() .. "/lorem/bin/lorem"
end

return M

