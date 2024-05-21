package wad

import (
	"encoding/binary"
	"io"
	"os"
	"regexp"
)

type WadFile struct {
	file       *os.File
	Identifier string
	Lumps      []Lump
	Levels     []Level
}

type Lump struct {
	Name string
	Data []byte
}

type Game int

const (
	GAME_DOOM Game = iota
	GAME_DOOM2
)

func makeThingsLump(things []Thing) Lump {
	return Lump{
		Name: LUMP_THINGS,
		Data: MarshalThings(things),
	}
}

func makeSidedefsLump(sidedefs []Sidedef) Lump {
	return Lump{
		Name: LUMP_SIDEDEFS,
		Data: MarshalSidedefs(sidedefs),
	}
}

func isLevelFromGame(name string, game Game) bool {
	switch game {
	case GAME_DOOM:
		d1LevelNameRegexp := regexp.MustCompile(`^E(\d)M(\d)$`)
		return d1LevelNameRegexp.MatchString(name)
	case GAME_DOOM2:
		d2LevelNameRegexp := regexp.MustCompile(`^MAP(\d+)$`)
		return d2LevelNameRegexp.MatchString(name)
	}

	return false
}

func OpenFile(filepath string) (*WadFile, error) {
	f, err := os.OpenFile(filepath, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	header, err := parseHeader(f)
	if err != nil {
		return nil, err
	}

	directory, err := parseDirectory(f, header.DirectoryOffset, header.LumpCount)
	if err != nil {
		return nil, err
	}

	levels := make([]Level, 0, 9)
	lumps := make([]Lump, 0, header.LumpCount)
	for i := 0; i < len(directory); i++ {
		dir := directory[i]
		lumpData, err := parseLumpData(f, dir.DataOffset, dir.DataLength)
		if err != nil {
			return nil, err
		}

		lump := Lump{
			Name: NameToStr(dir.LumpName[:]),
			Data: lumpData,
		}

		if isLevelFromGame(lump.Name, GAME_DOOM) || isLevelFromGame(lump.Name, GAME_DOOM2) {
			level, err := parseLevel(f, lump.Name, directory[i+1:i+11])
			if err != nil {
				return nil, err
			}
			levels = append(levels, level)
			i += 10
		} else {
			lumps = append(lumps, lump)
		}
	}

	return &WadFile{
		file:       f,
		Identifier: string(header.Identifier[:]),
		Lumps:      lumps,
		Levels:     levels,
	}, nil
}

func (wf WadFile) Save() error {
	_, err := wf.file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	err = wf.file.Truncate(0)
	if err != nil {
		return err
	}

	f := wf.file

	lumps := make([]Lump, 0, len(wf.Lumps)+len(wf.Levels)*11)

	for _, level := range wf.Levels {
		levelLumps := level.toLumps()
		lumps = append(lumps, levelLumps...)
	}

	lumps = append(lumps, wf.Lumps...)

	header := makeHeader(wf.Identifier, lumps)
	err = binary.Write(f, binary.LittleEndian, header)
	if err != nil {
		return err
	}

	for _, lump := range lumps {
		err = binary.Write(f, binary.LittleEndian, lump.Data)
		if err != nil {
			return err
		}
	}

	directory := makeDirectory(lumps)
	err = binary.Write(f, binary.LittleEndian, directory)
	if err != nil {
		return err
	}

	return f.Sync()
}

func (wf WadFile) Close() error {
	return wf.file.Close()
}

func makeHeader(identifier string, lumps []Lump) fileHeader {

	directoryOffset := SIZE_HEADER
	for _, lump := range lumps {
		directoryOffset += len(lump.Data)
	}

	ident := [4]byte{}
	copy(ident[:], []byte(identifier))

	return fileHeader{
		Identifier:      ident,
		LumpCount:       int32(len(lumps)),
		DirectoryOffset: int32(directoryOffset),
	}
}

func makeDirectory(lumps []Lump) []fileDirectoryEntry {
	directory := make([]fileDirectoryEntry, 0, len(lumps))

	offset := SIZE_HEADER
	for _, lump := range lumps {
		lumpName := [8]byte{}
		copy(lumpName[:], StrToName(lump.Name))
		lumpSize := len(lump.Data)

		dir := fileDirectoryEntry{
			DataOffset: int32(offset),
			DataLength: int32(lumpSize),
			LumpName:   lumpName,
		}

		directory = append(directory, dir)
		offset += lumpSize
	}

	return directory
}
