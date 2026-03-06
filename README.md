# tfcdash

[![Go Version](https://img.shields.io/github/go-mod/go-version/utkarsh/tfcdash)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/utkarsh/tfcdash)](https://github.com/utkarsh/tfcdash/releases)

A terminal user interface (TUI) for managing Terraform Cloud and Terraform Enterprise runs directly from your terminal.

## Overview

tfcdash provides a fast, keyboard-driven interface to browse workspaces, view run history, and approve/apply pending runs without leaving your terminal.

Modern TUI experience. Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- Browse workspaces in your Terraform Cloud/Enterprise organization
- View recent runs with status, duration, and resource changes
- Approve and apply runs that require confirmation
- Keyboard-driven navigation
- Dark mode support (automatic based on terminal)

## Installation

### Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/utkarsh867/tfcdash/releases).

### Homebrew

```bash
brew install utkarsh/tfcdash/tfcdash
```

### Go Install

```bash
go install github.com/utkarsh/tfcdash@latest
```

### Build from Source

```bash
git clone https://github.com/utkarsh/tfcdash.git
cd tfcdash
go build -o tfcdash main.go
```

## Configuration

### Authentication

tfcdash uses the Terraform Cloud/Enterprise API. You need to set up authentication using Terraform CLI. Refer to the docs for [terraform login](https://developer.hashicorp.com/terraform/cli/commands/login)

## Usage

Run tfcdash:

```bash
./tfcdash
```

### Controls

| Key | Action |
|-----|--------|
| `Enter` | Select workspace/run |
| `Esc` / `Backspace` | Go back |
| `a` | Apply the selected run (in Run Detail view) |
| `q` / `Ctrl+c` | Quit |

## Tech Stack

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [go-tfe](https://github.com/hashicorp/go-tfe) - Terraform Cloud API client

## Contributing

Contributions are welcome! Please read our [contributing guidelines](CONTRIBUTING.md) before submitting PRs.

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Charm](https://charm.sh/) for the amazing TUI libraries
- [HashiCorp](https://www.hashicorp.com/) for Terraform Cloud
