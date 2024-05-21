package wad

import (
	"bytes"
	"encoding/binary"
)

const SIZE_SIDEDEF int = 30

type Sidedef struct {
	XOffset      int16
	YOffset      int16
	UpperTex     string
	LowerTex     string
	MiddleTex    string
	FacingSector int16
}

func (s *Sidedef) unmarshalBinary(data []byte) error {
	s.XOffset = int16(binary.LittleEndian.Uint16(data[0:2]))
	s.YOffset = int16(binary.LittleEndian.Uint16(data[2:4]))
	s.UpperTex = NameToStr(data[4:12])
	s.LowerTex = NameToStr(data[12:20])
	s.MiddleTex = NameToStr(data[20:28])
	s.FacingSector = int16(binary.LittleEndian.Uint16(data[28:30]))

	return nil
}

func (s Sidedef) marshalBinary() ([]byte, error) {
	sbytes := [SIZE_SIDEDEF]byte{}
	binary.LittleEndian.PutUint16(sbytes[0:2], uint16(s.XOffset))
	binary.LittleEndian.PutUint16(sbytes[2:4], uint16(s.YOffset))
	copy(sbytes[4:12], StrToName(s.UpperTex))
	copy(sbytes[12:20], StrToName(s.LowerTex))
	copy(sbytes[20:28], StrToName(s.MiddleTex))
	binary.LittleEndian.PutUint16(sbytes[28:30], uint16(s.FacingSector))

	return sbytes[:], nil
}

func unmarshalSidedefs(sidedefs []Sidedef, data []byte) {
	buf := bytes.NewBuffer(data)
	for i, s := range sidedefs {
		sbytes := buf.Next(SIZE_SIDEDEF)
		s.unmarshalBinary(sbytes)
		sidedefs[i] = s
	}
}

func marshalSidedefs(sidedefs []Sidedef) []byte {
	buf := make([]byte, 0, len(sidedefs)*SIZE_SIDEDEF)
	for _, s := range sidedefs {
		sbytes, _ := s.marshalBinary()
		buf = append(buf, sbytes...)
	}

	return buf
}
