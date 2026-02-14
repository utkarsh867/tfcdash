package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hashicorp/go-tfe"
	"github.com/utkarsh/tfcdash/internal/ui/theme"
)

var (
	logoStyle = lipgloss.NewStyle().
			Foreground(theme.CurrentTheme.Blue).
			Bold(true).
			Padding(0, 1)

	headerTitleStyle = lipgloss.NewStyle().
				Foreground(theme.CurrentTheme.Crust).
				Background(theme.CurrentTheme.Lavender).
				Bold(true).
				Padding(0, 1)

	profileStyle = lipgloss.NewStyle().
			Foreground(theme.CurrentTheme.Subtext0).
			Padding(0, 1)

	profileOrgStyle = lipgloss.NewStyle().
			Foreground(theme.CurrentTheme.Green).
			Bold(true)

	profileUserStyle = lipgloss.NewStyle().
				Foreground(theme.CurrentTheme.Blue)
)

type HeaderModel struct {
	width int
	title string
	org   string
	user  *tfe.User
}

func NewHeaderModel(org string) HeaderModel {
	return HeaderModel{
		org: org,
	}
}

func (m HeaderModel) SetTitle(title string) HeaderModel {
	m.title = title
	return m
}

func (m HeaderModel) SetUser(user *tfe.User) HeaderModel {
	m.user = user
	return m
}

func (m HeaderModel) SetOrg(org string) HeaderModel {
	m.org = org
	return m
}

func (m HeaderModel) Init() tea.Cmd {
	return nil
}

func (m HeaderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}
	return m, nil
}

func (m HeaderModel) View() string {
	logo := logoStyle.Render("⛅ tfcdash")
	titleStr := headerTitleStyle.Render(m.title)

	var profileInfo string
	if m.user != nil {
		org := profileOrgStyle.Render(m.org)
		username := profileUserStyle.Render(m.user.Username)
		profileInfo = profileStyle.Render("Org: " + org + " | User: " + username)
	} else {
		profileInfo = profileStyle.Render("Org: " + profileOrgStyle.Render(m.org) + " | Loading user...")
	}

	// Calculate spacing for right alignment of profile info
	availableWidth := max(m.width-lipgloss.Width(logo)-lipgloss.Width(titleStr)-lipgloss.Width(profileInfo)-4, 0)
	spacing := strings.Repeat(" ", availableWidth)

	return lipgloss.JoinHorizontal(lipgloss.Top, logo, titleStr, spacing, profileInfo)
}
