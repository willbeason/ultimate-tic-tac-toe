package main

import (
	"github.com/spf13/cobra"
	"os"
	"ultimate-tic-tac-toe/pkg/ttt"
)

func main() {
	err := mainCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}

func mainCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `battle`,
		Short: `Battles ultimate tic-tac-toe opponents.`,
		RunE:  runCmd,
	}

	return cmd
}

func runCmd(cmd *cobra.Command, _ []string) error {
	cmd.SilenceUsage = true

	return nil
}

// Battle runs n battles between self and opponent, returning the proportion of wins
// self had.
func Battle(self, opponent any, n int) float64 {
	return 0.0
}

// battle runs a battle between self and opponent.
// Returns 1.0 if self won, 0.0 if opponent, and 0.5 if a tie.
func battle(self, opponent any) float64 {
	g := &ttt.Game{
		Boards:  [3][3]*ttt.Board{{{}, {}, {}}, {{}, {}, {}}, {{}, {}, {}}},
		Winners: &ttt.Board{},
	}

	nMoves := 81
	moves := make([]ttt.Move, nMoves)
	i := 0
	for a := uint8(0); a < 3; a++ {
		for b := uint8(0); b < 3; b++ {
			for x := uint8(0); x < 3; x++ {
				for y := uint8(0); y < 3; y++ {
					moves[i] = ttt.ToMove(a, b, x, y)
					i++
				}
			}
		}
	}

	return 0.0
}
