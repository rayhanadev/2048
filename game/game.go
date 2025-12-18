package game

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type TileMove struct {
	From   Position
	To     Position
	Value  int
	Merged bool
}

type MoveResult struct {
	Moved       bool
	Moves       []TileMove
	Score       int
	BoardBefore [BoardSize][BoardSize]int
	BoardState  [BoardSize][BoardSize]int
	NewTile     *Position
}

type Game struct {
	Board     *Board
	Score     int
	BestScore int
	GameOver  bool
	Won       bool
}

func NewGame(bestScore int) *Game {
	g := &Game{
		Board:     NewBoard(),
		Score:     0,
		BestScore: bestScore,
		GameOver:  false,
		Won:       false,
	}

	g.Board.SpawnTile()
	g.Board.SpawnTile()

	return g
}

func (g *Game) Move(dir Direction) *MoveResult {
	if g.GameOver {
		return nil
	}

	oldBoard := g.Board.Clone()
	boardBefore := g.Board.Grid

	var result *MoveResult
	switch dir {
	case Up:
		result = g.moveUp()
	case Down:
		result = g.moveDown()
	case Left:
		result = g.moveLeft()
	case Right:
		result = g.moveRight()
	}

	if g.Board.Equals(oldBoard) {
		return &MoveResult{Moved: false, BoardBefore: boardBefore, BoardState: g.Board.Grid}
	}

	result.Moved = true
	result.BoardBefore = boardBefore
	g.Score += result.Score
	if g.Score > g.BestScore {
		g.BestScore = g.Score
	}

	if g.Board.MaxTile() >= 2048 && !g.Won {
		g.Won = true
	}

	newTilePos := g.Board.SpawnTile()
	if newTilePos != nil {
		result.NewTile = newTilePos
	}
	result.BoardState = g.Board.Grid

	if !g.canMove() {
		g.GameOver = true
	}

	return result
}

func (g *Game) moveLeft() *MoveResult {
	result := &MoveResult{Moves: make([]TileMove, 0)}

	for row := 0; row < BoardSize; row++ {
		line := make([]int, BoardSize)
		positions := make([]int, BoardSize)
		for col := 0; col < BoardSize; col++ {
			line[col] = g.Board.Grid[row][col]
			positions[col] = col
		}

		newLine, moves, lineScore := compressAndMergeWithTracking(line, positions)
		result.Score += lineScore

		for _, move := range moves {
			result.Moves = append(result.Moves, TileMove{
				From:   Position{Row: row, Col: move.fromIdx},
				To:     Position{Row: row, Col: move.toIdx},
				Value:  move.value,
				Merged: move.merged,
			})
		}

		for col := 0; col < BoardSize; col++ {
			g.Board.Grid[row][col] = newLine[col]
		}
	}

	return result
}

func (g *Game) moveRight() *MoveResult {
	result := &MoveResult{Moves: make([]TileMove, 0)}

	for row := 0; row < BoardSize; row++ {
		line := make([]int, BoardSize)
		positions := make([]int, BoardSize)
		for col := 0; col < BoardSize; col++ {
			line[col] = g.Board.Grid[row][BoardSize-1-col]
			positions[col] = BoardSize - 1 - col
		}

		newLine, moves, lineScore := compressAndMergeWithTracking(line, positions)
		result.Score += lineScore

		for _, move := range moves {
			result.Moves = append(result.Moves, TileMove{
				From:   Position{Row: row, Col: move.fromIdx},
				To:     Position{Row: row, Col: BoardSize - 1 - move.toIdx},
				Value:  move.value,
				Merged: move.merged,
			})
		}

		for col := 0; col < BoardSize; col++ {
			g.Board.Grid[row][BoardSize-1-col] = newLine[col]
		}
	}

	return result
}

func (g *Game) moveUp() *MoveResult {
	result := &MoveResult{Moves: make([]TileMove, 0)}

	for col := 0; col < BoardSize; col++ {
		line := make([]int, BoardSize)
		positions := make([]int, BoardSize)
		for row := 0; row < BoardSize; row++ {
			line[row] = g.Board.Grid[row][col]
			positions[row] = row
		}

		newLine, moves, lineScore := compressAndMergeWithTracking(line, positions)
		result.Score += lineScore

		for _, move := range moves {
			result.Moves = append(result.Moves, TileMove{
				From:   Position{Row: move.fromIdx, Col: col},
				To:     Position{Row: move.toIdx, Col: col},
				Value:  move.value,
				Merged: move.merged,
			})
		}

		for row := 0; row < BoardSize; row++ {
			g.Board.Grid[row][col] = newLine[row]
		}
	}

	return result
}

func (g *Game) moveDown() *MoveResult {
	result := &MoveResult{Moves: make([]TileMove, 0)}

	for col := 0; col < BoardSize; col++ {
		line := make([]int, BoardSize)
		positions := make([]int, BoardSize)
		for row := 0; row < BoardSize; row++ {
			line[row] = g.Board.Grid[BoardSize-1-row][col]
			positions[row] = BoardSize - 1 - row
		}

		newLine, moves, lineScore := compressAndMergeWithTracking(line, positions)
		result.Score += lineScore

		for _, move := range moves {
			result.Moves = append(result.Moves, TileMove{
				From:   Position{Row: move.fromIdx, Col: col},
				To:     Position{Row: BoardSize - 1 - move.toIdx, Col: col},
				Value:  move.value,
				Merged: move.merged,
			})
		}

		for row := 0; row < BoardSize; row++ {
			g.Board.Grid[BoardSize-1-row][col] = newLine[row]
		}
	}

	return result
}

type moveInfo struct {
	fromIdx int
	toIdx   int
	value   int
	merged  bool
}

func compressAndMergeWithTracking(line []int, originalPositions []int) ([]int, []moveInfo, int) {
	score := 0
	moves := make([]moveInfo, 0)

	type tileInfo struct {
		value   int
		origPos int
	}

	compressed := make([]tileInfo, 0, BoardSize)
	for i, val := range line {
		if val != 0 {
			compressed = append(compressed, tileInfo{value: val, origPos: originalPositions[i]})
		}
	}

	result := make([]int, BoardSize)
	resultIdx := 0

	for i := 0; i < len(compressed); i++ {
		if i+1 < len(compressed) && compressed[i].value == compressed[i+1].value {
			newVal := compressed[i].value * 2
			result[resultIdx] = newVal
			score += newVal

			moves = append(moves, moveInfo{
				fromIdx: compressed[i].origPos,
				toIdx:   resultIdx,
				value:   compressed[i].value,
				merged:  true,
			})
			moves = append(moves, moveInfo{
				fromIdx: compressed[i+1].origPos,
				toIdx:   resultIdx,
				value:   compressed[i+1].value,
				merged:  true,
			})

			resultIdx++
			i++
		} else {
			result[resultIdx] = compressed[i].value

			if compressed[i].origPos != resultIdx {
				moves = append(moves, moveInfo{
					fromIdx: compressed[i].origPos,
					toIdx:   resultIdx,
					value:   compressed[i].value,
					merged:  false,
				})
			}

			resultIdx++
		}
	}

	return result, moves, score
}

func (g *Game) canMove() bool {
	if !g.Board.IsFull() {
		return true
	}

	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize-1; col++ {
			if g.Board.Grid[row][col] == g.Board.Grid[row][col+1] {
				return true
			}
		}
	}

	for col := 0; col < BoardSize; col++ {
		for row := 0; row < BoardSize-1; row++ {
			if g.Board.Grid[row][col] == g.Board.Grid[row+1][col] {
				return true
			}
		}
	}

	return false
}

func (g *Game) MaxTile() int {
	return g.Board.MaxTile()
}

func (g *Game) Reset() {
	g.Board = NewBoard()
	g.Score = 0
	g.GameOver = false
	g.Won = false
	g.Board.SpawnTile()
	g.Board.SpawnTile()
}
