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
	_ = &ttt.Game{
		Boards:  [3][3]*ttt.Board{{{}, {}, {}}, {{}, {}, {}}, {{}, {}, {}}},
		Winners: &ttt.Board{},
	}

	return 0.0
}
