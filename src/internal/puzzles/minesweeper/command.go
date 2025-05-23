package minesweeper

import (
	"context"
	"fmt"
	"github.com/Liuuner/go-puzzles/src/internal/puzzles"
	"github.com/urfave/cli/v3"
	"slices"
)

func Command(run puzzles.RunFunc) *cli.Command {
	return &cli.Command{
		Name:     "minesweeper",
		Usage:    "Play minesweeper",
		Aliases:  []string{"ms"},
		Category: "Games",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return run(NewWithCmd(cmd))
		},
		Flags: []cli.Flag{
			cli.HelpFlag,
			&cli.IntFlag{
				Name:        "width",
				Aliases:     []string{"w"},
				Usage:       "Width of the minesweeper field",
				Required:    false,
				HideDefault: true,
			},
			&cli.IntFlag{
				Name:        "height",
				Aliases:     []string{"h"},
				Usage:       "Height of the minesweeper field",
				Required:    false,
				HideDefault: true,
			},
			&cli.IntFlag{
				Name:        "mines",
				Aliases:     []string{"m"},
				Usage:       "Number of mines in the field",
				Required:    false,
				HideDefault: true,
			},
			&cli.StringFlag{
				Name:        "difficulty",
				Usage:       "Difficulty level of the minesweeper game (beginner, intermediate, expert)",
				Aliases:     []string{"d"},
				Required:    false,
				HideDefault: true,
				Validator: func(d string) error {
					difficulties := []string{"beginner", "intermediate", "expert"}
					if !slices.Contains(difficulties, d) {
						return fmt.Errorf("invalid difficulty level: %s, valid levels are: %v", d, difficulties)
					}
					return nil
				},
			},
		},
	}
}
