package tui

import (
	"fmt"
	"strings"

	"sync/atomic"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ksckaan1/templ-iconify/internal/core/domain"
)

type DownloadModel struct {
	total   int
	current atomic.Int64

	logs    []string
	OnStart chan struct{}

	percent  float64
	progress progress.Model

	// styles
	logStyle     lipgloss.Style
	doneMsgStyle lipgloss.Style
}

func NewDownloadModel(total int) *DownloadModel {
	return &DownloadModel{
		total:   total,
		current: atomic.Int64{},
		logs:    make([]string, 0),
		OnStart: make(chan struct{}),

		logStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("242")),
		doneMsgStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
	}
}

func (m *DownloadModel) Init() tea.Cmd {
	m.progress = progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
	close(m.OnStart)
	return nil
}

func (m *DownloadModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch message := msg.(type) {
	case tea.KeyMsg:
		switch message.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case SaveMsg:
		if len(m.logs) < 10 {
			m.logs = append(m.logs, message.Icon.Prefix+":"+message.Icon.Name)
		} else {
			m.logs = append(m.logs[1:], message.Icon.Prefix+":"+message.Icon.Name)
		}
		m.percent = float64(m.current.Load()) / float64(m.total)
		m.current.Add(1)
	}
	return m, nil
}

func (m *DownloadModel) View() string {
	var result string
	if m.current.Load() != int64(m.total) {
		result = "Downloading icons...\n\n"

		result += m.logStyle.Render(strings.Join(m.logs, "\n"))

		result += "\n\n"
		result += m.progress.ViewAs(m.percent)
		result += fmt.Sprintf(" %d/%d\n", m.current.Load(), m.total)
	} else {
		result += m.doneMsgStyle.Render(fmt.Sprintf("%d icons downloaded\n", m.total))
	}
	return result
}

type SaveMsg struct {
	Icon *domain.Icon
}
