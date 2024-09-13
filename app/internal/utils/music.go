package utils

import (
	"os"
)

type Chunk []byte

type MusicReader struct {
	file *os.File
}

func NewMusicReader(path string) (*MusicReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &MusicReader{
		file: file,
	}, nil
}

func (m *MusicReader) Read(p []byte) (int, error) {
	n, err := m.file.Read(p)
	if err != nil {
		return n, err
	}

	return n, nil
}

func (m *MusicReader) Close() error {
	return m.file.Close()
}
