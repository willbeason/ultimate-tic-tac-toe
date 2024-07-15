package ttt

import (
	"fmt"
	"math"
	"os"
)

type Move struct {
	X, Y int8
}

func (m Move) String() string {
	return fmt.Sprintf("%d %d", m.Y, m.X)
}

type Player = int8

const (
	None     Player = 0
	Self            = 1
	Opponent        = -1
)

type Board struct {
	Columns   [3]int8
	Rows      [3]int8
	Diagonals [2]int8
	Taken     [3][3]bool
}

func (b *Board) WithMove(move Move, player Player) bool {
	b.Taken[move.X][move.Y] = true

	win := false

	b.Columns[move.X] += player
	if b.Columns[move.X] == 3 || b.Columns[move.X] == -3 {
		win = true
	}

	b.Rows[move.Y] += player
	if b.Rows[move.Y] == 3 || b.Rows[move.Y] == -3 {
		win = true
	}

	if move.X == move.Y {
		b.Diagonals[0] += player
		if b.Diagonals[0] == 3 || b.Diagonals[0] == -3 {
			win = true
		}
	}

	if move.X+move.Y == 2 {
		b.Diagonals[1] += player
		if b.Diagonals[1] == 3 || b.Diagonals[1] == -3 {
			win = true
		}
	}

	return win
}

func (b *Board) WithoutMove(move Move, player Player) {
	b.Taken[move.X][move.Y] = false

	b.Columns[move.X] -= player
	b.Rows[move.Y] -= player

	if move.X == move.Y {
		b.Diagonals[0] -= player
	}

	if move.X+move.Y == 2 {
		b.Diagonals[1] -= player
	}
}

func (b *Board) LegalMoves(out []Move) int {
	nMoves := 0
	for x, col := range b.Taken {
		for y, taken := range col {
			if !taken {
				out[nMoves] = Move{X: int8(x), Y: int8(y)}
				nMoves++
			}
		}
	}
	return nMoves
}

type Game struct {
	Boards  [3][3]*Board
	Winners *Board
}

const (
	Cell0 = 1
	Cell1 = Cell0 * 3
	Cell2 = Cell1 * 3
	Cell3 = Cell2 * 3
	Cell4 = Cell3 * 3
	Cell5 = Cell4 * 3
	Cell6 = Cell5 * 3
	Cell7 = Cell6 * 3
	Cell8 = Cell7 * 3
)

var Cells = [3][3]uint16{{Cell0, Cell1, Cell2}, {Cell3, Cell4, Cell5}, {Cell6, Cell7, Cell8}}

func (g *Game) WithMove(move [2]Move, player Player) (bool, bool) {
	boardWinner := g.Boards[move[0].X][move[0].Y].WithMove(move[1], player)
	var gameWinner bool
	if boardWinner {
		gameWinner = g.Winners.WithMove(move[0], player)
	}

	return gameWinner, boardWinner
}

func (g *Game) WithoutMove(move [2]Move, player Player, wasBoardWin bool) {
	g.Boards[move[0].X][move[0].Y].WithoutMove(move[1], player)

	if wasBoardWin {
		g.Winners.WithoutMove(move[0], player)
	}
}

func (g *Game) LegalMoves(previous Move, out [][2]Move) int {
	legalBoards := make([]Move, 9)

	legalBoards[0] = previous
	nBoards := 1

	if g.Winners.Taken[previous.X][previous.Y] {
		nBoards = g.Winners.LegalMoves(legalBoards)
	}

	moves := make([]Move, 9)
	nLegalMoves := 0
	for i, legalBoard := range legalBoards {
		if i >= nBoards {
			break
		}

		nMoves := g.Boards[legalBoard.X][legalBoard.Y].LegalMoves(moves)
		for j, move := range moves {
			if j >= nMoves {
				break
			}
			out[nLegalMoves] = [2]Move{legalBoard, move}
			nLegalMoves++
		}
	}

	return nLegalMoves
}

func Minimax(game *Game, depth int, player Player, move Move) float64 {
	if depth == 0 {
		score := 0.0
		for _, col := range game.Winners.Columns {
			score += float64(col)
		}
		for _, row := range game.Winners.Rows {
			score += float64(row)
		}
		for _, diagonal := range game.Winners.Diagonals {
			score += float64(diagonal)
		}
		return score
	}

	var value float64
	legalMoves := make([][2]Move, 81)
	if player == Self {
		// Evaluate own moves.
		value = -100
		nLegalMoves := game.LegalMoves(move, legalMoves)
		for i, nextMove := range legalMoves {
			if i >= nLegalMoves {
				break
			}
			isWin, winsBoard := game.WithMove(nextMove, Self)
			if isWin {
				// We can win the game.
				game.WithoutMove(nextMove, Self, winsBoard)
				return math.Inf(1.0)
			}

			nextMoveValue := Minimax(game, depth-1, Opponent, nextMove[1])
			game.WithoutMove(nextMove, Self, winsBoard)

			if winsBoard {
				// We can win a board.
				nextMoveValue += 1.0
			}

			value = math.Max(value, nextMoveValue)
		}
	} else {
		value = math.Inf(1.0)
		// Evaluate opponent moves.
		nLegalMoves := game.LegalMoves(move, legalMoves)
		for i, nextMove := range legalMoves {
			if i >= nLegalMoves {
				break
			}
			isWin, winsBoard := game.WithMove(nextMove, Opponent)
			if isWin {
				game.WithoutMove(nextMove, Opponent, winsBoard)
				// Opponent can win the game.
				return math.Inf(-1.0)
			}

			nextMoveValue := Minimax(game, depth-1, Self, nextMove[1])
			game.WithoutMove(nextMove, Opponent, winsBoard)

			if winsBoard {
				// Opponent can win a board.
				nextMoveValue -= 1.0
			}

			if nextMoveValue < value {
				value = math.Min(value, nextMoveValue)
			}
		}
	}

	return value
}

func PickMove(moves [][2]Move, game *Game, depth int) [2]Move {
	// Default to first valid move.
	choice := moves[0]

	value := math.Inf(-1.0)

	for i, move := range moves {
		if debug {
			fmt.Fprintf(os.Stderr, "%d/%d: %s", i, len(moves), move)
		}
		isWin, winsBoard := game.WithMove(move, Self)
		if isWin {
			choice = move
			value = math.Inf(1.0)
			game.WithoutMove(move, Self, winsBoard)
			fmt.Fprintln(os.Stderr, "Wins game")
			break
		}

		moveValue := Minimax(game, depth-1, Opponent, move[1])
		game.WithoutMove(move, Self, winsBoard)

		if winsBoard {
			moveValue += 1.0
		}

		if moveValue > value {
			choice = move
			value = moveValue
		}

		if debug {
			fmt.Fprintf(os.Stderr, ": %f\n", moveValue)
			if winsBoard {
				fmt.Fprintln(os.Stderr, "Wins board")
			}
		}
	}
	return choice
}
