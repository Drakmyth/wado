package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(analyzeCmd)
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze <path>",
	Short: "Print the type",
	Long:  `Display IWAD or PWAD type of specified file.`,
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

		// header := parser.Header{}
		// err = binary.Read(f, binary.LittleEndian, &header)
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Printf("%s, %d, %d", header.Identification, header.NumLumps, header.InfoTableOffset)
	},
}
