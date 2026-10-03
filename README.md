# TerminalOS 🖥️

A secure, high-performance terminal emulator and shell interpreter designed for local use on Android (Termux), Linux, macOS, and Unix-like environments. The project focuses on correct shell syntax, clean command execution, strong security boundaries, and a modular library structure.

## Highlights

- Bash-like command parsing and execution
- Built-in commands for filesystem and text processing
- Environment-variable support
- Pipe and redirection handling
- Library/plugin registry for extension
- Firewall, honeypot, sandbox, and logging layers
- Local-first security model

## Security Warning

This project is intended for authorized, local, and educational use only.

Do not use it to:
- Access systems without permission
- Perform malware, phishing, or credential theft activities
- Launch network attacks or denial-of-service tools
- Bypass legitimate security protections or policy controls

If remote access is enabled, always enforce:
- Firewall rules
- Authentication
- Encryption (TLS/SSL)
- IP allowlisting
- Strong access logging
- External monitoring

## Requirements

- Go 1.22 or newer
- Linux / macOS / BSD / Android (Termux)
- Terminal support

## Quick Start

```bash
git clone https://github.com/Sunbr99/TerminalOS.git
cd TerminalOS
go build -o terminalos main.go
./terminalos
```

## Basic Usage

```bash
# Print working directory
pwd

# List files
ls
ls -la

# Change directory
cd /tmp

# View file contents
cat file.txt

# Search text
grep "pattern" file.txt

# Use variables
echo $HOME
NAME="demo"
echo $NAME

# Pipe output
ls | grep .go
cat file.txt | wc -l
```

## Project Structure

```text
TerminalOS/
├── main.go
├── go.mod
├── README.md
├── config/
│   ├── config.go
│   └── default_config.json
├── core/
│   ├── builtin.go
│   ├── command_parser.go
│   ├── environment.go
│   ├── executor.go
│   └── shell_engine.go
├── commands/
│   ├── file_commands.go
│   ├── system_commands.go
│   ├── text_commands.go
│   └── util_commands.go
├── library/
│   ├── loader.go
│   ├── plugin.go
│   └── stdlib.go
├── security/
│   ├── auth.go
│   ├── firewall.go
│   ├── honeypot.go
│   ├── logger.go
│   └── sandbox.go
├── examples/
│   ├── hello.sh
│   └── process_files.sh
├── tests/
│   └── shell_test.go
├── LICENSE
└── README.md
```

## Example Scripts

```bash
# Hello world
./examples/hello.sh

# File processing
./examples/process_files.sh
```

## Recommended Safe Mode

Keep the default configuration local-only unless you explicitly need a proxy or controlled access path.

```json
{
  "security": {
    "mode": "local_only",
    "firewall": {
      "enabled": true,
      "default_policy": "deny",
      "allowed_ips": ["127.0.0.1"]
    },
    "honeypot": {
      "enabled": true,
      "port": 2222
    }
  }
}
```

## License

MIT License.

## Disclaimer

This software is provided "as is" without warranty. Use only in authorized environments. The project maintainers are not responsible for any misuse, improper configuration, or unlawful activity.
