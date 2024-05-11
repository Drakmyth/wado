package wad

import (
	"encoding/binary"
	"io"
	"os"
	"regexp"
	"strings"
)

type Wad struct {
	Type string // IWAD/PWAD
	Maps []LevelMap
	// Flats []Flat
	// Sprites []Sprite
	// Patches []Patch
	// Palettes []Palette
	// ColorMaps []ColorMap
	// Demos []Demo
	UnknownLumps []Lump
}

type Lump struct {
	Name string
	Data []byte
}

func (l *Lump) IsMapHeader() bool {
	matched, err := regexp.MatchString(`^E\dM\d$`, l.Name)
	if err != nil {
		return false
	}

	return matched
}

type LevelMap struct {
	Header       *Lump
	MapDataLumps []Lump
}

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

func OpenWadFile(filepath string) (*Wad, error) {
	// Open File
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read WAD Header
	header := wadFileHeader{}
	err = binary.Read(f, binary.LittleEndian, &header)
	if err != nil {
		return nil, err
	}

	// Create WAD
	wad := Wad{
		Type:         strings.Trim(string(header.Identifier[:]), "\x00"),
		UnknownLumps: []Lump{},
	}

	// Position cursor at beginning of lump directory
	_, err = f.Seek(int64(header.DirectoryOffset), io.SeekStart)
	if err != nil {
		return nil, err
	}

	var currentMap *LevelMap = nil

	// For each lump...
	for i := int32(0); i < header.LumpCount; i++ {
		// Read the directory entry for this lump
		dir := wadDirectoryEntry{}
		err = binary.Read(f, binary.LittleEndian, &dir)
		if err != nil {
			return nil, err
		}

		// If lump contains data...
		var lumpData []byte
		if dir.DataLength > 0 {
			// Store current cursor position
			currentOffset, err := f.Seek(0, io.SeekCurrent)
			if err != nil {
				return nil, err
			}

			// Position cursor at beginning of lump data
			_, err = f.Seek(int64(dir.DataOffset), io.SeekStart)
			if err != nil {
				return nil, err
			}

			// Read lump data
			lumpData = make([]byte, dir.DataLength)
			err = binary.Read(f, binary.LittleEndian, &lumpData)
			if err != nil {
				return nil, err
			}

			// Put cursor back where it was
			_, err = f.Seek(currentOffset, io.SeekStart)
			if err != nil {
				return nil, err
			}
		} else {
			// If not, no lump data to store
			lumpData = []byte{}
		}

		// Build the lump object
		lump := Lump{
			Name: strings.Trim(string(dir.LumpName[:]), "\x00"),
			Data: lumpData,
		}

		if lump.IsMapHeader() {
			// fmt.Println("Found map:", lump.Name)
			level := LevelMap{
				Header:       &lump,
				MapDataLumps: []Lump{},
			}
			currentMap = &level
		} else if currentMap != nil {
			currentMap.MapDataLumps = append(currentMap.MapDataLumps, lump)

			if lump.Name == "BLOCKMAP" {
				wad.Maps = append(wad.Maps, *currentMap)
				currentMap = nil
			}
		} else {
			// TODO: Detect lump type
			wad.UnknownLumps = append(wad.UnknownLumps, lump)
		}

	}

	return &wad, nil
}
