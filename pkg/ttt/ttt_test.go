package ttt_test

import (
	"github.com/google/go-cmp/cmp"
	"math"
	"testing"
	"ultimate-tic-tac-toe/pkg/ttt"
)

type Board [3][3]ttt.Player

type Game struct {
	Boards  [3][3]*Board
	Winners *Board
}

var startingGame = &Game{
	Boards: [3][3]*Board{
		{{{0, 0, 0}, {0, 1, 0}, {0, 0, 1}},
			{{0, 0, 1}, {0, 0, 1}, {0, 0, 1}},
			{{2, 0, 0}, {0, 0, 0}, {0, 0, 0}},
		},
		{{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			{{2, 0, 0}, {0, 0, 0}, {0, 1, 0}},
			{{0, 0, 0}, {0, 0, 0}, {2, 0, 0}},
		},
		{{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			{{0, 2, 0}, {0, 0, 0}, {0, 0, 0}},
			{{0, 2, 0}, {0, 2, 0}, {0, 0, 0}},
		},
	},
	Winners: &Board{{0, 1, 0}, {0, 0, 0}, {0, 0, 0}},
}

func NewBoard(start *Board) *ttt.Board {
	b := &ttt.Board{}

	for xCell := uint8(0); xCell < 3; xCell++ {
		for yCell := uint8(0); yCell < 3; yCell++ {
			owner := start[xCell][yCell]
			if owner != ttt.None {
				b.WithMove(xCell, yCell, owner)
			}
		}
	}

	return b
}

func NewGame(start [3][3]*Board) *ttt.Game {
	g := &ttt.Game{
		Boards:  [3][3]*ttt.Board{{{}, {}, {}}, {{}, {}, {}}, {{}, {}, {}}},
		Winners: &ttt.Board{},
	}

	for xBoard := uint8(0); xBoard < 3; xBoard++ {
		for yBoard := uint8(0); yBoard < 3; yBoard++ {
			for xCell := uint8(0); xCell < 3; xCell++ {
				for yCell := uint8(0); yCell < 3; yCell++ {
					owner := start[xBoard][yBoard][xCell][yCell]
					if owner != ttt.None {
						g.WithMove(xBoard, yBoard, xCell, yCell, owner)
					}
				}
			}
		}
	}

	return g
}

func BenchmarkPickMove(b *testing.B) {
	startingGame2 := NewGame(startingGame.Boards)
	moves := make([]ttt.Move, 81)
	nMoves := startingGame2.LegalMoves(2, 0, moves)
	moves = moves[:nMoves]

	for i := 0; i < b.N; i++ {
		_ = ttt.PickMove(moves[:nMoves], startingGame2, 6)
	}
}

func BenchmarkLegalMoves(b *testing.B) {
	startingGame2 := NewGame(startingGame.Boards)
	moves := make([]ttt.Move, 81)

	for i := 0; i < b.N; i++ {
		_ = startingGame2.LegalMoves(2, 0, moves)
	}
}

func BenchmarkBoard_WithMove(b *testing.B) {
	board := &ttt.Board{}

	for i := 0; i < b.N; i++ {
		_ = board.WithMove(0, 0, ttt.Self)
		board.WithoutMove(0, 0, ttt.Self)
	}
}

func TestMinimax(t *testing.T) {
	tt := []struct {
		name     string
		game     *ttt.Game
		depth    int
		player   ttt.Player
		lastMove ttt.Move
		wantEval float64
	}{
		{
			name: "obvious win",
			game: NewGame([3][3]*Board{
				{
					{{0, 1, 1}, {0, 0, 0}, {0, 0, 0}},
					{{1, 1, 1}, {0, 0, 0}, {0, 0, 0}},
					{{1, 1, 1}, {0, 0, 0}, {0, 0, 0}},
				},
				{
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
				},
				{
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
				},
			}),
			depth:    1,
			player:   ttt.Self,
			lastMove: ttt.ToMove(0, 0, 0, 0),
			wantEval: math.Inf(1.0),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := ttt.Minimax(tc.game, tc.depth, tc.player, tc.lastMove)
			if got != tc.wantEval {
				t.Errorf("got %v, want %v", got, tc.wantEval)
			}
		})
	}
}

func TestPickMove(t *testing.T) {
	tt := []struct {
		name     string
		game     *ttt.Game
		depth    int
		player   ttt.Player
		lastMove ttt.Move
		wantMove ttt.Move
	}{
		{
			name: "obvious win game",
			game: NewGame([3][3]*Board{
				{
					{{0, 1, 1}, {0, 0, 0}, {0, 0, 0}},
					{{1, 1, 1}, {0, 0, 0}, {0, 0, 0}},
					{{1, 1, 1}, {0, 0, 0}, {0, 0, 0}},
				},
				{
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
				},
				{
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
				},
			}),
			depth:    1,
			player:   ttt.Self,
			lastMove: ttt.ToMove(0, 0, 0, 0),
			wantMove: ttt.ToMove(0, 0, 0, 0),
		},
		{
			name: "obvious win board",
			game: NewGame([3][3]*Board{
				{
					{{0, 1, 1}, {0, 0, 0}, {0, 0, 0}},
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
				},
				{
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
				},
				{
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
				},
			}),
			depth:    1,
			player:   ttt.Self,
			lastMove: ttt.ToMove(0, 0, 0, 0),
			wantMove: ttt.ToMove(0, 0, 0, 0),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			moves := make([]ttt.Move, 81)
			a := tc.lastMove.XCell()
			b := tc.lastMove.YCell()
			nMoves := tc.game.LegalMoves(a, b, moves)
			moves = moves[:nMoves]
			got := ttt.PickMove(moves, tc.game, tc.depth)
			if got != tc.wantMove {
				t.Errorf("got %v, want %v", got, tc.wantMove)
			}
		})
	}
}

func TestGame_LegalMoves(t *testing.T) {
	tt := []struct {
		name      string
		game      *ttt.Game
		depth     int
		player    ttt.Player
		lastMove  ttt.Move
		wantMoves []ttt.Move
	}{
		{
			name: "obvious win",
			game: NewGame([3][3]*Board{
				{
					{{0, 1, 1}, {0, 0, 0}, {0, 0, 0}},
					{{1, 1, 1}, {0, 0, 0}, {0, 0, 0}},
					{{1, 1, 1}, {0, 0, 0}, {0, 0, 0}},
				},
				{
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
				},
				{
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
					{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
				},
			}),
			depth:    1,
			player:   ttt.Self,
			lastMove: ttt.ToMove(0, 0, 0, 0),
			wantMoves: []ttt.Move{
				ttt.ToMove(0, 0, 0, 0),
				ttt.ToMove(0, 0, 1, 0),
				ttt.ToMove(0, 0, 1, 1),
				ttt.ToMove(0, 0, 1, 2),
				ttt.ToMove(0, 0, 2, 0),
				ttt.ToMove(0, 0, 2, 1),
				ttt.ToMove(0, 0, 2, 2),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			moves := make([]ttt.Move, 81)
			a := tc.lastMove.XCell()
			b := tc.lastMove.YCell()
			nMoves := tc.game.LegalMoves(a, b, moves)
			moves = moves[:nMoves]

			if len(moves) != len(tc.wantMoves) {
				t.Errorf("got %v, want %v", len(moves), len(tc.wantMoves))
			}

			if diff := cmp.Diff(moves, tc.wantMoves); diff != "" {
				t.Errorf("got %v, want %v", moves, tc.wantMoves)
			}
		})
	}
}

func TestBoard_LegalMoves(t *testing.T) {
	tt := []struct {
		name      string
		board     *ttt.Board
		depth     int
		player    ttt.Player
		wantMoves []ttt.Move
	}{
		{
			name: "obvious win",
			board: NewBoard(&Board{
				{0, 1, 1}, {0, 0, 0}, {0, 0, 0},
			}),
			depth:  1,
			player: ttt.Self,
			wantMoves: []ttt.Move{
				ttt.ToMove(0, 0, 0, 0),
				ttt.ToMove(0, 0, 1, 0),
				ttt.ToMove(0, 0, 1, 1),
				ttt.ToMove(0, 0, 1, 2),
				ttt.ToMove(0, 0, 2, 0),
				ttt.ToMove(0, 0, 2, 1),
				ttt.ToMove(0, 0, 2, 2),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			moves := make([]ttt.Move, 9)
			nMoves := tc.board.LegalMoves(moves)
			if nMoves != len(tc.wantMoves) {
				t.Errorf("got %v, want %v", nMoves, len(tc.wantMoves))
			}
			moves = moves[:nMoves]

			if diff := cmp.Diff(moves, tc.wantMoves); diff != "" {
				t.Errorf("got %v, want %v", moves, tc.wantMoves)
			}
		})
	}
}

func TestBoard_WithMove(t *testing.T) {
	tt := []struct {
		name         string
		board        *ttt.Board
		winsSelf     [9]bool
		winsOpponent [9]bool
	}{
		{
			name:  "no wins on empty board",
			board: NewBoard(&Board{}),
			winsSelf: [9]bool{
				false, false, false,
				false, false, false,
				false, false, false,
			},
			winsOpponent: [9]bool{
				false, false, false,
				false, false, false,
				false, false, false,
			},
		},
		{
			name: "self win row",
			board: NewBoard(&Board{
				{ttt.Self, ttt.Self, ttt.None},
				{ttt.None, ttt.None, ttt.None},
				{ttt.None, ttt.None, ttt.None},
			}),
			winsSelf: [9]bool{
				false, false, true,
				false, false, false,
				false, false, false,
			},
			winsOpponent: [9]bool{
				false, false, false,
				false, false, false,
				false, false, false,
			},
		},
		{
			name: "self win column",
			board: NewBoard(&Board{
				{ttt.None, ttt.Self, ttt.Self},
				{ttt.None, ttt.None, ttt.None},
				{ttt.None, ttt.None, ttt.None},
			}),
			winsSelf: [9]bool{
				true, false, false,
				false, false, false,
				false, false, false,
			},
			winsOpponent: [9]bool{
				false, false, false,
				false, false, false,
				false, false, false,
			},
		},
		{
			name: "self win row",
			board: NewBoard(&Board{
				{ttt.None, ttt.None, ttt.Self},
				{ttt.None, ttt.None, ttt.Self},
				{ttt.None, ttt.None, ttt.None},
			}),
			winsSelf: [9]bool{
				false, false, false,
				false, false, false,
				false, false, true,
			},
			winsOpponent: [9]bool{
				false, false, false,
				false, false, false,
				false, false, false,
			},
		},
		{
			name: "self win row with opponent",
			board: NewBoard(&Board{
				{ttt.Self, ttt.Opponent, ttt.Self},
				{ttt.None, ttt.None, ttt.Self},
				{ttt.None, ttt.None, ttt.None},
			}),
			winsSelf: [9]bool{
				false, false, false,
				false, false, false,
				false, false, true,
			},
			winsOpponent: [9]bool{
				false, false, false,
				false, false, false,
				false, false, false,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var gotSelf [9]bool
			var gotOpponent [9]bool

			for x := uint8(0); x < 3; x++ {
				for y := uint8(0); y < 3; y++ {
					if tc.board.Taken[x][y] {
						continue
					}

					gotSelf[3*x+y] = tc.board.WithMove(x, y, ttt.Self)
					tc.board.WithoutMove(x, y, ttt.Self)
					gotOpponent[3*x+y] = tc.board.WithMove(x, y, ttt.Opponent)
					tc.board.WithoutMove(x, y, ttt.Opponent)
				}
			}

			if diff := cmp.Diff(gotSelf, tc.winsSelf); diff != "" {
				t.Errorf(diff)
			}
			if diff := cmp.Diff(gotOpponent, tc.winsOpponent); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
