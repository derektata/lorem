-- shell script to grab the latest release from github and install it.
local install_script = [[

os_name="$(uname -s)"
os_arch="$(uname -m)"

case "$os_name" in
    "Darwin"|"Linux") latest_release_file="lorem_${os_name}_$os_arch.tar.gz";;
    "Windows") latest_release_file="lorem_${os_name}_$os_arch.zip";;
    *) echo "Unsupported OS: $os_name"; exit 1 ;;
esac

curl -s -L -o "$latest_release_file" $(curl -s https://api.github.com/repos/derektata/lorem/releases/latest \
    | grep 'browser_' \
    | cut -d\" -f4 \
    | grep "$os_name" \
    | grep "$os_arch")

case "$os_name" in
    "Darwin"|"Linux") tar -xzf "$latest_release_file" ;;
    "Windows") unzip "$latest_release_file" ;;
esac

rm "$latest_release_file"
]]

local M = {}

-- function M.bootstrap()
-- 	vim.fn.system(install_script)
-- end

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

return M
