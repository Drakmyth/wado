package wad

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type wadFileHeader struct {
	Identifier      [4]byte
	LumpCount       int32
	DirectoryOffset int32
}

type wadDirectoryEntry struct {
	DataOffset int32
	DataLength int32
	LumpName   [8]byte
}

func strToName(str string) []byte {
	name := [8]byte{}
	paddedName := strings.ReplaceAll(fmt.Sprintf("%-8s", str), " ", "\x00")
	copy(name[:], paddedName)
	return name[:]
}

func nameToStr(name []byte) string {
	return strings.Trim(string(name), "\x00")
}

const (
	LUMP_THINGS   = "THINGS"
	LUMP_SIDEDEFS = "SIDEDEFS"
)

const (
	SIZE_DIRENTRY = 16
	SIZE_SECTOR   = 26
)

const (
	THING_SHOTGUN    int16 = 2001
	THING_SSG        int16 = 82
	THING_MEDKIT     int16 = 2012
	THING_STIM       int16 = 2011
	THING_HEALTH     int16 = 2014
	THING_MEGASPHERE int16 = 83
	THING_BERSERK    int16 = 2023
)

const (
	ENEMY_SHOTGUN int16 = 9
	ENEMY_IMP     int16 = 3001
	ENEMY_PINKY   int16 = 3002
	ENEMY_BARON   int16 = 3003
	ENEMY_PISTOL  int16 = 3004
	ENEMY_CACO    int16 = 3005
	ENEMY_SOUL    int16 = 3006

	ENEMY_ARCHVILE    int16 = 64
	ENEMY_CHAINGUNNER int16 = 65
	ENEMY_REVENANT    int16 = 66
	ENEMY_ARACH       int16 = 68
	ENEMY_KNIGHT      int16 = 69
	ENEMY_PAIN        int16 = 71
)

var TEXTURE_REPLACEMENTS = map[string]string{
	"AASTINKY": "DOORSTOP",
	"ASHWALL":  "ASHWALL2",
	"BLODGR1":  "CEMENT9",
	"BLODGR2":  "CEMENT9",
	"BLODGR3":  "CEMENT9",
	"BLODGR4":  "CEMENT9",
	"BRNBIGC":  "MIDGRATE",
	"BRNBIGL":  "MIDGRATE",
	"BRNBIGR":  "MIDGRATE",
	"BRNPOIS2": "BROWN96",
	"BROVINE":  "BROWN1",
	"BROWNWEL": "BROWNHUG",
	"CEMPOIS":  "CEMENT1",
	"COMP2":    "COMPTALL",
	"COMPOHSO": "COMPWERD",
	"COMPTILE": "COMPWERD",
	"COMPUTE1": "COMPSTA1",
	"COMPUTE2": "COMPTALL",
	"COMPUTE3": "COMPTALL",
	"DOORHI":   "TEKBRON2",
	"GRAYDANG": "GRAY5",
	"ICKDOOR1": "DOOR1",
	"ICKWALL6": "ICKWALL5",
	"LITE2":    "BROWN1",
	"LITE4":    "LITE5",
	"LITE96":   "BROWN96",
	"LITEBLU2": "LITEBLU1",
	"LITEBLU3": "LITEBLU1",
	"LITEMET":  "METAL1",
	"LITERED":  "DOORRED",
	"LITESTON": "STONE2",
	"MIDVINE1": "MIDGRATE",
	"MIDVINE2": "MIDGRATE",
	"NUKESLAD": "SLADWALL",
	"PLANET1":  "COMPSTA2",
	"REDWALL1": "REDWALL",
	"SKINBORD": "SKINMET1",
	"SKINTEK1": "SKINMET2",
	"SKINTEK2": "SKSPINE1",
	"SKULWAL3": "SKSPINE1",
	"SKULWALL": "SKSPINE1",
	"SLADRIP1": "SLADWALL",
	"SLADRIP2": "SLADWALL",
	"SLADRIP3": "SLADWALL",
	"SP_DUDE3": "SP_DUDE4",
	"SP_DUDE6": "SP_DUDE4",
	"SP_ROCK2": "SP_ROCK1",
	"STARTAN1": "STARTAN2",
	"STONGARG": "STONE3",
	"STONPOIS": "STONE",
	"TEKWALL2": "TEKWALL4",
	"TEKWALL3": "TEKWALL4",
	"TEKWALL5": "TEKWALL4",
	"WOODSKUL": "WOODGARG",
}

// These changed from 64x128 to 128x128 in Doom 2
var SHIFT_TEXTURES = []string{"BRNPOIS", "NUKEPOIS", "SW1BRN1", "SW1STON2", "SW1STONE", "SW2BRN1", "SW2STON2", "SW2STONE"}

var LUMP_REPLACEMENTS = map[string]string{
	"D_INTER":  "D_DM2INT",
	"D_INTRO":  "D_DM2TTL",
	"D_VICTOR": "D_READ_M",
	"SKY1":     "RSKY1",
	"SKY2":     "RSKY2",
	"SKY3":     "RSKY3",
	"DEMO1":    "DEMO1_D",
	"DEMO2":    "DEMO2_D",
	"DEMO3":    "DEMO3_D",
}

var LUMPS_TO_PROCESS = []string{LUMP_THINGS, LUMP_SIDEDEFS}
var D2_REPLACEMENT_CANDIDATES = []int16{ENEMY_SHOTGUN, ENEMY_IMP, ENEMY_PINKY, ENEMY_BARON, ENEMY_PISTOL, ENEMY_CACO, ENEMY_SOUL}

func findAllIndices(things []Thing, thingTypes ...int16) []int {
	indices := make([]int, 0, 10) // Arbitrarily start with 10 capacity since we don't know how many objects we'll find
	for i, thing := range things {
		if slices.Contains(thingTypes, thing.Type) {
			indices = append(indices, i)
		}
	}

	return indices
}

func Convert(in_filepath string, out_filepath string) error {
	// Copy to output file so we don't have to worry about messing up the format or the source file
	err := CopyFile(in_filepath, out_filepath)
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
	header := wadFileHeader{}
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
		dir := wadDirectoryEntry{}
		err = binary.Read(f, binary.LittleEndian, &dir)
		if err != nil {
			return err
		}

		// If lump is a map header...
		lumpName := strings.Trim(string(dir.LumpName[:]), "\x00")
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
			_, err = f.Seek(-SIZE_DIRENTRY, io.SeekCurrent)
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
		if slices.Contains(LUMPS_TO_PROCESS, lumpName) {
			switch lumpName {
			case LUMP_THINGS:
				err = updateThings(f, dir)
				if err != nil {
					return err
				}
			case LUMP_SIDEDEFS:
				err = updateSidedefs(f, dir)
				if err != nil {
					return err
				}
			}
		}

		// If lump needs to be renamed...
		newNameStr, rename := LUMP_REPLACEMENTS[lumpName]
		if rename {
			// Update lump name in dir entry
			newName := strings.ReplaceAll(fmt.Sprintf("%-8s", newNameStr), " ", "\x00")
			copy(dir.LumpName[:], newName)

			// Rewind cursor to beginning of dir entry
			_, err = f.Seek(-SIZE_DIRENTRY, io.SeekCurrent)
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

func updateThings(f *os.File, dir wadDirectoryEntry) error {
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
	tData := make([]byte, dir.DataLength)
	numThings := dir.DataLength / int32(SIZE_THING)
	things := make([]Thing, numThings)

	_, err = f.Read(tData)
	if err != nil {
		return err
	}
	unmarshalThings(things, tData)

	// Replace all shotguns with SSGs
	shotgunIndices := findAllIndices(things, THING_SHOTGUN)
	for _, i := range shotgunIndices {
		shotgun := things[i]
		shotgun.Type = THING_SSG
		things[i] = shotgun
	}

	// Generate 1 Megasphere, 1 Archvile, 1 Berserk, and 1 SSG
	itemsToGenerate := []int16{THING_MEGASPHERE, ENEMY_ARCHVILE, THING_BERSERK, THING_SSG}
	candidateIndices := findAllIndices(things, D2_REPLACEMENT_CANDIDATES...)
	for done := false; !done; done = len(itemsToGenerate) == 0 || len(candidateIndices) == 0 {
		// Pick a random index
		candidateIndexIndex := rand.Intn(len(candidateIndices))
		candidateIndex := candidateIndices[candidateIndexIndex]

		// Replace the thing
		candidate := things[candidateIndex]
		generatedThingIndex := rand.Intn(len(itemsToGenerate))
		candidate.Type = itemsToGenerate[generatedThingIndex]
		itemsToGenerate = removeIndex(itemsToGenerate, int16(generatedThingIndex))
		things[candidateIndex] = candidate

		// Remove the index of the replaced thing from the candidate list
		candidateIndices = removeIndex(candidateIndices, candidateIndexIndex)
	}

	// Get Imps, Cacodemons, Barons, and Troopers
	impIndices := findAllIndices(things, ENEMY_IMP)
	cacoIndices := findAllIndices(things, ENEMY_CACO)
	baronIndices := findAllIndices(things, ENEMY_BARON)
	pistolIndices := findAllIndices(things, ENEMY_PISTOL)

	// Replace 20% of Imps with Chaingunners
	impReplacements := []int16{}
	impReplacements = append(impReplacements, repeatedSlice(ENEMY_CHAINGUNNER, int16(math.Round(float64(len(impIndices))*0.2)))...)
	for done := len(impReplacements) == 0 || len(impIndices) == 0; !done; done = len(impReplacements) == 0 || len(impIndices) == 0 {
		// Pick a random index
		impIndexIndex := rand.Intn(len(impIndices))
		impIndex := impIndices[impIndexIndex]

		// Replace the imp
		imp := things[impIndex]
		replacementIndex := rand.Intn(len(impReplacements))
		imp.Type = impReplacements[replacementIndex]
		impReplacements = removeIndex(impReplacements, int16(replacementIndex))
		things[impIndex] = imp

		// Remove the index of the replaced imp from the candidate list
		impIndices = removeIndex(impIndices, impIndexIndex)
	}

	// Replace 10% of Cacodemons with Pain Elementals
	cacoReplacements := []int16{}
	cacoReplacements = append(cacoReplacements, repeatedSlice(ENEMY_PAIN, int16(math.Round(float64(len(cacoIndices))*0.1)))...)
	for done := len(cacoReplacements) == 0 || len(cacoIndices) == 0; !done; done = len(cacoReplacements) == 0 || len(cacoIndices) == 0 {
		// Pick a random index
		cacoIndexIndex := rand.Intn(len(cacoIndices))
		cacoIndex := cacoIndices[cacoIndexIndex]

		// Replace the caco
		caco := things[cacoIndex]
		replacementIndex := rand.Intn(len(cacoReplacements))
		caco.Type = cacoReplacements[replacementIndex]
		cacoReplacements = removeIndex(cacoReplacements, int16(replacementIndex))
		things[cacoIndex] = caco

		// Remove the index of the replaced caco from the candidate list
		cacoIndices = removeIndex(cacoIndices, cacoIndexIndex)
	}

	// Replace 10% of BaronsImps with Arachnotrons, 10% with Revenants, and 30% with Hell Knights
	baronReplacements := []int16{}
	baronReplacements = append(baronReplacements, repeatedSlice(ENEMY_ARACH, int16(math.Round(float64(len(baronIndices))*0.1)))...)
	baronReplacements = append(baronReplacements, repeatedSlice(ENEMY_REVENANT, int16(math.Round(float64(len(baronIndices))*0.1)))...)
	baronReplacements = append(baronReplacements, repeatedSlice(ENEMY_KNIGHT, int16(math.Round(float64(len(baronIndices))*0.3)))...)
	for done := len(baronReplacements) == 0 || len(baronIndices) == 0; !done; done = len(baronReplacements) == 0 || len(baronIndices) == 0 {
		// Pick a random index
		baronIndexIndex := rand.Intn(len(baronIndices))
		baronIndex := baronIndices[baronIndexIndex]

		// Replace the baron
		baron := things[baronIndex]
		replacementIndex := rand.Intn(len(baronReplacements))
		baron.Type = baronReplacements[replacementIndex]
		baronReplacements = removeIndex(baronReplacements, int16(replacementIndex))
		things[baronIndex] = baron

		// Remove the index of the replaced baron from the candidate list
		baronIndices = removeIndex(baronIndices, baronIndexIndex)
	}

	// Replace 10% of Pistol Zombies with Chaingunners, 5% with Medikits, 10% with Stimpacks, and 20% with Health Pots
	pistolReplacements := []int16{}
	pistolReplacements = append(pistolReplacements, repeatedSlice(ENEMY_CHAINGUNNER, int16(math.Round(float64(len(pistolIndices))*0.1)))...)
	pistolReplacements = append(pistolReplacements, repeatedSlice(THING_MEDKIT, int16(math.Round(float64(len(pistolIndices))*0.05)))...)
	pistolReplacements = append(pistolReplacements, repeatedSlice(THING_STIM, int16(math.Round(float64(len(pistolIndices))*0.1)))...)
	pistolReplacements = append(pistolReplacements, repeatedSlice(THING_HEALTH, int16(math.Round(float64(len(pistolIndices))*0.2)))...)
	for done := len(pistolReplacements) == 0 || len(pistolIndices) == 0; !done; done = len(pistolReplacements) == 0 || len(pistolIndices) == 0 {
		// Pick a random index
		pistolIndexIndex := rand.Intn(len(pistolIndices))
		pistolIndex := pistolIndices[pistolIndexIndex]

		// Replace the pistol
		pistol := things[pistolIndex]
		replacementIndex := rand.Intn(len(pistolReplacements))
		pistol.Type = pistolReplacements[replacementIndex]
		pistolReplacements = removeIndex(pistolReplacements, int16(replacementIndex))
		things[pistolIndex] = pistol

		// Remove the index of the replaced pistol from the candidate list
		pistolIndices = removeIndex(pistolIndices, pistolIndexIndex)
	}

	// Move cursor to lump data
	_, err = f.Seek(int64(dir.DataOffset), io.SeekStart)
	if err != nil {
		return err
	}

	// Overwrite old things lump with updated one
	_, err = f.Write(marshalThings(things))
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

func repeatedSlice[E int | int16](value, n E) []E {
	arr := make([]E, n)
	for i := E(0); i < n; i++ {
		arr[i] = value
	}
	return arr
}

func removeIndex[E int | int16](s []E, index E) []E {
	ret := make([]E, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func updateSidedefs(f *os.File, dir wadDirectoryEntry) error {
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
	numSidedefs := dir.DataLength / int32(SIZE_SIDEDEF)
	sidedefs := make([]Sidedef, numSidedefs)

	_, err = f.Read(sData)
	if err != nil {
		return err
	}
	unmarshalSidedefs(sidedefs, sData)

	// Update texture names in sidedefs
	for i, sidedef := range sidedefs {
		if shouldShiftTex(sidedef) {
			sidedef.XOffset += 32
		}

		sidedef.UpperTex = GetNewTexName(sidedef.UpperTex)
		sidedef.MiddleTex = GetNewTexName(sidedef.MiddleTex)
		sidedef.LowerTex = GetNewTexName(sidedef.LowerTex)
		sidedefs[i] = sidedef
	}

	// Move cursor to lump data
	_, err = f.Seek(int64(dir.DataOffset), io.SeekStart)
	if err != nil {
		return err
	}

	// Overwrite old sidedefs lump with updated one
	_, err = f.Write(marshalSidedefs(sidedefs))
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

func shouldShiftTex(sidedef Sidedef) bool {
	return slices.ContainsFunc(SHIFT_TEXTURES, func(tex string) bool {
		return tex == sidedef.UpperTex || tex == sidedef.LowerTex || tex == sidedef.MiddleTex
	})
}

func GetNewTexName(oldName string) string {
	newName, replaced := TEXTURE_REPLACEMENTS[oldName]
	if !replaced {
		newName = oldName
	}

	return newName
}

func CopyFile(srcPath string, destPath string) error {
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
