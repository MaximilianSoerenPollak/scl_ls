# scl_ls
Simple language server exploration for a previously build sphinx extension called 'source code linker'.

## What is this?
This is a small project to learn more about LSP's and maybe implement it into a real world use.

Currently this is the MVP with very simple logic and internals. Missing tests etc.

## How to use it? 

Currently I only tested it crudely on Neovim (0.10+). Though plugins/integration for VSCode and Neovim are planned.

Here is how you can implement it in neovim. 
1. Clone The Repo
```bash
git clone git@github.com:MaximilianSoerenPollak/scl_ls.git
```
2. Inside your init.lua or whever you load your configuration add the following: 
```lua
local client = vim.lsp.start_client {
  name = "sclls",
  cmd = { "<path to the binary>" },
  priority = 1, -- low priority as to not disturb other 'real' LSP's
}
if not client then
  vim.notify("didn't do client thing good")
  return
end

vim.api.nvim_create_autocmd("FileType", {
  pattern = '*', -- Change this to any file type you want, or keep it as is to enable the LSP in all files
  callback = function()
    vim.lsp.buf_attach_client(0, client)
  end,
})
```


## What can it do? 

### Diagnostics
It can publish diagnostics as errors or warnings. It looks like this: 
![](_assets/diagnostis_prev.png)

### Go To Definition
If you have a 'need' it knows defined, it can go to the definition of said need inside of your sphinx documentation (rst files)

### Completion
It has completion suggestions for template strings and needs it knows (from the needs.json)

### Logging
Very simple logging inside the LSP.

## What Is Missing? 

Currently there is quiet some things missing I would like to integrate into this LSP.

- [ ] More & better tests
- [ ] Better Documentation (of everything, inside and outside the code)
- [ ] VSCode integration (via a plugin)
- [ ] Better Neovim integration (plugin?)
- [ ] Persitent Datastorage (maybe sqlite3 db or so?)
- [ ] Debouncing of spaming messages (Diagnostics mainly)

- [ ] Further improvements based on feedback


