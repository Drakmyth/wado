package cmd

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strconv"

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
		err := convert(args[0], args[1])
		if err != nil {
			panic(err)
		}
	},
}

func convert(in_filepath string, out_filepath string) error {
	// Copy to output file so we don't have to worry about messing up the format or the source file
	err := copyFile(in_filepath, out_filepath)
	if err != nil {
		return err
	}

	// Open file
	f, err := os.OpenFile(out_filepath, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	// Read WAD header
	header := wad.WadFileHeader{}
	err = binary.Read(f, binary.LittleEndian, &header)
	if err != nil {
		return err
	}

	// Position cursor at beginning of lump directory
	_, err = f.Seek(int64(header.DirectoryOffset), io.SeekStart)
	if err != nil {
		return err
	}

	mapNameRegexp := regexp.MustCompile(`^E(\d)M(\d)$`)

	// For each lump...
	for i := int32(0); i < header.LumpCount; i++ {
		// Read the directory entry for this lump
		dir := wad.WadDirectoryEntry{}
		err = binary.Read(f, binary.LittleEndian, &dir)
		if err != nil {
			return err
		}

		// If lump is a map header...
		lumpName := wad.NameToStr(dir.LumpName[:])
		parts := mapNameRegexp.FindStringSubmatch(lumpName)
		if parts != nil {
			episodeNumber, err := strconv.Atoi(parts[1])
			if err != nil {
				return err
			}

			missionNumber, err := strconv.Atoi(parts[2])
			if err != nil {
				return err
			}

			// Convert map name from ExMy to MAPxx
			mapNumber := ((episodeNumber - 1) * 9) + missionNumber
			mapName := fmt.Sprintf("MAP%02d\x00\x00\x00", mapNumber)
			copy(dir.LumpName[:], mapName)

			// Rewind cursor to beginning of lump
			_, err = f.Seek(-wad.SIZE_DIRENTRY, io.SeekCurrent)
			if err != nil {
				return err
			}

			// Overwrite old map header with updated one
			err = binary.Write(f, binary.LittleEndian, dir)
			if err != nil {
				return err
			}
		}

		// If lump needs to be processed...
		if slices.Contains(wad.LUMPS_TO_PROCESS, lumpName) {
			switch lumpName {
			case wad.LUMP_THINGS:
				err = updateThings(f, dir)
				if err != nil {
					return err
				}
			case wad.LUMP_SIDEDEFS:
				err = updateSidedefs(f, dir)
				if err != nil {
					return err
				}
			}
		}

		// If lump needs to be renamed...
		newName, rename := wad.LUMP_REPLACEMENTS[lumpName]
		if rename {
			// Update lump name in dir entry
			copy(dir.LumpName[:], wad.StrToName(newName))

			// Rewind cursor to beginning of dir entry
			_, err = f.Seek(-wad.SIZE_DIRENTRY, io.SeekCurrent)
			if err != nil {
				return err
			}

			// Overwrite old dir entry with updated one
			err = binary.Write(f, binary.LittleEndian, dir)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(srcPath string, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	return err
}

func updateThings(f *os.File, dir wad.WadDirectoryEntry) error {
	// Remember current cursor position
	currentPosition, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	// Move cursor to lump data
	_, err = f.Seek(int64(dir.DataOffset), io.SeekStart)
	if err != nil {
		return err
	}

	// Read all things from file
	tData := make([]byte, dir.DataLength)
	numThings := dir.DataLength / int32(wad.SIZE_THING)
	things := make([]wad.Thing, numThings)

	_, err = f.Read(tData)
	if err != nil {
		return err
	}
	wad.UnmarshalThings(things, tData)

	// Replace all shotguns with SSGs
	shotguns := wad.FindAllThings(things, wad.THING_SHOTGUN)
	for _, shotgun := range shotguns {
		shotgun.Type = wad.THING_SSG
	}

	// Generate 1 Megasphere, 1 Archvile, 1 Berserk, and 1 SSG
	wad.ReplaceThingsCount(&things, wad.D2_REPLACEMENT_CANDIDATES, map[int16]int16{
		wad.THING_MEGASPHERE: 1,
		wad.ENEMY_ARCHVILE:   1,
		wad.THING_BERSERK:    1,
		wad.THING_SSG:        1,
	})

	// Replace 20% of Imps with Chaingunners
	wad.ReplaceThingsWeighted(&things, []int16{wad.ENEMY_IMP}, map[int16]float64{
		wad.ENEMY_CHAINGUNNER: 0.2,
	})

	// Replace 10% of Cacodemons with Pain Elementals
	wad.ReplaceThingsWeighted(&things, []int16{wad.ENEMY_CACO}, map[int16]float64{
		wad.ENEMY_PAIN: 0.1,
	})

	// Replace 10% of Barons with Arachnotrons, 10% with Revenants, and 30% with Hell Knights
	wad.ReplaceThingsWeighted(&things, []int16{wad.ENEMY_BARON}, map[int16]float64{
		wad.ENEMY_ARACH:    0.1,
		wad.ENEMY_REVENANT: 0.1,
		wad.ENEMY_KNIGHT:   0.3,
	})

	// Replace 10% of Pistol Zombies with Chaingunners, 5% with Medikits, 10% with Stimpacks, and 20% with Health Pots
	wad.ReplaceThingsWeighted(&things, []int16{wad.ENEMY_PISTOL}, map[int16]float64{
		wad.ENEMY_CHAINGUNNER: 0.1,
		wad.THING_MEDKIT:      0.05,
		wad.THING_STIM:        0.1,
		wad.THING_HEALTH:      0.2,
	})

	// Move cursor to lump data
	_, err = f.Seek(int64(dir.DataOffset), io.SeekStart)
	if err != nil {
		return err
	}

	// Overwrite old things lump with updated one
	_, err = f.Write(wad.MarshalThings(things))
	if err != nil {
		return err
	}

	// Return cursor to original position
	_, err = f.Seek(currentPosition, io.SeekStart)
	if err != nil {
		return err
	}

	return nil
}

func updateSidedefs(f *os.File, dir wad.WadDirectoryEntry) error {
	// Remember current cursor position
	currentPosition, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	// Move cursor to lump data
	_, err = f.Seek(int64(dir.DataOffset), io.SeekStart)
	if err != nil {
		return err
	}

	// Read all sidedefs from file
	sData := make([]byte, dir.DataLength)
	numSidedefs := dir.DataLength / int32(wad.SIZE_SIDEDEF)
	sidedefs := make([]wad.Sidedef, numSidedefs)

	_, err = f.Read(sData)
	if err != nil {
		return err
	}
	wad.UnmarshalSidedefs(sidedefs, sData)

	// Update texture names in sidedefs
	for i, sidedef := range sidedefs {
		if wad.ShouldShiftTex(sidedef) {
			sidedef.XOffset += 32
		}

		sidedef.UpperTex = wad.GetNewTexName(sidedef.UpperTex)
		sidedef.MiddleTex = wad.GetNewTexName(sidedef.MiddleTex)
		sidedef.LowerTex = wad.GetNewTexName(sidedef.LowerTex)
		sidedefs[i] = sidedef
	}

	// Move cursor to lump data
	_, err = f.Seek(int64(dir.DataOffset), io.SeekStart)
	if err != nil {
		return err
	}

	// Overwrite old sidedefs lump with updated one
	_, err = f.Write(wad.MarshalSidedefs(sidedefs))
	if err != nil {
		return err
	}

	// Return cursor to original position
	_, err = f.Seek(currentPosition, io.SeekStart)
	if err != nil {
		return err
	}

	return nil
}
