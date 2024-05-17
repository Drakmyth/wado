package wad

import (
	"encoding/binary"
	"io"
	"os"
)

type WadFile struct {
	file       *os.File
	Identifier string
	Lumps      []Lump
}

type Lump struct {
	Name string
	Data []byte
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

	lumps := make([]Lump, 0, header.LumpCount)
	for _, dir := range directory {
		lumpData, err := parseLumpData(f, dir.DataOffset, dir.DataLength)
		if err != nil {
			return nil, err
		}

		lumps = append(lumps, Lump{
			Name: NameToStr(dir.LumpName[:]),
			Data: lumpData,
		})
	}

	return &WadFile{
		file:       f,
		Identifier: string(header.Identifier[:]),
		Lumps:      lumps,
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
