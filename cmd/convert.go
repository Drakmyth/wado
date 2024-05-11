package cmd

import (
	"errors"
	"fmt"

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
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		wad, err := wad.OpenWadFile(args[0])
		if err != nil {
			panic(err)
		}

		fmt.Println("Type:", wad.Type)

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

const (
	THING_SHOTGUN    = 2001
	THING_SSG        = 82
	THING_MEDKIT     = 2012
	THING_STIM       = 2011
	THING_HEALTH     = 2014
	THING_MEGASPHERE = 83
	THING_BERSERK    = 2023

	ENEMY_SHOTGUN = 9
	ENEMY_IMP     = 3001
	ENEMY_PINKY   = 3002
	ENEMY_BARON   = 3003
	ENEMY_PISTOL  = 3004
	ENEMY_CACO    = 3005
	ENEMY_SOUL    = 3006

	ENEMY_ARCHVILE    = 64
	ENEMY_CHAINGUNNER = 65
	ENEMY_REVENANT    = 66
	ENEMY_ARACH       = 68
	ENEMY_KNIGHT      = 69
	ENEMY_PAIN        = 71
)

func ParseThings() {

}
