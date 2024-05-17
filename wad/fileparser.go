package wad

import (
	"encoding/binary"
	"io"
	"os"
)

const SIZE_HEADER int = 12

type fileHeader struct {
	Identifier      [4]byte
	LumpCount       int32
	DirectoryOffset int32
}

type fileDirectoryEntry struct {
	DataOffset int32
	DataLength int32
	LumpName   [8]byte
}

func parseHeader(f *os.File) (fileHeader, error) {
	header := fileHeader{}

	// Position cursor at beginning of file
	_, err := f.Seek(0, io.SeekStart)
	if err != nil {
		return header, err
	}

	// Read the file header
	err = binary.Read(f, binary.LittleEndian, &header)
	if err != nil {
		return header, err
	}
	return header, nil
}

func parseDirectory(f *os.File, offset int32, count int32) ([]fileDirectoryEntry, error) {
	directory := make([]fileDirectoryEntry, 0, count)

	// Position cursor at beginning of lump directory
	_, err := f.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return directory, err
	}

	// For each lump...
	for i := int32(0); i < count; i++ {
		// Read the directory entry for this lump
		entry := fileDirectoryEntry{}
		err = binary.Read(f, binary.LittleEndian, &entry)
		if err != nil {
			return directory, err
		}

		directory = append(directory, entry)
	}

	return directory, nil
}

func parseLumpData(f *os.File, offset int32, length int32) ([]byte, error) {
	lumpData := make([]byte, length)

	// Position cursor at beginning of lump data
	_, err := f.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return lumpData, err
	}

	// Read the lump data
	err = binary.Read(f, binary.LittleEndian, lumpData)
	if err != nil {
		return lumpData, err
	}
	return lumpData, nil
}