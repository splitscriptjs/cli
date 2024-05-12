package main

import (
	"github.com/charmbracelet/lipgloss"
)

var info = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
var errMessage = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
var warning = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
var block = lipgloss.NewStyle().Width(30).Height(5).Align(lipgloss.Center, lipgloss.Center).Border(lipgloss.NormalBorder()).Foreground(lipgloss.Color("10"))
