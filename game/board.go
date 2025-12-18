package game

import (
	"math/rand"
)

const BoardSize = 4

// Position represents a cell position on the board
type Position struct {
	Row, Col int
}

// Board represents the 4x4 game grid
type Board struct {
	Grid [BoardSize][BoardSize]int
}

// NewBoard creates an empty board
func NewBoard() *Board {
	return &Board{}
}

// GetEmptyCells returns all positions with value 0
func (b *Board) GetEmptyCells() []Position {
	var empty []Position
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			if b.Grid[row][col] == 0 {
				empty = append(empty, Position{Row: row, Col: col})
			}
		}
	}
	return empty
}

func (b *Board) SpawnTile() *Position {
	empty := b.GetEmptyCells()
	if len(empty) == 0 {
		return nil
	}

	pos := empty[rand.Intn(len(empty))]

	value := 2
	if rand.Float64() < 0.1 {
		value = 4
	}

	b.Grid[pos.Row][pos.Col] = value
	return &pos
}

// Clone creates a deep copy of the board
func (b *Board) Clone() *Board {
	newBoard := &Board{}
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			newBoard.Grid[row][col] = b.Grid[row][col]
		}
	}
	return newBoard
}

// Equals checks if two boards have the same state
func (b *Board) Equals(other *Board) bool {
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			if b.Grid[row][col] != other.Grid[row][col] {
				return false
			}
		}
	}
	return true
}

// MaxTile returns the highest tile value on the board
func (b *Board) MaxTile() int {
	max := 0
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			if b.Grid[row][col] > max {
				max = b.Grid[row][col]
			}
		}
	}
	return max
}

// IsFull returns true if there are no empty cells
func (b *Board) IsFull() bool {
	return len(b.GetEmptyCells()) == 0
}
