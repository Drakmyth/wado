package wad

import (
	"fmt"
	"strings"
)

type WadFileHeader struct {
	Identifier      [4]byte
	LumpCount       int32
	DirectoryOffset int32
}

type WadDirectoryEntry struct {
	DataOffset int32
	DataLength int32
	LumpName   [8]byte
}

func StrToName(str string) []byte {
	name := [8]byte{}
	paddedName := strings.ReplaceAll(fmt.Sprintf("%-8s", str), " ", "\x00")
	copy(name[:], paddedName)
	return name[:]
}

func NameToStr(name []byte) string {
	return strings.Trim(string(name), "\x00")
}

const (
	LUMP_THINGS   = "THINGS"
	LUMP_SIDEDEFS = "SIDEDEFS"
)

const (
	SIZE_DIRENTRY = 16
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
