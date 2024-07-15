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

type Player int8

const (
	None Player = iota
	Self
	Opponent
)

type Board [3][3]Player

func (b *Board) LegalMoves(out []Move) int {
	nMoves := 0
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if b[x][y] == None {
				out[nMoves] = Move{X: int8(x), Y: int8(y)}
				nMoves++
			}
		}
	}

	return nMoves
}

func (b *Board) WithMove(move Move, player Player) bool {
	b[move.X][move.Y] = player

	if b[move.X][0] == b[move.X][1] && b[move.X][1] == b[move.X][2] {
		return true
	}

	if b[0][move.Y] == b[1][move.Y] && b[1][move.Y] == b[2][move.Y] {
		return true
	}

	// On a diagonal if parity of X and Y is the same.
	return (b[1][1] == player) && ((move.X+move.Y)%2 == 0) && (b[0][0] == player && b[2][2] == player ||
		b[2][0] == player && b[0][2] == player)
}

func (b *Board) WithoutMove(move Move) {
	b[move.X][move.Y] = None
}

func (b *Board) String() string {
	return fmt.Sprintf("{{%d, %d, %d}, {%d, %d, %d}, {%d, %d, %d}}",
		b[0][0], b[0][1], b[0][2],
		b[1][0], b[1][1], b[1][2],
		b[2][0], b[2][1], b[2][2],
	)
}

type Game struct {
	Boards  [3][3]*Board
	Winners *Board
}

func (g *Game) WithMove(move [2]Move, player Player) (bool, bool) {
	boardWinner := g.Boards[move[0].X][move[0].Y].WithMove(move[1], player)

	var gameWinner bool
	if boardWinner {
		gameWinner = g.Winners.WithMove(move[0], player)
	}

	return gameWinner, boardWinner
}

func (g *Game) WithoutMove(move [2]Move, wasBoardWin bool) {
	g.Boards[move[0].X][move[0].Y].WithoutMove(move[1])

	if wasBoardWin {
		g.Winners.WithoutMove(move[0])
	}
}

// LegalMoves outputs the list of legal moves given the current game state
// and the sub-indices of the previously-made move.
func (g *Game) LegalMoves(previous Move, out [][2]Move) int {
	legalBoards := make([]Move, 9)

	legalBoards[0] = previous
	nBoards := 1

	if g.Winners[previous.X][previous.Y] != None {
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

func (g *Game) String() string {
	result := ""
	for i := 0; i < 3; i++ {
		result += "{\n"
		for j := 0; j < 3; j++ {
			result += "\t" + g.Boards[i][j].String() + ",\n"
		}
		result += "},\n"
	}

	return result
}

func Minimax(game *Game, depth int, player Player, move Move) float64 {
	if depth == 0 {
		selfBoards := 1.0
		opponentBoards := 1.0
		for _, col := range game.Winners {
			for _, winner := range col {
				switch winner {
				case Self:
					selfBoards *= 2.0
				case Opponent:
					opponentBoards *= 2.0
				}
			}
		}
		return selfBoards - opponentBoards
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
				game.WithoutMove(nextMove, winsBoard)
				return math.Inf(1.0)
			}

			nextMoveValue := Minimax(game, depth-1, Opponent, nextMove[1])
			game.WithoutMove(nextMove, winsBoard)

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
				game.WithoutMove(nextMove, winsBoard)
				// Opponent can win the game.
				return math.Inf(-1.0)
			}

			nextMoveValue := Minimax(game, depth-1, Self, nextMove[1])
			game.WithoutMove(nextMove, winsBoard)

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
			game.WithoutMove(move, winsBoard)
			break
		}

		moveValue := Minimax(game, depth-1, Opponent, move[1])
		game.WithoutMove(move, winsBoard)

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
