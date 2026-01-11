package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/leeineian/gauth/internal/model"
	"github.com/leeineian/gauth/internal/service"
)

type tickMsg time.Time

type LiveModel struct {
	accounts []model.Account
	otpSvc   *service.OTPService
	width    int
	height   int
}

func NewLiveModel(accounts []model.Account) *LiveModel {
	return &LiveModel{
		accounts: accounts,
		otpSvc:   service.NewOTPService(),
	}
}

func (m *LiveModel) Init() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *LiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	case tickMsg:
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *LiveModel) View() string {
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true).PaddingRight(1)
	rowStyle := lipgloss.NewStyle().PaddingRight(1)

	tbl := table.New().
		Border(lipgloss.HiddenBorder()).
		Headers("ISSUER", "LABEL", "TYPE", "CODE", "REMAINING").
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return headerStyle
			}
			return rowStyle
		})

	for _, acc := range m.accounts {
		res, err := m.otpSvc.Generate(&acc)
		code := "ERROR"
		remaining := "-"

		if err == nil {
			codeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
			if res.Remaining < 5 && acc.Type == model.TypeTOTP {
				codeStyle = codeStyle.Foreground(lipgloss.Color("9"))
			}
			code = codeStyle.Render(res.Code)
			if acc.Type == model.TypeTOTP {
				remaining = fmt.Sprintf("%ds", res.Remaining)
			}
		}

		tbl.Row(
			acc.Issuer,
			acc.DisplayLabel(),
			strings.ToUpper(string(acc.Type)),
			code,
			remaining,
		)
	}

	return "\n" + tbl.Render() + "\n\nPress 'q' to exit\n"
}

func RunLiveView(accounts []model.Account) error {
	p := tea.NewProgram(NewLiveModel(accounts))
	_, err := p.Run()
	return err
}
