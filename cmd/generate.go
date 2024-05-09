package cmd

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"github.com/Drakmyth/wado/models"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate <path>",
	Short: "Make a random wad",
	Long:  `Picks levels from WADs in path to generate a custom episode`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Open(args[0])
		if err != nil {
			panic(err)
		}
		defer f.Close()

		header := models.WadHeader{}
		err = binary.Read(f, binary.LittleEndian, &header)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s, %d, %d", header.Identification, header.NumLumps, header.InfoTableOffset)
	},
}
