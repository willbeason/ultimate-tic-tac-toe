package ttt

import (
	"fmt"
	"math"
	"os"
)

type Move uint8

const (
	XBoard Move = 0b11000000
	YBoard Move = 0b00110000
	XCell  Move = 0b00001100
	YCell  Move = 0b00000011
)

func (m Move) String() string {
	x := ((m&XBoard)>>6)*3 + ((m & XCell) >> 2)
	y := ((m&YBoard)>>4)*3 + (m & YCell)
	return fmt.Sprintf("%d %d", y, x)
}

func ToMove(a, b, x, y uint8) Move {
	return Move(a<<6 + b<<4 + x<<2 + y)
}

func (m Move) XBoard() uint8 {
	return uint8(m&XBoard) >> 6
}

func (m Move) YBoard() uint8 {
	return uint8(m&YBoard) >> 6
}

func (m Move) XCell() uint8 {
	return uint8(m&XCell) >> 2
}

func (m Move) YCell() uint8 {
	return uint8(m & YCell)
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

func (b *Board) WithMove(x, y uint8, player Player) bool {
	b.Taken[x][y] = true

	win := false

	b.Columns[x] += player
	if b.Columns[x] == 3 || b.Columns[x] == -3 {
		win = true
	}

	b.Rows[y] += player
	if b.Rows[y] == 3 || b.Rows[y] == -3 {
		win = true
	}

	if x == y {
		b.Diagonals[0] += player
		if b.Diagonals[0] == 3 || b.Diagonals[0] == -3 {
			win = true
		}
	}

	if x+y == 2 {
		b.Diagonals[1] += player
		if b.Diagonals[1] == 3 || b.Diagonals[1] == -3 {
			win = true
		}
	}

	return win
}

func (b *Board) WithoutMove(x, y uint8, player Player) {
	b.Taken[x][y] = false

	b.Columns[x] -= player
	b.Rows[y] -= player

	if x == y {
		b.Diagonals[0] -= player
	}

	if x+y == 2 {
		b.Diagonals[1] -= player
	}
}

func (b *Board) LegalMoves(out []Move) int {
	nMoves := 0
	for x, col := range b.Taken {
		for y, taken := range col {
			if !taken {
				out[nMoves] = Move((x << 2) + y)
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

func (g *Game) WithMove(a, b, x, y uint8, player Player) (bool, bool) {
	boardWinner := g.Boards[a][b].WithMove(x, y, player)
	var gameWinner bool
	if boardWinner {
		gameWinner = g.Winners.WithMove(x, y, player)
	}

	return gameWinner, boardWinner
}

func (g *Game) WithoutMove(a, b, x, y uint8, player Player, wasBoardWin bool) {
	g.Boards[a][b].WithoutMove(x, y, player)

	if wasBoardWin {
		g.Winners.WithoutMove(x, y, player)
	}
}

func (g *Game) LegalMoves(x, y uint8, out []Move) int {
	legalBoards := make([]Move, 9)

	legalBoards[0] = Move((x << 2) + y)
	nBoards := 1

	if g.Winners.Taken[x][y] {
		nBoards = g.Winners.LegalMoves(legalBoards)
	}

	moves := make([]Move, 9)
	nLegalMoves := 0
	for i, legalBoard := range legalBoards {
		a := legalBoard.XCell()
		b := legalBoard.YCell()

		if i >= nBoards {
			break
		}

		nMoves := g.Boards[a][b].LegalMoves(moves)
		for j, move := range moves {
			if j >= nMoves {
				break
			}
			out[nLegalMoves] = (legalBoard << 4) + move
			nLegalMoves++
		}
	}

	return nLegalMoves
}

func (b *Board) Score() int8 {
	return b.Columns[0] + b.Columns[1] + b.Columns[2] + b.Rows[0] + b.Rows[1] + b.Rows[2] + b.Diagonals[0] + b.Diagonals[1]
}

func Minimax(game *Game, depth int, player Player, move Move) float64 {
	if depth == 0 {
		var score int16
		score = int16(game.Winners.Score()) * 100

		for _, row := range game.Boards {
			for _, b := range row {
				score += int16(b.Score())
			}
		}

		return float64(score)
	}

	var value float64
	legalMoves := make([]Move, 81)
	a := move.XCell()
	b := move.YCell()
	if player == Self {
		// Evaluate own moves.
		value = -100
		nLegalMoves := game.LegalMoves(a, b, legalMoves)
		for i, nextMove := range legalMoves {
			if i >= nLegalMoves {
				break
			}

			a := nextMove.XBoard()
			b := nextMove.YBoard()
			x := nextMove.XBoard()
			y := nextMove.YBoard()

			isWin, winsBoard := game.WithMove(a, b, x, y, Self)
			if isWin {
				// We can win the game.
				game.WithoutMove(a, b, x, y, Self, winsBoard)
				return math.Inf(1.0)
			}

			nextMoveValue := Minimax(game, depth-1, Opponent, nextMove)
			game.WithoutMove(a, b, x, y, Self, winsBoard)

			if winsBoard {
				// We can win a board.
				nextMoveValue += 1.0
			}

			value = math.Max(value, nextMoveValue)
		}
	} else {
		value = math.Inf(1.0)
		// Evaluate opponent moves.
		nLegalMoves := game.LegalMoves(a, b, legalMoves)
		for i, nextMove := range legalMoves {
			if i >= nLegalMoves {
				break
			}

			a := nextMove.XBoard()
			b := nextMove.YBoard()
			x := nextMove.XBoard()
			y := nextMove.YBoard()

			isWin, winsBoard := game.WithMove(a, b, x, y, Opponent)
			if isWin {
				game.WithoutMove(a, b, x, y, Opponent, winsBoard)
				// Opponent can win the game.
				return math.Inf(-1.0)
			}

			nextMoveValue := Minimax(game, depth-1, Self, nextMove)
			game.WithoutMove(a, b, x, y, Opponent, winsBoard)

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

func PickMove(moves []Move, game *Game, depth int) Move {
	// Default to first valid move.
	choice := moves[0]

	value := math.Inf(-1.0)

	for i, move := range moves {
		if debug {
			fmt.Fprintf(os.Stderr, "%d/%d: %s", i, len(moves), move)
		}

		a := move.XBoard()
		b := move.YBoard()
		x := move.XCell()
		y := move.YCell()

		isWin, winsBoard := game.WithMove(a, b, x, y, Self)
		if isWin {
			choice = move
			value = math.Inf(1.0)
			game.WithoutMove(a, b, x, y, Self, winsBoard)
			fmt.Fprintln(os.Stderr, "Wins game")
			break
		}

		moveValue := Minimax(game, depth-1, Opponent, move)
		game.WithoutMove(a, b, x, y, Self, winsBoard)

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
