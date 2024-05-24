package wad

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
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

func CreateFile(filepath string) (*WadFile, error) {
	f, err := os.OpenFile(filepath, os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &WadFile{
		file:       f,
		Identifier: "PWAD",
		Lumps:      []Lump{},
		Levels:     []Level{},
	}, nil
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
			Name: nameToStr(dir.LumpName[:]),
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
	levelInfos := make([]LevelInfo, 0, len(wf.Levels))

	for _, level := range wf.Levels {
		levelLumps := level.toLumps()
		lumps = append(lumps, levelLumps...)
		levelInfos = append(levelInfos, level.LevelInfo)
	}

	lumps = append(lumps, wf.Lumps...)
	lumps = append(lumps, makeUMapInfoLump(levelInfos))

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
		copy(lumpName[:], strToName(lump.Name))
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

func makeUMapInfoLump(levelInfos []LevelInfo) Lump {
	builder := strings.Builder{}
	for i, levelInfo := range levelInfos {
		levelSlot := fmt.Sprintf("MAP%02d", i+1)
		lastSlot := strconv.FormatBool(levelInfo.EndGame)

		builder.WriteString(fmt.Sprintf(
			`MAP %s
{
    levelname = "%s"
    label = "%s"
`, levelSlot, levelInfo.Name, levelInfo.Label))
		if !levelInfo.EndGame {
			builder.WriteString(fmt.Sprintf(
				`    next = "%s"
    nextsecret = "%s"
`, levelInfo.Next, levelInfo.NextSecret))
		}
		builder.WriteString(fmt.Sprintf(
			`    intertext = clear
    intertextsecret = clear
    endgame = %s
    endcast = %s
    nointermission = %s
    bossaction = clear
`, lastSlot, lastSlot, lastSlot))

		for _, bossAction := range levelInfo.BossActions {
			builder.WriteString(fmt.Sprintf("    bossaction = %s\n", bossAction))
		}
		builder.WriteString("}\n\n")
	}

	mapInfoStr := builder.String()
	return Lump{
		Name: "UMAPINFO",
		Data: []byte(mapInfoStr),
	}
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
