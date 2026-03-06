package ui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/go-tfe"
	"github.com/utkarsh867/tfcdash/internal/tfc"
	"github.com/utkarsh867/tfcdash/internal/ui/components"
)

type Model struct {
	client        *tfc.Client
	state         state
	organizations OrganizationListModel
	workspaces    WorkspaceListModel
	runs          RunsListModel
	runDetail     RunDetailModel
	header        components.HeaderModel
	err           error
	width         int
	height        int
	user          *tfe.User
}

func NewModel(client *tfc.Client) Model {
	o := NewOrganizationListModel(client)
	w := NewWorkspaceListModel(client)
	r := NewRunsListModel(client)
	d := NewRunDetailModel(client)
	h := components.NewHeaderModel("") // Org will be set later

	return Model{
		client:        client,
		state:         stateOrganizations,
		organizations: o,
		workspaces:    w,
		runs:          r,
		runDetail:     d,
		header:        h,
	}
}

func (m *Model) getUpdatedOrganizationListModel(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var orgModel tea.Model
	orgModel, cmd = m.organizations.Update(msg)
	m.organizations = orgModel.(OrganizationListModel)
	return cmd
}

func (m *Model) getUpdatedWorkspaceListModel(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var workspaceModel tea.Model
	workspaceModel, cmd = m.workspaces.Update(msg)
	m.workspaces = workspaceModel.(WorkspaceListModel)
	return cmd
}

func (m *Model) getUpdatedRunsListModel(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var runsModel tea.Model
	runsModel, cmd = m.runs.Update(msg)
	m.runs = runsModel.(RunsListModel)
	return cmd
}

func (m *Model) getUpdatedRunDetailModel(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var detailModel tea.Model
	detailModel, cmd = m.runDetail.Update(msg)
	m.runDetail = detailModel.(RunDetailModel)
	return cmd
}

func (m *Model) getUpdatedHeaderModel(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var headerModel tea.Model
	headerModel, cmd = m.header.Update(msg)
	m.header = headerModel.(components.HeaderModel)
	return cmd
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.organizations.Init(),
		m.fetchCurrentUser,
	)
}

func (m Model) fetchCurrentUser() tea.Msg {
	user, err := m.client.GetCurrentUser(context.Background())
	if err != nil {
		return errMsg(err)
	}
	return userMsg(user)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case selectOrganizationMsg:
		// Organization selected, set it in client and move to workspaces
		m.client.SetOrg(msg.organization.Name)
		m.header = m.header.SetOrg(msg.organization.Name)
		m.state = stateWorkspaces
		// Re-initialize workspaces to fetch for the selected org
		return m, m.workspaces.Init()

	case selectWorkspaceMsg:
		// Workspace selected, set it and fetch runs
		m.workspaces.selectedWS = msg.workspace
		m.state = stateRuns
		m.runs.SetWorkspace(msg.workspace)
		return m, m.runs.Init()

	case runsMsg:
		m.state = stateRuns
		cmd := m.getUpdatedRunsListModel(msg)
		return m, cmd

	case selectRunMsg:
		m.state = stateRunDetail
		cmd := m.getUpdatedRunDetailModel(msg)
		return m, cmd

	case navigateBackMsg:
		switch m.state {
		case stateRunDetail:
			m.state = stateRuns
		case stateRuns:
			m.state = stateWorkspaces
		case stateWorkspaces:
			m.state = stateOrganizations
			m.client.SetOrg("")
			m.header = m.header.SetOrg("")
		}
		return m, nil

	case runAppliedMsg:
		m.state = stateRuns
		return m, nil

	case userMsg:
		m.user = msg
		m.header = m.header.SetUser(msg)
		return m, nil

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.getUpdatedOrganizationListModel(msg)
		m.getUpdatedWorkspaceListModel(msg)
		m.getUpdatedRunsListModel(msg)
		m.getUpdatedRunDetailModel(msg)
		m.getUpdatedHeaderModel(msg)
		return m, nil

	case errMsg:
		m.err = msg
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	}

	switch m.state {
	case stateOrganizations:
		cmd := m.getUpdatedOrganizationListModel(msg)
		return m, cmd
	case stateWorkspaces:
		cmd := m.getUpdatedWorkspaceListModel(msg)
		return m, cmd
	case stateRuns:
		cmd := m.getUpdatedRunsListModel(msg)
		return m, cmd
	case stateRunDetail:
		cmd := m.getUpdatedRunDetailModel(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) currentTitle() string {
	switch m.state {
	case stateOrganizations:
		return "Organizations"
	case stateWorkspaces:
		return "Workspaces"
	case stateRuns:
		return "Runs"
	case stateRunDetail:
		return "Run Details"
	}
	return ""
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	// Update header title based on current state
	m.header = m.header.SetTitle(m.currentTitle())
	headerView := m.header.View()

	switch m.state {
	case stateOrganizations:
		return headerView + "\n" + m.organizations.View()
	case stateWorkspaces:
		return headerView + "\n" + m.workspaces.View()
	case stateRuns:
		return headerView + "\n" + m.runs.View()
	case stateRunDetail:
		return headerView + "\n" + m.runDetail.View()
	}
	return ""
}
