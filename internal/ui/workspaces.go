package ui

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/go-tfe"
	"github.com/utkarsh867/tfcdash/internal/tfc"
	"github.com/utkarsh867/tfcdash/internal/ui/theme"
)

type WorkspaceListModel struct {
	client     *tfc.Client
	list       list.Model
	selectedWS *tfe.Workspace

	workspacesTable table.Model

	width, height int
}

func NewWorkspaceListModel(client *tfc.Client) WorkspaceListModel {
	workspacesList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	workspacesList.Title = "Workspaces"
	return WorkspaceListModel{
		client:          client,
		list:            workspacesList,
		workspacesTable: initWorkspaceTableBase(),
	}
}

func initWorkspaceTableBase() table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "Title", Width: 20},
		{Title: "Description", Width: 20},
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

func (m WorkspaceListModel) fetchWorkspaces() tea.Msg {
	ws, err := m.client.ListWorkspaces(context.Background())
	if err != nil {
		return errMsg(err)
	}
	return workspacesMsg(ws)
}

func (m WorkspaceListModel) fetchRuns(workspaceID string) tea.Cmd {
	return func() tea.Msg {
		runs, err := m.client.ListRuns(context.Background(), workspaceID)
		if err != nil {
			return errMsg(err)
		}
		return runsMsg(runs)
	}
}

func (m WorkspaceListModel) Init() tea.Cmd {
	return tea.Batch(
		m.fetchWorkspaces,
	)
}

func (m WorkspaceListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		hor, _ := tableContainerStyles.GetFrameSize()
		m.workspacesTable.SetWidth(msg.Width - hor)
		for i, col := range m.workspacesTable.Columns() {
			if i == 0 {
				col.Width = m.workspacesTable.Width() / 6
			}
			if i == 1 {
				col.Width = m.workspacesTable.Width() / 3
			}
			if i == 2 {
				col.Width = m.workspacesTable.Width() / 3
			}
			m.workspacesTable.Columns()[i] = col
		}

		m.workspacesTable.UpdateViewport()
		m.workspacesTable, cmd = m.workspacesTable.Update(msg)
		return m, cmd

	case workspacesMsg:
		var items []table.Row
		for _, ws := range msg {
			items = append(items, tablerow{id: ws.ID, title: ws.Name, desc: ws.Description}.toRow())
		}
		m.workspacesTable.SetRows(items)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace":
			return m, func() tea.Msg { return navigateBackMsg{} }
		case "enter":
			i := m.workspacesTable.SelectedRow()
			if len(i) > 0 {
				m.selectedWS = &tfe.Workspace{ID: i[0], Name: i[1]}
				return m, func() tea.Msg {
					return selectWorkspaceMsg{workspace: m.selectedWS}
				}
			}
		}
	}

	m.workspacesTable, cmd = m.workspacesTable.Update(msg)
	return m, cmd
}

func (m WorkspaceListModel) View() string {
	return tableContainerStyles.Render(m.workspacesTable.View())
}
