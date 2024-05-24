package wad

import (
	"bytes"
	"encoding/binary"
)

const SIZE_LINEDEF int = 14

var SECRET_EXIT_LINETYPES = []int16{51, 124, 198}

type Linedefs []Linedef
type Linedef struct {
	Start       int16
	End         int16
	Flags       int16
	SpecialType int16
	Tag         int16
	Front       int16
	Back        int16
}

func (l *Linedef) fromBytes(data []byte) {
	l.Start = int16(binary.LittleEndian.Uint16(data[0:2]))
	l.End = int16(binary.LittleEndian.Uint16(data[2:4]))
	l.Flags = int16(binary.LittleEndian.Uint16(data[4:6]))
	l.SpecialType = int16(binary.LittleEndian.Uint16(data[6:8]))
	l.Tag = int16(binary.LittleEndian.Uint16(data[8:10]))
	l.Front = int16(binary.LittleEndian.Uint16(data[10:12]))
	l.Back = int16(binary.LittleEndian.Uint16(data[12:14]))
}

func (l Linedef) toBytes() []byte {
	lbytes := [SIZE_LINEDEF]byte{}
	binary.LittleEndian.PutUint16(lbytes[0:2], uint16(l.Start))
	binary.LittleEndian.PutUint16(lbytes[2:4], uint16(l.End))
	binary.LittleEndian.PutUint16(lbytes[4:6], uint16(l.Flags))
	binary.LittleEndian.PutUint16(lbytes[6:8], uint16(l.SpecialType))
	binary.LittleEndian.PutUint16(lbytes[8:10], uint16(l.Tag))
	binary.LittleEndian.PutUint16(lbytes[10:12], uint16(l.Front))
	binary.LittleEndian.PutUint16(lbytes[12:14], uint16(l.Back))

	return lbytes[:]
}

func parseLinedefs(data []byte) []Linedef {
	numLinedefs := len(data) / SIZE_LINEDEF
	linedefs := make([]Linedef, numLinedefs)

	buf := bytes.NewBuffer(data)
	for i, l := range linedefs {
		lbytes := buf.Next(SIZE_LINEDEF)
		l.fromBytes(lbytes)
		linedefs[i] = l
	}

	return linedefs
}

func (linedefs Linedefs) toLump() Lump {
	buf := make([]byte, 0, len(linedefs)*SIZE_LINEDEF)
	for _, l := range linedefs {
		lbytes := l.toBytes()
		buf = append(buf, lbytes...)
	}

	return Lump{
		Name: LUMP_LINEDEFS,
		Data: buf,
	}
}
