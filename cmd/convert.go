package cmd

import (
	"errors"

	"github.com/Drakmyth/wado/wad"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(convertCmd)
}

var convertCmd = &cobra.Command{
	Use:   "convert <path>",
	Short: "Convert the WAD",
	Long:  `Converts Doom WADs to Doom 2 WADs`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires at least two args")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := wad.Convert(args[0], args[1])
		if err != nil {
			panic(err)
		}

		// fmt.Println("Break")

		// f, err := os.Open(args[0])
		// defer f.Close()

		// header := parser.Header{}
		// err = binary.Read(f, binary.LittleEndian, &header)
		// if err != nil {
		// 	panic(err)
		// }

		// _, err = f.Seek(int64(header.InfoTableOffset), io.SeekStart)
		// if err != nil {
		// 	panic(err)
		// }

		// directory := []parser.DirectoryEntry{}

		// for i := 0; i < int(header.NumLumps); i++ {
		// 	entry := parser.DirectoryEntry{}
		// 	err = binary.Read(f, binary.LittleEndian, &entry)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	directory = append(directory, entry)
		// }

		// fmt.Println("Header:", string(header.Identification[:]), header.NumLumps, header.InfoTableOffset)
		// for i := 0; i < len(directory); i++ {
		// 	entry := directory[i]
		// 	fmt.Println("Entry:", string(entry.Name[:]), entry.Size, entry.Position)
		// }
	},
}

func ParseThings() {

}
