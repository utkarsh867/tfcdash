package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/utkarsh/tfcdash/internal/config"
	"github.com/utkarsh/tfcdash/internal/credentials"
	"github.com/utkarsh/tfcdash/internal/tfc"
	"github.com/utkarsh/tfcdash/internal/ui"
	"github.com/utkarsh/tfcdash/internal/ui/theme"
)

func main() {
	cfg := config.Load()

	// Set theme if specified in config
	if cfg.Theme != "" {
		theme.SetTheme(cfg.Theme)
	}

	// Try to load token from Terraform CLI credentials
	creds, err := credentials.LoadDefault()
	if err != nil {
		fmt.Printf("Error loading Terraform CLI credentials: %v\n", err)
		fmt.Println("Please run 'terraform login' to authenticate.")
		os.Exit(1)
	}

	token, found := creds.GetDefaultToken()
	if !found {
		fmt.Println("Error: No Terraform Cloud credentials found.")
		fmt.Println("Please run 'terraform login' to authenticate with Terraform Cloud.")
		os.Exit(1)
	}

	client, err := tfc.NewClient(token)
	if err != nil {
		fmt.Printf("Error creating TFC client: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(ui.NewModel(client), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
