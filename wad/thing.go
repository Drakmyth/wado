package wad

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/rand"
)

const SIZE_THING int = 10

type Thing struct {
	X     int16
	Y     int16
	Angle int16
	Type  int16
	Flags int16
}

func (t *Thing) UnmarshalBinary(data []byte) error {
	t.X = int16(binary.LittleEndian.Uint16(data[0:2]))
	t.Y = int16(binary.LittleEndian.Uint16(data[2:4]))
	t.Angle = int16(binary.LittleEndian.Uint16(data[4:6]))
	t.Type = int16(binary.LittleEndian.Uint16(data[6:8]))
	t.Flags = int16(binary.LittleEndian.Uint16(data[8:10]))

	return nil
}

func (t Thing) MarshalBinary() ([]byte, error) {
	tbytes := [SIZE_THING]byte{}
	binary.LittleEndian.PutUint16(tbytes[0:2], uint16(t.X))
	binary.LittleEndian.PutUint16(tbytes[2:4], uint16(t.Y))
	binary.LittleEndian.PutUint16(tbytes[4:6], uint16(t.Angle))
	binary.LittleEndian.PutUint16(tbytes[6:8], uint16(t.Type))
	binary.LittleEndian.PutUint16(tbytes[8:10], uint16(t.Flags))

	return tbytes[:], nil
}

func UnmarshalThings(things []Thing, data []byte) {
	buf := bytes.NewBuffer(data)
	for i, t := range things {
		tbytes := buf.Next(SIZE_THING)
		t.UnmarshalBinary(tbytes)
		things[i] = t
	}
}

func MarshalThings(things []Thing) []byte {
	buf := make([]byte, 0, len(things)*SIZE_THING)
	for _, t := range things {
		tbytes, _ := t.MarshalBinary()
		buf = append(buf, tbytes...)
	}

	return buf
}

func ReplaceThingsWeighted(candidates []*Thing, weights map[int16]float64) {
	// Build bag of replacements to replace candidates with according to weights
	replacements := []int16{}
	for k, v := range weights {
		cnt := int16(math.Round(float64(len(candidates)) * v))
		replacements = append(replacements, repeatedSlice(k, cnt)...)
	}

	executeReplacements(candidates, replacements)
}

func ReplaceThingsCount(candidates []*Thing, counts map[int16]int16) {
	// Build bag of replacements to replace candidates with according to counts
	replacements := []int16{}
	for k, cnt := range counts {
		replacements = append(replacements, repeatedSlice(k, cnt)...)
	}

	executeReplacements(candidates, replacements)
}

func executeReplacements(candidates []*Thing, replacements []int16) {
	// Replace candidates with replacements until we're out of one or the other
	for done := len(replacements) == 0 || len(candidates) == 0; !done; done = len(replacements) == 0 || len(candidates) == 0 {
		// Pick a random index
		candidateIndex := rand.Intn(len(candidates))

		// Replace the candidate
		candidate := candidates[candidateIndex]
		replacementIndex := rand.Intn(len(replacements))
		candidate.Type = replacements[replacementIndex]
		replacements = append(replacements[:replacementIndex], replacements[replacementIndex+1:]...)

		// Remove the index of the replaced candidate from the candidate list
		candidates = append(candidates[:candidateIndex], candidates[candidateIndex:]...)
	}
}

func repeatedSlice[E int | int16](value, n E) []E {
	arr := make([]E, n)
	for i := E(0); i < n; i++ {
		arr[i] = value
	}
	return arr
}
