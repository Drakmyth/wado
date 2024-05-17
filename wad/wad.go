package wad

import (
	"fmt"
	"slices"
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

func ShouldShiftTex(sidedef Sidedef) bool {
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
