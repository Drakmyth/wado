package wad

import (
	"encoding/binary"
	"fmt"
)

type Sidedef struct {
	XOffset      int16
	YOffset      int16
	UpperTex     string
	LowerTex     string
	MiddleTex    string
	FacingSector int16
}

func (s *Sidedef) UnmarshalBinary(data []byte) error {
	s.XOffset = int16(binary.LittleEndian.Uint16(data[0:2]))
	s.YOffset = int16(binary.LittleEndian.Uint16(data[2:4]))
	s.UpperTex = nameToStr(data[4:12])
	s.LowerTex = nameToStr(data[12:20])
	s.MiddleTex = nameToStr(data[20:28])
	s.FacingSector = int16(binary.LittleEndian.Uint16(data[28:30]))

	fmt.Println(s)
	return nil
}

func (s Sidedef) MarshalBinary() ([]byte, error) {
	sbytes := [30]byte{}
	binary.LittleEndian.PutUint16(sbytes[0:2], uint16(s.XOffset))
	binary.LittleEndian.PutUint16(sbytes[2:4], uint16(s.YOffset))
	copy(sbytes[4:12], strToName(s.UpperTex))
	copy(sbytes[12:20], strToName(s.LowerTex))
	copy(sbytes[20:28], strToName(s.MiddleTex))
	binary.LittleEndian.PutUint16(sbytes[28:30], uint16(s.FacingSector))

	return sbytes[:], nil
}
