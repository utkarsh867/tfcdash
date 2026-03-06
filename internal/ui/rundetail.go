package ui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hashicorp/go-tfe"
	"github.com/hashicorp/terraform-json"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/utkarsh867/tfcdash/internal/tfc"
)

type planDetailsMsg struct {
	plan *tfjson.Plan
	err  error
}

type applyLogsMsg struct {
	status string
	err    error
}

// Style definitions
var (
	runDetailTitleStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#89b4fa"))
	runDetailLabelStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086"))
	runDetailValueStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4"))
	runDetailCreateStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#40a02b"))
	runDetailUpdateStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#df8e1d"))
	runDetailDeleteStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#d20f39"))
	runDetailReplaceStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#d20f39"))
	runDetailHelpStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086"))
	runDetailSeparatorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#45475a"))
)

type RunDetailModel struct {
	client       *tfc.Client
	run          *tfe.Run
	planJSON     *tfjson.Plan
	applyLog     string
	loadingPlan  bool
	loadingApply bool
	width        int
	height       int
}

func NewRunDetailModel(client *tfc.Client) RunDetailModel {
	return RunDetailModel{
		client: client,
	}
}

func (m *RunDetailModel) SetRun(run *tfe.Run) {
	m.run = run
	m.planJSON = nil
	m.applyLog = ""
	m.loadingPlan = false
	m.loadingApply = false
}

func (m RunDetailModel) Init() tea.Cmd {
	return nil
}

func (m RunDetailModel) fetchPlanDetails() tea.Msg {
	if m.run == nil || m.run.Plan == nil {
		return planDetailsMsg{err: fmt.Errorf("no plan available")}
	}

	jsonData, err := m.client.GetPlanJSONOutput(context.Background(), m.run.Plan.ID)
	if err != nil {
		return planDetailsMsg{err: err}
	}

	var plan tfjson.Plan
	if err := json.Unmarshal(jsonData, &plan); err != nil {
		return planDetailsMsg{err: err}
	}

	return planDetailsMsg{plan: &plan}
}

func (m RunDetailModel) fetchApplyLogs() tea.Msg {
	if m.run == nil || m.run.Apply == nil {
		return applyLogsMsg{err: fmt.Errorf("no apply available")}
	}

	status, err := m.client.GetApplyLogs(context.Background(), m.run.Apply.ID)
	if err != nil {
		return applyLogsMsg{err: err}
	}

	return applyLogsMsg{status: status}
}

func (m RunDetailModel) applyRun(runID string) tea.Cmd {
	return func() tea.Msg {
		err := m.client.ApplyRun(context.Background(), runID, "Applied via tfcdash")
		if err != nil {
			return errMsg(err)
		}
		return runAppliedMsg{}
	}
}

func (m RunDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case selectRunMsg:
		m.SetRun(msg.run)
		m.loadingPlan = true
		m.loadingApply = false

		var cmds []tea.Cmd
		cmds = append(cmds, m.fetchPlanDetails)

		if m.run != nil && m.run.Apply != nil {
			m.loadingApply = true
			cmds = append(cmds, m.fetchApplyLogs)
		}
		return m, tea.Batch(cmds...)

	case planDetailsMsg:
		m.loadingPlan = false
		if msg.err == nil {
			m.planJSON = msg.plan
		}
		return m, nil

	case applyLogsMsg:
		m.loadingApply = false
		if msg.err == nil {
			m.applyLog = msg.status
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace":
			return m, func() tea.Msg { return navigateBackMsg{} }
		case "a":
			if m.run != nil && (m.run.Status == tfe.RunPlanned || m.run.Status == tfe.RunPlannedAndFinished) {
				return m, m.applyRun(m.run.ID)
			}
		}
	}
	return m, nil
}

func (m RunDetailModel) renderResourceChangeCategory(title string, resources []*tfjson.ResourceChange, style lipgloss.Style, prefix string) []string {
	if len(resources) == 0 {
		return nil
	}

	var section []string
	// Category header
	header := lipgloss.JoinHorizontal(lipgloss.Left,
		style.Render(fmt.Sprintf("  %s %s", prefix, title)),
		runDetailValueStyle.Render(fmt.Sprintf(" (%d)", len(resources))),
	)
	section = append(section, header)

	// Items with detailed changes
	for _, rc := range resources {
		resourceSection := m.renderResourceChange(rc, style, prefix)
		section = append(section, resourceSection...)
	}

	return section
}

func (m RunDetailModel) renderResourceChange(rc *tfjson.ResourceChange, style lipgloss.Style, prefix string) []string {
	var sections []string

	// Resource header line
	resourceHeader := fmt.Sprintf("    %s %s (%s)", prefix, rc.Address, rc.Type)
	sections = append(sections, style.Render(resourceHeader))

	// If there's no change data, return just the header
	if rc.Change == nil {
		return sections
	}

	// For updates and replaces, show field-level changes
	if rc.Change.Actions.Update() || rc.Change.Actions.Replace() {
		fieldChanges := m.extractFieldChanges(rc)
		for _, change := range fieldChanges {
			// Indent field changes more
			fieldLine := fmt.Sprintf("      %s", change)
			sections = append(sections, runDetailLabelStyle.Render(fieldLine))
		}
	}

	return sections
}

func (m RunDetailModel) extractFieldChanges(rc *tfjson.ResourceChange) []string {
	var changes []string

	if rc.Change == nil || rc.Change.Before == nil || rc.Change.After == nil {
		return changes
	}

	// Before and After are interface{}, need to marshal then unmarshal
	beforeBytes, err := json.Marshal(rc.Change.Before)
	if err != nil {
		return changes
	}
	afterBytes, err := json.Marshal(rc.Change.After)
	if err != nil {
		return changes
	}

	var prettyJSONBefore, prettyJSONAfter bytes.Buffer
	if err = json.Indent(&prettyJSONBefore, beforeBytes, "", " "); err != nil {
		return changes
	}
	if err = json.Indent(&prettyJSONAfter, afterBytes, "", " "); err != nil {
		return changes
	}

	dmp := diffmatchpatch.New()
	lineText1, lineText2, lineArray := dmp.DiffLinesToChars(prettyJSONBefore.String(), prettyJSONAfter.String())
	diffs := dmp.DiffMain(lineText1, lineText2, false)
	diffs = dmp.DiffCharsToLines(diffs, lineArray)

	for _, diff := range diffs {
		diffText := strings.TrimSpace(diff.Text)
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			lines := strings.Split(diffText, "\n")
			var changeText string
			for _, line := range lines {
				changeText = changeText + "+ " + line
			}
			changes = append(changes, changeText)
		case diffmatchpatch.DiffDelete:
			lines := strings.Split(diffText, "\n")
			var changeText string
			for _, line := range lines {
				changeText = changeText + "- " + line
			}
			changes = append(changes, changeText)
		}
	}
	return changes
}

func (m RunDetailModel) renderResourceChanges() []string {
	if m.planJSON == nil {
		return nil
	}

	if len(m.planJSON.ResourceChanges) == 0 {
		return []string{runDetailLabelStyle.Render("No resource changes")}
	}

	// Group changes by action
	var creates, updates, deletes, replacements []*tfjson.ResourceChange

	for _, rc := range m.planJSON.ResourceChanges {
		if rc.Change == nil {
			continue
		}

		actions := rc.Change.Actions

		switch {
		case actions.Create():
			creates = append(creates, rc)
		case actions.Update():
			updates = append(updates, rc)
		case actions.Delete():
			deletes = append(deletes, rc)
		case actions.Replace():
			replacements = append(replacements, rc)
		}
	}

	var sections []string

	// Build each category section
	if createSection := m.renderResourceChangeCategory("Create", creates, runDetailCreateStyle, "+"); createSection != nil {
		sections = append(sections, createSection...)
	}

	if updateSection := m.renderResourceChangeCategory("Update", updates, runDetailUpdateStyle, "~"); updateSection != nil {
		sections = append(sections, updateSection...)
	}

	if deleteSection := m.renderResourceChangeCategory("Delete", deletes, runDetailDeleteStyle, "-"); deleteSection != nil {
		sections = append(sections, deleteSection...)
	}

	if replaceSection := m.renderResourceChangeCategory("Replace", replacements, runDetailReplaceStyle, "±"); replaceSection != nil {
		sections = append(sections, replaceSection...)
	}

	return sections
}

func (m RunDetailModel) renderHeaderSection() []string {
	var sections []string

	// Title
	sections = append(sections, runDetailTitleStyle.Render("Run Details"))
	sections = append(sections, runDetailSeparatorStyle.Render(strings.Repeat("─", 50)))

	// Run ID
	sections = append(sections, lipgloss.JoinHorizontal(lipgloss.Left,
		runDetailLabelStyle.Render("Run ID:    "),
		runDetailValueStyle.Render(m.run.ID),
	))

	// Status
	statusColor := "#cdd6f4"
	switch m.run.Status {
	case tfe.RunPlanned, tfe.RunPlannedAndFinished:
		statusColor = "#df8e1d"
	case tfe.RunApplied:
		statusColor = "#40a02b"
	case tfe.RunErrored:
		statusColor = "#d20f39"
	}
	statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(statusColor))
	sections = append(sections, lipgloss.JoinHorizontal(lipgloss.Left,
		runDetailLabelStyle.Render("Status:    "),
		statusStyle.Render(string(m.run.Status)),
	))

	// Message (if present)
	if m.run.Message != "" {
		sections = append(sections, lipgloss.JoinHorizontal(lipgloss.Left,
			runDetailLabelStyle.Render("Message:   "),
			runDetailValueStyle.Render(m.run.Message),
		))
	}

	return sections
}

func (m RunDetailModel) renderPlanSection() []string {
	var sections []string

	if m.loadingPlan {
		sections = append(sections, runDetailLabelStyle.Render("Loading plan details..."))
		return sections
	}

	if m.planJSON == nil {
		sections = append(sections, runDetailLabelStyle.Render("Plan details not available"))
		return sections
	}

	// Title
	sections = append(sections, runDetailTitleStyle.Render("Resource Changes"))
	sections = append(sections, runDetailSeparatorStyle.Render(strings.Repeat("─", 50)))

	// Resource changes
	resourceChanges := m.renderResourceChanges()
	sections = append(sections, resourceChanges...)

	return sections
}

func (m RunDetailModel) renderApplySection() []string {
	var sections []string

	if m.run.Apply == nil {
		return sections
	}

	// Title
	sections = append(sections, runDetailTitleStyle.Render("Apply Status"))
	sections = append(sections, runDetailSeparatorStyle.Render(strings.Repeat("─", 50)))

	if m.loadingApply {
		sections = append(sections, runDetailLabelStyle.Render("Loading apply status..."))
	} else if m.applyLog != "" {
		sections = append(sections, lipgloss.JoinHorizontal(lipgloss.Left,
			runDetailLabelStyle.Render("Status: "),
			runDetailValueStyle.Render(m.applyLog),
		))
	}

	return sections
}

func (m RunDetailModel) renderFooter() string {
	if m.run.Status == tfe.RunPlanned || m.run.Status == tfe.RunPlannedAndFinished {
		return runDetailHelpStyle.Render("Press 'a' to Apply  •  'esc' to go back")
	}
	return runDetailHelpStyle.Render("Press 'esc' to go back")
}

func (m RunDetailModel) View() string {
	if m.run == nil {
		return docStyle.Render("No run selected")
	}

	// Build all sections
	var allSections []string

	// Header section
	headerSections := m.renderHeaderSection()
	allSections = append(allSections, headerSections...)

	// Plan section
	planSections := m.renderPlanSection()
	if len(planSections) > 0 {
		allSections = append(allSections, "") // Spacer
		allSections = append(allSections, planSections...)
	}

	// Apply section
	applySections := m.renderApplySection()
	if len(applySections) > 0 {
		allSections = append(allSections, "") // Spacer
		allSections = append(allSections, applySections...)
	}

	// Footer section
	allSections = append(allSections, "") // Spacer
	allSections = append(allSections, runDetailSeparatorStyle.Render(strings.Repeat("─", 50)))
	allSections = append(allSections, m.renderFooter())

	// Join all sections vertically
	content := lipgloss.JoinVertical(lipgloss.Left, allSections...)

	return docStyle.Render(content)
}
