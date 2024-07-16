package main

import (
	"fmt"
	"math"
)

func main() {
	g := &Game{}
	selfMoves := make([]Move, 100)
	var nSelfMoves int

	for {
		// opponentRow: The coordinates of your opponent's last move
		var opponentRow, opponentCol int8
		_, _ = fmt.Scan(&opponentRow, &opponentCol)

		if opponentRow != -1 {
			g.WithMove(NewMove(opponentCol, opponentRow), Opponent)
		}

		// validActionCount: the number of possible actions for your next move
		_, _ = fmt.Scan(&nSelfMoves)

		for i := 0; i < nSelfMoves; i++ {
			// row: The coordinates of a possible next move
			var row, col int8
			_, _ = fmt.Scan(&row, &col)
			selfMoves[i] = NewMove(col, row)
		}

		choice := selfMoves[0]
		score := math.Inf(-1.0)
		for i := 0; i < nSelfMoves; i++ {
			move := selfMoves[i]
			dSelfScore := g.WithMove(move, Self)
			g.WithoutMove(move, Self, dSelfScore)
			dOpponentScore := g.WithMove(move, Opponent)
			g.WithoutMove(move, Self, dOpponentScore)

			moveScore := 10 * float64(dSelfScore+dOpponentScore)

			// Check neighbors
			x, y := move.XY()
			for dx := int8(-1); dx <= int8(1); dx++ {
				if x+dx < 0 || x+dx >= 10 {
					// Penalize left and right border.
					moveScore -= 3.0
					continue
				}

				for dy := int8(-1); dy <= 1; dy++ {
					if y+dy < 0 || y+dy >= 10 {
						// Penalize top and bottom border.
						moveScore -= 1.0
						continue
					}
					if dx == 0 && dy == 0 {
						continue
					}

					// Best if next to self, otherwise ok if next to many opponent.
					switch g.Cells[x+dx][y+dy] {
					case Self:
						moveScore += 2.0
					case Opponent:
						moveScore += 1.0
					}
				}
			}

			if moveScore > score {
				score = moveScore
				choice = move
			}

		}

		g.WithMove(choice, Self)

		// fmt.Fprintln(os.Stderr, "Debug messages...")
		fmt.Println(choice) // <row> <column>
	}
}

type Player uint8

const (
	None Player = iota
	Self
	Opponent
)

type Move uint8

func (m Move) XY() (int8, int8) {
	return int8(m >> 4), int8(m & 0b1111)
}

func (m Move) String() string {
	x, y := m.XY()
	return fmt.Sprintf("%d %d", y, x)
}

func NewMove(x, y int8) Move {
	return Move((uint8(x) << 4) + uint8(y))
}

type Game struct {
	Cells [10][10]Player

	// Scores is the set of scores for the board.
	// First score is Self, second is opponent.
	Scores [2]int
}

func (g *Game) ValidMoves(out []Move) int {
	nMoves := 0
	for x := int8(0); x < 10; x++ {
		for y := int8(0); y < 10; y++ {
			if g.Cells[x][y] == None {
				out[nMoves] = NewMove(x, y)
				nMoves++
			}
		}
	}
	return nMoves
}

func (g *Game) WithMove(m Move, player Player) int {
	dScore := 0
	x, y := m.XY()
	g.Cells[x][y] = player

	// Left to Right
	for left := x - 2; left <= x; left++ {
		if left < 0 {
			continue
		} else if left >= 8 {
			break
		}

		if g.Cells[left][y] == player && g.Cells[left+1][y] == player && g.Cells[left+2][y] == player {
			dScore++
		}
	}

	// Up to Down
	for up := y - 2; up <= x; up++ {
		if up < 0 {
			continue
		} else if up >= 8 {
			break
		}

		if g.Cells[x][up] == player && g.Cells[x][up+1] == player && g.Cells[x][up+2] == player {
			dScore++
		}
	}

	// UpperLeft Diagonal
	for ul := int8(-2); ul <= 0; ul++ {
		if x+ul < 0 || y+ul < 0 {
			continue
		} else if x+ul >= 8 || y+ul >= 8 {
			break
		}

		if g.Cells[x+ul][y+ul] == player && g.Cells[x+ul+1][y+ul+1] == player && g.Cells[x+ul+2][y+ul+2] == player {
			dScore++
		}
	}

	// UpperRight Diagonal
	for ur := int8(-2); ur <= 0; ur++ {
		if x+ur < 0 || y-ur-2 < 0 {
			continue
		} else if x+ur >= 8 || y-ur >= 8 {
			break
		}

		//fmt.Fprintln(os.Stderr, x, y, ur)

		if g.Cells[x+ur][y-ur] == player && g.Cells[x+ur+1][y-ur-1] == player && g.Cells[x+ur+2][y-ur-2] == player {
			dScore++
		}
	}

	g.Scores[player-1] += dScore
	return dScore
}

func (g *Game) WithoutMove(m Move, player Player, dScore int) {
	x, y := m.XY()
	g.Cells[x][y] = None
	g.Scores[player-1] -= dScore
}
