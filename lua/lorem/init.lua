-- TODO: Booststrap Code

--[[

local getOSAndArch = function()
	local os_name = "Unknown OS"
	local os_arch = "Unknown Architecture"

	local capture_output = function(command)
		local file = io.popen(command)
		if file then
			local output = file:read("*l")
			file:close()
			return output
		end
		return nil
	end

	os_name = capture_output("uname -s") or os_name
	os_arch = capture_output("uname -m") or os_arch

	return os_name, os_arch
end

local os_name, os_arch = getOSAndArch()

local tag_name = "v0.0.1"

local url = "https://github.com/derektata/lorem/releases/download/" .. tag_name .. "/lorem_"

local getDownloadURL = function()
	if os_name == "Windows" or os_name == "Linux" or os_name == "Darwin" then
		local extension = (os_name == "Windows") and ".zip" or ".tar.gz"
		return url .. os_name .. "_" .. os_arch .. extension
		-- e.g. lorem_Darwin_arm64.tar.gz or lorem_Windows_x86_64.zip
	end
end

]]

local M = {}

--@param args table
--@param amount number
--@param textType string
--@return nil
function M.generate(...)
	local args = { ... }
	local amount = tonumber(args[1])
	local textType = tostring(args[2])

	if type(amount) ~= "number" or amount <= 0 then
		print("Invalid amount. Expected a positive number.")
		return
	end

	local handle, err = io.popen("lorem" .. " --" .. textType .. " " .. amount, "r")
	if not handle then
		print("Failed to execute command: " .. err)
		return
	end

	local result = handle:read("*a")
	handle:close()

	local lines = {}
	for line in result:gmatch("[^\r\n]+") do
		table.insert(lines, line)
	end

	vim.api.nvim_put(lines, "l", true, true)
end

--@param arg_lead string
--@return table
function M.lorem_complete(arg_lead)
	local opts = { "words", "paragraphs" }
	local filtered_opts = {}
	for _, opt in ipairs(opts) do
		if opt:find("^" .. arg_lead) then
			table.insert(filtered_opts, opt)
		end
	end
	return filtered_opts
end

vim.cmd([[
    function! LoremCompleteCmd(ArgLead, CmdLine, CursorPos)
        return luaeval('require("lorem").lorem_complete(_A)', a:ArgLead)
    endfunction

    command! -nargs=* -complete=customlist,LoremCompleteCmd LoremIpsum lua require("lorem").generate(<f-args>)
]])

return M
