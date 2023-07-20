" API

function! LoremComplete(ArgLead, CmdLine, CursorPos)
    return luaeval('require("lorem").lorem_complete(_A)', a:ArgLead)
endfunction


" Interface

command! -nargs=* -complete=customlist,LoremComplete LoremIpsum lua require("lorem").generate(<f-args>)