package ui

import (
	"context"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/go-tfe"
	"github.com/utkarsh867/tfcdash/internal/tfc"
	"github.com/utkarsh867/tfcdash/internal/ui/theme"
)

type OrganizationListModel struct {
	client *tfc.Client

	orgsTable   table.Model
	selectedOrg *tfe.Organization

	width, height int
}

func NewOrganizationListModel(client *tfc.Client) OrganizationListModel {
	return OrganizationListModel{
		client:    client,
		orgsTable: initOrgsTableBase(),
	}
}

func initOrgsTableBase() table.Model {
	columns := []table.Column{
		{Title: "Name", Width: 25},
		{Title: "Email", Width: 30},
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

func (m OrganizationListModel) fetchOrganizations() tea.Msg {
	orgs, err := m.client.ListOrganizations(context.Background())
	if err != nil {
		return errMsg(err)
	}
	return organizationsMsg(orgs)
}

func (m OrganizationListModel) Init() tea.Cmd {
	return m.fetchOrganizations
}

func (m OrganizationListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		hor, _ := tableContainerStyles.GetFrameSize()
		m.orgsTable.SetWidth(msg.Width - hor)
		cols := m.orgsTable.Columns()
		if len(cols) >= 2 {
			cols[0].Width = m.orgsTable.Width() / 2
			cols[1].Width = m.orgsTable.Width() / 2
			for i := range cols {
				m.orgsTable.Columns()[i] = cols[i]
			}
		}

		m.orgsTable.UpdateViewport()
		m.orgsTable, cmd = m.orgsTable.Update(msg)
		return m, cmd

	case organizationsMsg:
		var items []table.Row
		for _, org := range msg {
			items = append(items, []string{org.Name, org.Email})
		}
		m.orgsTable.SetRows(items)
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "enter" {
			i := m.orgsTable.SelectedRow()
			if len(i) > 0 {
				m.selectedOrg = &tfe.Organization{Name: i[0], Email: i[1]}
				return m, func() tea.Msg {
					return selectOrganizationMsg{organization: m.selectedOrg}
				}
			}
		}
	}

	m.orgsTable, cmd = m.orgsTable.Update(msg)
	return m, cmd
}

func (m OrganizationListModel) View() string {
	return tableContainerStyles.Render(m.orgsTable.View())
}
