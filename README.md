# tfcdash

A TUI for managing Terraform Cloud runs.

## Features
- List workspaces in your organization.
- View recent runs for a selected workspace.
- Approve/Apply runs that are pending approval.

## Requirements
- Go 1.21+
- Terraform CLI


## Installation
```bash
go build -o tfcdash main.go
```

## Usage
Run the binary:
```bash
./tfcdash
```

### Controls
- `Enter`: Select workspace/run.
- `Esc` or `Backspace`: Go back.
- `a`: Apply the selected run (in Run Detail view).
- `q` or `Ctrl+c`: Quit.

## Tech Stack
- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- [go-tfe](https://github.com/hashicorp/go-tfe)
