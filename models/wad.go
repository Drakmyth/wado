package models

type WadHeader struct {
	Identification  [4]byte
	NumLumps        int32
	InfoTableOffset int32
}
