package utils

import (
	"github.com/charmbracelet/lipgloss"
)

var Info = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
var Error = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
var Warning = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
var Block = lipgloss.NewStyle().Width(30).Height(5).Align(lipgloss.Center, lipgloss.Center).Border(lipgloss.NormalBorder()).Foreground(lipgloss.Color("10"))
