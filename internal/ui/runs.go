package ui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/go-tfe"
	"github.com/utkarsh867/tfcdash/internal/tfc"
	"github.com/utkarsh867/tfcdash/internal/ui/theme"
)

type RunsListModel struct {
	client      *tfc.Client
	workspace   *tfe.Workspace
	runs        []*tfe.Run
	runsTable   table.Model
	selectedRun *tfe.Run

	width, height int
}

func NewRunsListModel(client *tfc.Client) RunsListModel {
	return RunsListModel{
		client:    client,
		runsTable: initRunsTableBase(),
	}
}

func initRunsTableBase() table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 10},
		{Title: "Status", Width: 15},
		{Title: "Message", Width: 30},
		{Title: "Changes", Width: 15},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(5),
	)

	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		Foreground(theme.CurrentTheme.Text).
		Background(theme.CurrentTheme.Surface).
		Bold(true)
	styles.Selected = styles.Selected.
		Foreground(theme.CurrentTheme.Crust).
		Background(theme.CurrentTheme.Lavender).
		Bold(true)
	t.SetStyles(styles)
	return t
}

func (m RunsListModel) Init() tea.Cmd {
	if m.workspace == nil {
		return nil
	}
	return m.fetchRuns(m.workspace.ID)
}

func (m RunsListModel) fetchRuns(workspaceID string) tea.Cmd {
	return func() tea.Msg {
		runs, err := m.client.ListRuns(context.Background(), workspaceID)
		if err != nil {
			return errMsg(err)
		}
		return runsMsg(runs)
	}
}

func (m *RunsListModel) SetWorkspace(ws *tfe.Workspace) {
	m.workspace = ws
}

func (m RunsListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		hor, _ := tableContainerStyles.GetFrameSize()
		m.runsTable.SetWidth(msg.Width - hor)

		// Distribute column widths proportionally
		totalWidth := m.runsTable.Width()
		cols := m.runsTable.Columns()
		if len(cols) >= 4 {
			cols[0].Width = totalWidth / 8 // ID
			cols[1].Width = totalWidth / 6 // Status
			cols[2].Width = totalWidth / 3 // Message
			cols[3].Width = totalWidth / 6 // Changes
			for i := range cols {
				m.runsTable.Columns()[i] = cols[i]
			}
		}

		m.runsTable.UpdateViewport()
		m.runsTable, cmd = m.runsTable.Update(msg)
		return m, cmd

	case runsMsg:
		m.runs = msg
		var items []table.Row
		for _, r := range msg {
			status := string(r.Status)
			changes := ""
			if r.Plan != nil {
				changes = fmt.Sprintf("+%d ~%d -%d",
					r.Plan.ResourceAdditions,
					r.Plan.ResourceChanges,
					r.Plan.ResourceDestructions)
			}
			items = append(items, []string{
				r.ID,
				status,
				r.Message,
				changes,
			})
		}
		m.runsTable.SetRows(items)
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "esc" || msg.String() == "backspace" {
			return m, func() tea.Msg { return navigateBackMsg{} }
		} else if msg.String() == "enter" {
			i := m.runsTable.SelectedRow()
			if len(i) > 0 {
				for _, r := range m.runs {
					if r.ID == i[0] {
						m.selectedRun = r
						break
					}
				}
				return m, func() tea.Msg { return selectRunMsg{run: m.selectedRun} }
			}
		}
	}

	m.runsTable, cmd = m.runsTable.Update(msg)
	return m, cmd
}

func (m RunsListModel) View() string {
	return tableContainerStyles.Render(m.runsTable.View())
}
