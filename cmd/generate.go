package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"math/rand/v2"
	"path/filepath"
	"strings"

	"github.com/Drakmyth/wado/wad"
	"github.com/spf13/cobra"
)

var generateSeed uint64

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().Uint64VarP(&generateSeed, "seed", "s", 0,
		`Specify a seed value to influence randomization.
The same seed will produce the same results every
time.`)
}

var generateCmd = &cobra.Command{
	Use:   "generate [flags] <input-wad-folder> <output-wad-file>",
	Short: "Generate a new WAD with random levels",
	Long: `Generates a new WAD by randomly selecting levels
from WADs in the input folder. Will not perform
any conversion on the levels, so ensure the folder
only contains wads targetting the same game.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires input folder path and output file path")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("seed") {
			generateSeed = rand.Uint64()
		}

		err := generate(args[0], args[1])
		if err != nil {
			panic(err)
		}
	},
}

func generate(in_folderpath string, out_filepath string) error {
	wadPaths := []string{}

	// Find the wad files in the provided directory
	err := filepath.WalkDir(in_folderpath, func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(d.Name(), ".wad") {
			wadPaths = append(wadPaths, path)
		}
		return err
	})
	if err != nil {
		return err
	}

	// Read all levels from inputs wads and bucket by existance of secret exits
	levelsWithSecretExits := make([]wad.Level, 0, 9)
	levels := make([]wad.Level, 0, 9)
	for _, path := range wadPaths {
		// Open file
		wf, err := wad.OpenFile(path)
		if err != nil {
			return err
		}

		for _, level := range wf.Levels {
			if level.HasSecretExit() {
				levelsWithSecretExits = append(levelsWithSecretExits, level)
			} else {
				levels = append(levels, level)
			}
		}
	}

	// Create output wad
	wf, err := wad.CreateFile(out_filepath)
	if err != nil {
		return err
	}

	fmt.Printf("Seed: %d", generateSeed)
	rng := rand.New(rand.NewPCG(generateSeed, generateSeed))

	// Ensure exactly one level prior to level 8 has a secret exit
	secretExitLevelSlot := rng.IntN(7) + 1
	for i := 1; i < 10; i++ {
		var level wad.Level

		// Pull a random level
		if i == secretExitLevelSlot {
			level = levelsWithSecretExits[rng.IntN(len(levelsWithSecretExits))]
		} else {
			levelIndex := rng.IntN(len(levels))
			level = levels[levelIndex]
			levels = append(levels[:levelIndex], levels[levelIndex+1:]...)
		}

		// Identify the level slot
		level.Slot = fmt.Sprintf("MAP%02d", i)

		// Set exit
		next := fmt.Sprintf("MAP%02d", i+1)
		if i == 9 {
			next = fmt.Sprintf("MAP%02d", secretExitLevelSlot+1)
		}

		// Set secret exit
		nextSecret := next
		if i == secretExitLevelSlot {
			nextSecret = "MAP09"
		}

		level.LevelInfo.Next = next
		level.LevelInfo.NextSecret = nextSecret
		level.LevelInfo.EndGame = i == 8

		wf.Levels = append(wf.Levels, level)
	}

	return wf.Save()
}
