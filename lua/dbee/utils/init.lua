local M = {}

-- submodules exposed through here
M.layout = require("dbee.utils.layout")
M.once = require("dbee.utils.once")

-- Get random key from table
---@param tbl table key-value table
---@return any|nil key
function M.random_key(tbl)
  -- luacheck: push ignore 512
  for k, _ in pairs(tbl) do
    return k
  end
  -- luacheck: pop
end

-- Get type from alias
---@param alias string
---@return string type
function M.type_alias(alias)
  local aliases = {
    ["postgresql"] = "postgres",
    ["pg"] = "postgres",
    ["sqlite3"] = "sqlite",
    ["mongodb"] = "mongo",
    ["duckdb"] = "duck",
  }
  return aliases[alias] or alias or ""
end

-- Get cursor range of current selection
---@return integer start row
---@return integer start column
---@return integer end row
---@return integer end column
function M.visual_selection()
  -- return to normal mode ('< and '> become available only after you exit visual mode)
  local key = vim.api.nvim_replace_termcodes("<esc>", true, false, true)
  vim.api.nvim_feedkeys(key, "x", false)

  local _, srow, scol, _ = unpack(vim.fn.getpos("'<"))
  local _, erow, ecol, _ = unpack(vim.fn.getpos("'>"))
  if ecol > 200000 then
    ecol = 20000
  end
  if srow < erow or (srow == erow and scol <= ecol) then
    return srow - 1, scol - 1, erow - 1, ecol
  else
    return erow - 1, ecol - 1, srow - 1, scol
  end
end

---@param level "info"|"warn"|"error"
---@param message string
---@param subtitle? string
function M.log(level, message, subtitle)
  -- log level
  local l = vim.log.levels.OFF
  if level == "info" then
    l = vim.log.levels.INFO
  elseif level == "warn" then
    l = vim.log.levels.WARN
  elseif level == "error" then
    l = vim.log.levels.ERROR
  end

  -- subtitle
  if subtitle then
    subtitle = "[" .. subtitle .. "]:"
  else
    subtitle = ""
  end
  vim.notify(subtitle .. " " .. message, l, { title = "nvim-dbee" })
end

-- Gets keys of a map and sorts them by name
---@param obj table<string, any> map-like table
---@return string[]
function M.sorted_keys(obj)
  local keys = {}
  for k, _ in pairs(obj) do
    table.insert(keys, k)
  end
  table.sort(keys)
  return keys
end

-- Get number of parameters that a function takes
---@param fun fun(...):any function to get the number of parameters of
---@return integer # number of parameters
function M.get_function_param_number(fun)
  local info = debug.getinfo(fun)
  if info == nil then
    return 0
  end

  return info.nparams or 0
end

return M
