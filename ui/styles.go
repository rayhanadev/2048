package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	BoardBackground = lipgloss.Color("#303030")
	EmptyTileColor  = lipgloss.Color("#4a4a4a")

	// Text colors
	DarkText  = lipgloss.Color("#776e65")
	LightText = lipgloss.Color("#f9f6f2")

	// Tile colors - classic 2048 colors
	TileColors = map[int]lipgloss.Color{
		0:    lipgloss.Color("#4a4a4a"),
		2:    lipgloss.Color("#eee4da"),
		4:    lipgloss.Color("#ede0c8"),
		8:    lipgloss.Color("#f2b179"),
		16:   lipgloss.Color("#f59563"),
		32:   lipgloss.Color("#f67c5f"),
		64:   lipgloss.Color("#f65e3b"),
		128:  lipgloss.Color("#edcf72"),
		256:  lipgloss.Color("#edcc61"),
		512:  lipgloss.Color("#edc850"),
		1024: lipgloss.Color("#edc53f"),
		2048: lipgloss.Color("#edc22e"),
	}

	// Tile dimensions
	TileWidth  = 6
	TileHeight = 3

	// Base styles
	BoardStyle = lipgloss.NewStyle().
			Background(BoardBackground).
			Padding(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#5a5a5a"))

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#776e65")).
			Background(lipgloss.Color("#edc22e")).
			Padding(0, 2).
			MarginBottom(1)

	ScoreBoxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#4a4a4a")).
			Foreground(LightText).
			Padding(0, 1).
			Align(lipgloss.Center)

	ScoreLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#eee4da")).
			Bold(false).
			Align(lipgloss.Center)

	ScoreValueStyle = lipgloss.NewStyle().
			Foreground(LightText).
			Bold(true).
			Align(lipgloss.Center)

	InstructionsStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#776e65")).
				MarginTop(1)

	GameOverStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#f65e3b")).
			Background(lipgloss.Color("#faf8ef")).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#f65e3b"))

	WinStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#edc22e")).
			Background(lipgloss.Color("#faf8ef")).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#edc22e"))

	GameWonStyle = WinStyle

	UsernamePromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#776e65")).
				Bold(true).
				MarginBottom(1)

	UsernameInputStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#bbada0")).
				Padding(0, 1)

	LeaderboardTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#edc22e")).
				Background(lipgloss.Color("#776e65")).
				Padding(0, 2).
				MarginBottom(1)

	LeaderboardHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#776e65")).
				Background(lipgloss.Color("#eee4da")).
				Padding(0, 1)

	LeaderboardRowStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#776e65")).
				Padding(0, 1)

	LeaderboardHighlightStyle = lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("#edc22e")).
					Padding(0, 1)
)

// GetTileStyle returns the style for a specific tile value
func GetTileStyle(value int) lipgloss.Style {
	bgColor, ok := TileColors[value]
	if !ok {
		// For values > 2048, use the 2048 color
		bgColor = TileColors[2048]
	}

	// Use dark text for low values, light text for high values
	textColor := DarkText
	if value >= 8 {
		textColor = LightText
	}

	return lipgloss.NewStyle().
		Width(TileWidth).
		Height(TileHeight).
		Background(bgColor).
		Foreground(textColor).
		Bold(true).
		Align(lipgloss.Center, lipgloss.Center)
}

// GetEmptyTileStyle returns the style for an empty tile
func GetEmptyTileStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Width(TileWidth).
		Height(TileHeight).
		Background(EmptyTileColor).
		Align(lipgloss.Center, lipgloss.Center)
}
