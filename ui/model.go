package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/rayhanadev/2048/game"
	"github.com/rayhanadev/2048/storage"
)

type AppState int

const (
	StateUsernameEntry AppState = iota
	StatePlaying
	StateGameOver
	StateLeaderboard
)

type AnimationState struct {
	Active      bool
	Frame       int
	TotalFrames int
	Moves       []game.TileMove
	BoardBefore [4][4]int
	BoardAfter  [4][4]int
	NewTile     *game.Position
}

type Model struct {
	state       AppState
	game        *game.Game
	player      *storage.Player
	textInput   textinput.Model
	leaderboard []storage.LeaderboardEntry
	width       int
	height      int
	db          *storage.DB
	fingerprint string
	err         error
	animation   AnimationState
}

type tickMsg time.Time

func NewModel(db *storage.DB, fingerprint string, player *storage.Player, initialState AppState) Model {
	ti := textinput.New()
	ti.Placeholder = "Enter username"
	ti.Focus()
	ti.CharLimit = 20
	ti.Width = 20

	var bestScore int
	if player != nil {
		bestScore, _ = db.GetPlayerBestScore(player.ID)
	}

	m := Model{
		state:       initialState,
		game:        game.NewGame(bestScore),
		player:      player,
		textInput:   ti,
		db:          db,
		fingerprint: fingerprint,
	}

	return m
}

func (m Model) WithSize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}

func (m Model) Init() tea.Cmd {
	if m.state == StateUsernameEntry {
		return textinput.Blink
	}
	return nil
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*40, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		if m.animation.Active {
			m.animation.Frame++
			if m.animation.Frame >= m.animation.TotalFrames {
				m.animation.Active = false
				m.animation.Frame = 0
			} else {
				return m, tickCmd()
			}
		}
		return m, nil
	}

	if m.state == StateUsernameEntry {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	}

	switch m.state {
	case StateUsernameEntry:
		return m.handleUsernameInput(msg)
	case StatePlaying:
		return m.handleGameInput(msg)
	case StateGameOver:
		return m.handleGameOverInput(msg)
	case StateLeaderboard:
		return m.handleLeaderboardInput(msg)
	}

	return m, nil
}

func (m Model) handleUsernameInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		username := m.textInput.Value()
		if len(username) < 3 {
			return m, nil
		}

		player, err := m.db.CreatePlayer(m.fingerprint, username)
		if err != nil {
			m.err = err
			return m, nil
		}

		m.player = player
		m.state = StatePlaying
		m.game = game.NewGame(0)
		return m, nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) handleGameInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.animation.Active {
		return m, nil
	}

	var dir game.Direction
	var moved bool

	switch msg.String() {
	case "up", "w", "k":
		dir = game.Up
		moved = true
	case "down", "s", "j":
		dir = game.Down
		moved = true
	case "left", "a", "h":
		dir = game.Left
		moved = true
	case "right", "d", "l":
		dir = game.Right
		moved = true
	case "r":
		m.game.Reset()
		return m, nil
	case "b":
		entries, err := m.db.GetLeaderboard(10)
		if err == nil {
			m.leaderboard = entries
		}
		m.state = StateLeaderboard
		return m, nil
	}

	if moved {
		result := m.game.Move(dir)
		if result != nil && result.Moved {
			shouldAnimate := dir == game.Down || dir == game.Right

			if shouldAnimate {
				m.animation = AnimationState{
					Active:      true,
					Frame:       0,
					TotalFrames: 3,
					Moves:       result.Moves,
					BoardBefore: result.BoardBefore,
					BoardAfter:  result.BoardState,
					NewTile:     result.NewTile,
				}
			}

			if m.game.GameOver {
				if m.player != nil {
					m.db.SaveScore(m.player.ID, m.game.Score, m.game.MaxTile())
				}
				m.state = StateGameOver
			}

			if shouldAnimate {
				return m, tickCmd()
			}
		}
	}

	return m, nil
}

func (m Model) handleGameOverInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "r":
		m.game.Reset()
		m.state = StatePlaying
		return m, nil
	case "b":
		entries, err := m.db.GetLeaderboard(10)
		if err == nil {
			m.leaderboard = entries
		}
		m.state = StateLeaderboard
		return m, nil
	}
	return m, nil
}

func (m Model) handleLeaderboardInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "escape", "b", "enter", " ":
		if m.game.GameOver {
			m.state = StateGameOver
		} else {
			m.state = StatePlaying
		}
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case StateUsernameEntry:
		return m.renderUsernameEntry()
	case StatePlaying:
		return m.renderGame()
	case StateGameOver:
		return m.renderGameOver()
	case StateLeaderboard:
		return m.renderLeaderboard()
	}
	return ""
}
