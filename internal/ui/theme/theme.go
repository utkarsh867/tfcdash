package theme

import (
	"os"

	"github.com/charmbracelet/lipgloss"
)

// Catppuccin color palette
type Theme struct {
	Name string

	// Base colors
	Base    lipgloss.Color
	Mantle  lipgloss.Color
	Crust   lipgloss.Color
	Surface lipgloss.Color
	Overlay lipgloss.Color

	// Text colors
	Text       lipgloss.Color
	Subtext0   lipgloss.Color
	Subtext1   lipgloss.Color
	TextMuted  lipgloss.Color
	TextBright lipgloss.Color

	// Accent colors
	Rosewater lipgloss.Color
	Flamingo  lipgloss.Color
	Pink      lipgloss.Color
	Mauve     lipgloss.Color
	Red       lipgloss.Color
	Maroon    lipgloss.Color
	Peach     lipgloss.Color
	Yellow    lipgloss.Color
	Green     lipgloss.Color
	Teal      lipgloss.Color
	Sky       lipgloss.Color
	Sapphire  lipgloss.Color
	Blue      lipgloss.Color
	Lavender  lipgloss.Color
}

// Catppuccin Mocha (dark, default)
var Mocha = Theme{
	Name:       "mocha",
	Base:       "#1e1e2e",
	Mantle:     "#181825",
	Crust:      "#11111b",
	Surface:    "#313244",
	Overlay:    "#45475a",
	Text:       "#cdd6f4",
	Subtext0:   "#a6adc8",
	Subtext1:   "#bac2de",
	TextMuted:  "#6c7086",
	TextBright: "#f5e0dc",
	Rosewater:  "#f5e0dc",
	Flamingo:   "#f2cdcd",
	Pink:       "#f5c2e7",
	Mauve:      "#cba6f7",
	Red:        "#f38ba8",
	Maroon:     "#eba0ac",
	Peach:      "#fab387",
	Yellow:     "#f9e2af",
	Green:      "#a6e3a1",
	Teal:       "#94e2d5",
	Sky:        "#89dceb",
	Sapphire:   "#74c7ec",
	Blue:       "#89b4fa",
	Lavender:   "#b4befe",
}

// Catppuccin Macchiato (dark)
var Macchiato = Theme{
	Name:       "macchiato",
	Base:       "#24273a",
	Mantle:     "#1e2030",
	Crust:      "#181926",
	Surface:    "#363a4f",
	Overlay:    "#494d64",
	Text:       "#cad3f5",
	Subtext0:   "#a5adcb",
	Subtext1:   "#b8c0ec",
	TextMuted:  "#6e738d",
	TextBright: "#f4dbd6",
	Rosewater:  "#f4dbd6",
	Flamingo:   "#f0c6c6",
	Pink:       "#f5bde6",
	Mauve:      "#c6a0f6",
	Red:        "#ed8796",
	Maroon:     "#ee99a0",
	Peach:      "#f5a97f",
	Yellow:     "#eed49f",
	Green:      "#a6da95",
	Teal:       "#8bd5ca",
	Sky:        "#91d7e3",
	Sapphire:   "#7dc4e4",
	Blue:       "#8aadf4",
	Lavender:   "#b7bdf8",
}

// Catppuccin Frappé (dark)
var Frappe = Theme{
	Name:       "frappe",
	Base:       "#303446",
	Mantle:     "#292c3c",
	Crust:      "#232634",
	Surface:    "#414559",
	Overlay:    "#51576d",
	Text:       "#c6d0f5",
	Subtext0:   "#a5adce",
	Subtext1:   "#b5bfe2",
	TextMuted:  "#737994",
	TextBright: "#f2d5cf",
	Rosewater:  "#f2d5cf",
	Flamingo:   "#eebebe",
	Pink:       "#f4b8e4",
	Mauve:      "#ca9ee6",
	Red:        "#e78284",
	Maroon:     "#ea999c",
	Peach:      "#ef9f76",
	Yellow:     "#e5c890",
	Green:      "#a6d189",
	Teal:       "#81c8be",
	Sky:        "#99d1db",
	Sapphire:   "#85c1dc",
	Blue:       "#8caaee",
	Lavender:   "#babbf1",
}

// Catppuccin Latte (light)
var Latte = Theme{
	Name:       "latte",
	Base:       "#eff1f5",
	Mantle:     "#e6e9ef",
	Crust:      "#dce0e8",
	Surface:    "#ccd0da",
	Overlay:    "#9ca0b0",
	Text:       "#4c4f69",
	Subtext0:   "#6c6f85",
	Subtext1:   "#5c5f77",
	TextMuted:  "#8c8fa1",
	TextBright: "#dc8a78",
	Rosewater:  "#dc8a78",
	Flamingo:   "#dd7878",
	Pink:       "#ea76cb",
	Mauve:      "#8839ef",
	Red:        "#d20f39",
	Maroon:     "#e64553",
	Peach:      "#fe640b",
	Yellow:     "#df8e1d",
	Green:      "#40a02b",
	Teal:       "#179299",
	Sky:        "#04a5e5",
	Sapphire:   "#209fb5",
	Blue:       "#1e66f5",
	Lavender:   "#7287fd",
}

// Current theme instance
var CurrentTheme = Mocha

func init() {
	// Check for theme override via environment variable
	if themeName := os.Getenv("TFCDASH_THEME"); themeName != "" {
		switch themeName {
		case "mocha":
			CurrentTheme = Mocha
		case "macchiato":
			CurrentTheme = Macchiato
		case "frappe":
			CurrentTheme = Frappe
		case "latte":
			CurrentTheme = Latte
		}
	}
}

// GetTheme returns a theme by name
func GetTheme(name string) Theme {
	switch name {
	case "mocha":
		return Mocha
	case "macchiato":
		return Macchiato
	case "frappe":
		return Frappe
	case "latte":
		return Latte
	default:
		return Mocha
	}
}

// SetTheme sets the current theme
func SetTheme(name string) {
	CurrentTheme = GetTheme(name)
}
