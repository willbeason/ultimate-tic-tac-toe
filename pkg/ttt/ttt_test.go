package ttt_test

import (
	"github.com/google/go-cmp/cmp"
	"math"
	"testing"
	"ultimate-tic-tac-toe/pkg/ttt"
)

func NewGame(moves ...[2]ttt.Move) *ttt.Game {
	g := &ttt.Game{
		Boards: [3][3]*ttt.Board{
			{{}, {}, {}},
			{{}, {}, {}},
			{{}, {}, {}},
		},
		Winners: &ttt.Board{},
	}

	player := ttt.Self
	for _, m := range moves {
		g.WithMove(m, player)
		if player == ttt.Self {
			player = ttt.Opponent
		} else {
			player = ttt.Self
		}
	}

	return g
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
			game: &ttt.Game{
				Boards: [3][3]*ttt.Board{
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
				},
				Winners: &ttt.Board{{0, 1, 1}, {0, 0, 0}, {0, 0, 0}},
			},
			depth:    1,
			player:   ttt.Self,
			lastMove: ttt.Move{X: 0, Y: 0},
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
		wantMove [2]ttt.Move
	}{
		{
			name: "obvious win",
			game: &ttt.Game{
				Boards: [3][3]*ttt.Board{
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
				},
				Winners: &ttt.Board{{0, 1, 1}, {0, 0, 0}, {0, 0, 0}},
			},
			depth:    1,
			player:   ttt.Self,
			lastMove: ttt.Move{X: 0, Y: 0},
			wantMove: [2]ttt.Move{{X: 0, Y: 0}, {X: 0, Y: 0}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			moves := make([][2]ttt.Move, 81)
			nMoves := tc.game.LegalMoves(tc.lastMove, moves)
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
		wantMoves [][2]ttt.Move
	}{
		{
			name: "obvious win",
			game: &ttt.Game{
				Boards: [3][3]*ttt.Board{
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
				},
				Winners: &ttt.Board{{0, 1, 1}, {0, 0, 0}, {0, 0, 0}},
			},
			depth:    1,
			player:   ttt.Self,
			lastMove: ttt.Move{X: 0, Y: 0},
			wantMoves: [][2]ttt.Move{
				{{X: 0, Y: 0}, {X: 0, Y: 0}},
				{{X: 0, Y: 0}, {X: 1, Y: 0}},
				{{X: 0, Y: 0}, {X: 1, Y: 1}},
				{{X: 0, Y: 0}, {X: 1, Y: 2}},
				{{X: 0, Y: 0}, {X: 2, Y: 0}},
				{{X: 0, Y: 0}, {X: 2, Y: 1}},
				{{X: 0, Y: 0}, {X: 2, Y: 2}},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			moves := make([][2]ttt.Move, 81)
			nMoves := tc.game.LegalMoves(tc.lastMove, moves)
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
		lastMove  ttt.Move
		wantMoves []ttt.Move
	}{
		{
			name: "obvious win",
			board: &ttt.Board{
				{0, 1, 1}, {0, 0, 0}, {0, 0, 0},
			},
			depth:    1,
			player:   ttt.Self,
			lastMove: ttt.Move{X: 0, Y: 0},
			wantMoves: []ttt.Move{
				{X: 0, Y: 0},
				{X: 1, Y: 0},
				{X: 1, Y: 1},
				{X: 1, Y: 2},
				{X: 2, Y: 0},
				{X: 2, Y: 1},
				{X: 2, Y: 2},
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
