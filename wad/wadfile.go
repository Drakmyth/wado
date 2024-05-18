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

const (
	LEVEL_HEADER int = iota
	LEVEL_THINGS
	LEVEL_LINEDEFS
	LEVEL_SIDEDEFS
	LEVEL_VERTEXES
	LEVEL_SEGS
	LEVEL_SSECTORS
	LEVEL_NODES
	LEVEL_SECTORS
	LEVEL_REJECT
	LEVEL_BLOCKMAP
)

type Game int

const (
	GAME_DOOM Game = iota
	GAME_DOOM2
)

type Level struct {
	Name  *string
	Lumps []Lump
}

func (l Level) IsLevelFromGame(game Game) bool {
	return isLevelFromGame(l.Lumps[LEVEL_HEADER].Name, game)
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

	levelHeaderIndices := make([]int, 0, 9)
	lumps := make([]Lump, 0, header.LumpCount)
	for i, dir := range directory {
		lumpData, err := parseLumpData(f, dir.DataOffset, dir.DataLength)
		if err != nil {
			return nil, err
		}

		lump := Lump{
			Name: NameToStr(dir.LumpName[:]),
			Data: lumpData,
		}
		lumps = append(lumps, lump)

		if isLevelFromGame(lump.Name, GAME_DOOM) || isLevelFromGame(lump.Name, GAME_DOOM2) {
			levelHeaderIndices = append(levelHeaderIndices, i)
		}
	}

	levels := make([]Level, 0, len(levelHeaderIndices))
	for _, i := range levelHeaderIndices {
		levelLumps := lumps[i : i+11]
		level := Level{
			Name:  &levelLumps[LEVEL_HEADER].Name,
			Lumps: levelLumps,
		}

		levels = append(levels, level)
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

	header := wf.makeHeader()
	err = binary.Write(f, binary.LittleEndian, header)
	if err != nil {
		return err
	}

	for _, lump := range wf.Lumps {
		err = binary.Write(f, binary.LittleEndian, lump.Data)
		if err != nil {
			return err
		}
	}

	directory := wf.makeDirectory()
	err = binary.Write(f, binary.LittleEndian, directory)
	if err != nil {
		return err
	}

	return f.Sync()
}

func (wf WadFile) Close() error {
	return wf.file.Close()
}

func (wf WadFile) makeHeader() fileHeader {

	directoryOffset := SIZE_HEADER
	for _, lump := range wf.Lumps {
		directoryOffset += len(lump.Data)
	}

	ident := [4]byte{}
	copy(ident[:], []byte(wf.Identifier))

	return fileHeader{
		Identifier:      ident,
		LumpCount:       int32(len(wf.Lumps)),
		DirectoryOffset: int32(directoryOffset),
	}
}

func (wf WadFile) makeDirectory() []fileDirectoryEntry {
	directory := make([]fileDirectoryEntry, 0, len(wf.Lumps))

	offset := SIZE_HEADER
	for _, lump := range wf.Lumps {
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
