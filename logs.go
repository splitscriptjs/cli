package main

import (
	"github.com/charmbracelet/lipgloss"
)

var infoTag = lipgloss.NewStyle().Width(8).AlignHorizontal(lipgloss.Center).Background(lipgloss.Color("12"))
var errTag = lipgloss.NewStyle().Width(8).AlignHorizontal(lipgloss.Center).Background(lipgloss.Color("15"))
var info = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
var errMessage = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
var warning = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
var block = lipgloss.NewStyle().Width(30).Height(5).Align(lipgloss.Center, lipgloss.Center).Border(lipgloss.NormalBorder()).Foreground(lipgloss.Color("10"))
var checkMark = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
