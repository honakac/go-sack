package lib

import (
	"encoding/binary"
	"errors"
	"io"
	"unsafe"
)

type Reader struct {
	r     io.ReadSeeker
	Files map[string]FileHeaderInfo
}

func NewReader(r io.ReadSeeker) (*Reader, error) {
	reader := &Reader{r: r}

	if err := reader.readTOC(); err != nil {
		return nil, err
	}

	return reader, nil
}

func (r *Reader) OpenFile(header FileHeaderInfo) (io.Reader, error) {
	_, err := r.r.Seek(header.Offset, io.SeekStart)
	if err != nil {
		return nil, err
	}

	return io.LimitReader(r.r, header.Size), nil
}

// Read Table of Contents
func (r *Reader) readTOC() error {
	var buf [22]byte

	// 12 bytes
	// TOC offset (8 bytes) + SACK (4 bytes)
	if _, err := r.r.Seek(-int64(len(Magic))-8, io.SeekEnd); err != nil {
		return err
	}

	if _, err := io.ReadFull(r.r, buf[0:8]); err != nil {
		return err
	}
	tocStartOffset := int64(binary.LittleEndian.Uint64(buf[0:8]))

	var magic [4]byte
	if _, err := io.ReadFull(r.r, magic[:]); err != nil {
		return err
	}
	if magic != Magic {
		return errors.New("invalid magic header")
	}

	if _, err := r.r.Seek(tocStartOffset, io.SeekStart); err != nil {
		return err
	}

	if _, err := io.ReadFull(r.r, buf[0:4]); err != nil {
		return err
	}
	count := binary.LittleEndian.Uint32(buf[0:4])

	r.Files = make(map[string]FileHeaderInfo, count)
	for range count {
		if _, err := io.ReadFull(r.r, buf[0:22]); err != nil {
			return err
		}

		info := FileHeaderInfo{
			NameLength: binary.LittleEndian.Uint16(buf[0:2]),
			Size:       int64(binary.LittleEndian.Uint64(buf[2:10])),
			Offset:     int64(binary.LittleEndian.Uint64(buf[10:18])),
			Mode:       binary.LittleEndian.Uint32(buf[18:22]),
		}

		name := make([]byte, info.NameLength)
		if _, err := io.ReadFull(r.r, name); err != nil {
			return err
		}

		r.Files[unsafe.String(unsafe.SliceData(name), len(name))] = info
	}

	return nil
}
