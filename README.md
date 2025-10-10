# LinkedIn MCP Server

## Building

### macOS/Linux

To build the binary with proper execute permissions:

```bash
go build -o linkedin-mcp && chmod +x linkedin-mcp
```

This ensures the binary has execute permissions when distributed to other users.

### Windows

To build the binary for Windows:

```cmd
go build -o linkedin-mcp.exe
```

## Installing (macOS)

1. Download the release archive containing `linkedin-mcp` and `install-macos.command`.
2. Unzip and open the folder.
3. Double-click `install-macos.command` and follow the prompts.
   - The script copies the binary to `~/.local/bin`.
   - It removes Gatekeeper quarantine and sets execute permissions.
   - It updates `~/Library/Application Support/Claude/claude_desktop_config.json` with the MCP configuration.
4. Quit and relaunch Claude Desktop (Cmd+Q) to finish.

If you prefer the terminal, run:

```bash
~/Downloads/path-to-folder/install-macos.command
```

## Installing (Windows)

1. Download the release archive containing `linkedin-mcp.exe` and `install-windows.ps1`.
2. Extract the folder.
3. Run PowerShell and execute:

```powershell
Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass -Force
cd "C:\\path\\to\\folder"
./install-windows.ps1
```

   - The script copies the binary to `%LOCALAPPDATA%\LinkedInMCP`.
   - It removes the download block and updates `%APPDATA%\Claude\claude_desktop_config.json` with the MCP configuration.
4. Quit and relaunch Claude Desktop (Ctrl+Q) to finish.

## Troubleshooting

- **macOS warning about malware**: right-click the binary and choose Open once, or run `xattr -d com.apple.quarantine /full/path/to/linkedin-mcp`.
- **Windows execution policy**: if the script is blocked, ensure you used the `Set-ExecutionPolicy` command above.
- **Claude config reset**: rerun the installer for your platform to rewrite the configuration.
