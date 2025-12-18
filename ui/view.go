package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderUsernameEntry() string {
	title := TitleStyle.Render("Welcome to SSH 2048!")

	var content strings.Builder
	content.WriteString("\n\n")
	content.WriteString("This appears to be your first time playing.\n")
	content.WriteString("Please enter a username (3-20 characters):\n\n")
	content.WriteString(m.textInput.View())
	content.WriteString("\n\n")
	content.WriteString("Press Enter to continue")

	if m.err != nil {
		content.WriteString("\n\n")
		content.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Render(m.err.Error()))
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3d3d5c")).
		Padding(1, 2).
		Render(content.String())

	return lipgloss.JoinVertical(lipgloss.Center, title, box)
}

func (m Model) renderGame() string {
	header := m.renderHeader()
	board := m.renderBoard()
	footer := m.renderFooter()

	return lipgloss.JoinVertical(lipgloss.Center, header, board, footer)
}

func (m Model) renderGameOver() string {
	header := m.renderHeader()
	board := m.renderBoard()

	var msg string
	if m.game.Won {
		msg = GameWonStyle.Render("🎉 You Win! 🎉")
	} else {
		msg = GameOverStyle.Render("Game Over!")
	}

	instructions := InstructionsStyle.Render("Press R to restart • B for leaderboard • Q to quit")

	return lipgloss.JoinVertical(lipgloss.Center, header, board, msg, instructions)
}

func (m Model) renderLeaderboard() string {
	title := TitleStyle.Render("🏆 Top 10 Leaderboard 🏆")

	var rows []string
	headerRow := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#f9f6f2")).
		Render(fmt.Sprintf("%-4s %-15s %-8s %-6s", "Rank", "Player", "Score", "Tile"))
	rows = append(rows, headerRow)
	rows = append(rows, strings.Repeat("─", 40))

	for i, entry := range m.leaderboard {
		row := fmt.Sprintf("%-4d %-15s %-8d %-6d",
			i+1,
			truncateString(entry.Username, 15),
			entry.Score,
			entry.MaxTile)
		rows = append(rows, row)
	}

	if len(m.leaderboard) == 0 {
		rows = append(rows, "No scores yet!")
	}

	content := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f9f6f2")).
		Render(strings.Join(rows, "\n"))

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3d3d5c")).
		Padding(1, 2).
		Render(content)

	footer := InstructionsStyle.Render("Press Enter or B to return")

	return lipgloss.JoinVertical(lipgloss.Center, title, box, footer)
}

func (m Model) renderHeader() string {
	var playerName string
	if m.player != nil {
		playerName = m.player.Username
	} else {
		playerName = "Guest"
	}

	playerInfo := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f9f6f2")).
		Render(fmt.Sprintf("Player: %s", playerName))

	scoreLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#eee4da")).
		Bold(true).
		Render("SCORE")

	scoreValue := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff")).
		Bold(true).
		Render(fmt.Sprintf("%d", m.game.Score))

	scoreBox := ScoreBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Center, scoreLabel, scoreValue))

	bestLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#eee4da")).
		Bold(true).
		Render("BEST")

	bestValue := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff")).
		Bold(true).
		Render(fmt.Sprintf("%d", m.game.BestScore))

	bestBox := ScoreBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Center, bestLabel, bestValue))

	scores := lipgloss.JoinHorizontal(lipgloss.Center, scoreBox, "  ", bestBox)

	return lipgloss.JoinVertical(lipgloss.Center, playerInfo, "", scores, "")
}

func (m Model) renderBoard() string {
	var grid [4][4]int

	if m.animation.Active {
		grid = m.getAnimatedGrid()
	} else {
		grid = m.game.Board.Grid
	}

	var rows []string

	for row := 0; row < 4; row++ {
		var tiles []string
		for col := 0; col < 4; col++ {
			value := grid[row][col]
			tile := m.renderTile(value)
			tiles = append(tiles, tile)
		}
		rowStr := lipgloss.JoinHorizontal(lipgloss.Center, tiles...)
		rows = append(rows, rowStr)
	}

	board := lipgloss.JoinVertical(lipgloss.Center, rows...)
	return BoardStyle.Render(board)
}

func (m Model) getAnimatedGrid() [4][4]int {
	progress := float64(m.animation.Frame) / float64(m.animation.TotalFrames)

	if progress >= 1.0 {
		return m.animation.BoardAfter
	}

	var result [4][4]int

	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			result[row][col] = 0
		}
	}

	occupied := make(map[[2]int]int)

	for _, move := range m.animation.Moves {
		fromRow := float64(move.From.Row)
		fromCol := float64(move.From.Col)
		toRow := float64(move.To.Row)
		toCol := float64(move.To.Col)

		currentRow := int(fromRow + (toRow-fromRow)*progress + 0.5)
		currentCol := int(fromCol + (toCol-fromCol)*progress + 0.5)

		if currentRow < 0 {
			currentRow = 0
		} else if currentRow > 3 {
			currentRow = 3
		}
		if currentCol < 0 {
			currentCol = 0
		} else if currentCol > 3 {
			currentCol = 3
		}

		pos := [2]int{currentRow, currentCol}
		if existing, ok := occupied[pos]; ok && existing == move.Value {
			result[currentRow][currentCol] = move.Value * 2
		} else {
			result[currentRow][currentCol] = move.Value
		}
		occupied[pos] = move.Value
	}

	return result
}

func (m Model) renderTile(value int) string {
	var style lipgloss.Style
	var content string

	if value == 0 {
		style = GetEmptyTileStyle()
		content = ""
	} else {
		style = GetTileStyle(value)
		content = fmt.Sprintf("%d", value)
	}

	lines := make([]string, TileHeight)
	midLine := TileHeight / 2
	for i := 0; i < TileHeight; i++ {
		if i == midLine && value != 0 {
			padding := (TileWidth - len(content)) / 2
			if padding < 0 {
				padding = 0
			}
			lines[i] = strings.Repeat(" ", padding) + content + strings.Repeat(" ", TileWidth-padding-len(content))
		} else {
			lines[i] = strings.Repeat(" ", TileWidth)
		}
	}

	return style.Render(strings.Join(lines, "\n"))
}

func (m Model) renderFooter() string {
	instructions := "↑/↓/←/→: Move • R: Restart • B: Leaderboard • Q: Quit"
	return InstructionsStyle.Render(instructions)
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
