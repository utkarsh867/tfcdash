package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/hashicorp/go-tfe"
	"github.com/utkarsh/tfcdash/internal/ui/theme"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)

	// Table container styles (shared across table views)
	tableContainerStyles = lipgloss.NewStyle().
				Inherit(docStyle).
				Border(lipgloss.RoundedBorder(), true).
				BorderForeground(theme.CurrentTheme.Surface).
				Padding(1, 2)

	// Header styles
	headerStyle = lipgloss.NewStyle().
			Foreground(theme.CurrentTheme.Text).
			Padding(0, 1).
			Width(100)
)

type item struct {
	id    string
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type tablerow struct {
	id    string
	title string
	desc  string
}

func (i tablerow) toRow() []string {
	return []string{
		i.id,
		i.title,
		i.desc,
	}
}

type state int

const (
	stateOrganizations state = iota
	stateWorkspaces
	stateRuns
	stateRunDetail
)

type workspacesMsg []*tfe.Workspace
type runsMsg []*tfe.Run
type runAppliedMsg struct{}
type userMsg *tfe.User
type selectRunMsg struct {
	run *tfe.Run
}
type selectWorkspaceMsg struct {
	workspace *tfe.Workspace
}
type selectOrganizationMsg struct {
	organization *tfe.Organization
}
type organizationsMsg []*tfe.Organization
type navigateBackMsg struct{}
type errMsg error
