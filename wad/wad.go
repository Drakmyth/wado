package wad

import (
	"math"
	"math/rand/v2"
	"slices"
)

const (
	LUMP_THINGS     = "THINGS"
	LUMP_LINEDEFS   = "LINEDEFS"
	LUMP_SIDEDEFS   = "SIDEDEFS"
	LUMP_VERTEXES   = "VERTEXES"
	LUMP_SEGMENTS   = "SEGS"
	LUMP_SUBSECTORS = "SSECTORS"
	LUMP_NODES      = "NODES"
	LUMP_SECTORS    = "SECTORS"
	LUMP_REJECT     = "REJECT"
	LUMP_BLOCKMAP   = "BLOCKMAP"
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

func ReplaceThingsWeighted(candidates []*Thing, weights map[int16]float64, rng *rand.Rand) {
	// Map order is non-deterministic, so sort the keys first
	keys := make([]int16, 0, len(weights))
	for k := range weights {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	// Build bag of replacements to replace candidates with according to weights
	replacements := []int16{}
	for _, k := range keys {
		v := weights[k]
		cnt := int16(math.Round(float64(len(candidates)) * v))
		replacements = append(replacements, repeatedSlice(k, cnt)...)
	}

	executeReplacements(candidates, replacements, rng)
}

func ReplaceThingsCount(candidates []*Thing, counts map[int16]int16, rng *rand.Rand) {
	// Map order is non-deterministic, so sort the keys first
	keys := make([]int16, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	// Build bag of replacements to replace candidates with according to counts
	replacements := []int16{}
	for _, k := range keys {
		cnt := counts[k]
		replacements = append(replacements, repeatedSlice(k, cnt)...)
	}

	executeReplacements(candidates, replacements, rng)
}
