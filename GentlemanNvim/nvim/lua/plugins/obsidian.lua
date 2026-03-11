return {
  "obsidian-nvim/obsidian.nvim",
  version = "*", -- recommended, use latest release instead of latest commit
  lazy = false,
  enabled = function()
    -- Disable Obsidian when running from Oil Simple (to avoid path issues in Zed context)
    return not vim.g.disable_obsidian
  end,
  dependencies = {
    -- Required.
    "nvim-lua/plenary.nvim",
  },
  opts = function()
    -- Always include the personal vault
    local workspaces = {
      {
        name = "GentlemanNotes",
        path = os.getenv("HOME") .. "/.config/obsidian",
      },
    }

    -- Detect project vault: search upward from cwd for .obsidian-brain/
    local brain_dir = vim.fn.finddir(".obsidian-brain", vim.fn.getcwd() .. ";")
    if brain_dir ~= "" then
      local abs_path = vim.fn.fnamemodify(brain_dir, ":p")
      -- Derive project name from the parent directory of .obsidian-brain/
      local project_name = vim.fn.fnamemodify(abs_path:gsub("/$", ""), ":h:t")
      table.insert(workspaces, {
        name = project_name,
        path = abs_path,
      })
    end

    return {
      legacy_commands = false,
      workspaces = workspaces,
      detect_cwd = true,
      completion = {
        cmp = true,
      },
      picker = {
        -- Set your preferred picker. Can be one of 'telescope.nvim', 'fzf-lua', 'mini.pick' or 'snacks.pick'.
        name = "snacks.pick",
      },
      -- Optional, define your own callbacks to further customize behavior.
      callbacks = {
        -- Runs anytime you enter the buffer for a note.
        -- NOTE: Breaking change in obsidian.nvim - callback now receives only (note), not (client, note)
        enter_note = function(note)
          if not note then return end
          -- Setup keymaps for obsidian notes
          vim.keymap.set("n", "gf", function()
            return require("obsidian").util.gf_passthrough()
          end, { buffer = note.bufnr, expr = true, desc = "Obsidian follow link" })

          vim.keymap.set("n", "<leader>ch", function()
            return require("obsidian").util.toggle_checkbox()
          end, { buffer = note.bufnr, desc = "Toggle checkbox" })

          vim.keymap.set("n", "<cr>", function()
            return require("obsidian").util.smart_action()
          end, { buffer = note.bufnr, expr = true, desc = "Obsidian smart action" })
        end,
      },

      -- Settings for templates
      templates = {
        subdir = "templates", -- Subdirectory for templates
        date_format = "%Y-%m-%d-%a", -- Date format for templates
        gtime_format = "%H:%M", -- Time format for templates
        tags = "", -- Default tags for templates
      },
    }
  end,
}
