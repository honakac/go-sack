package lib

import (
	"encoding/binary"
	"fmt"
	"io"
	"unsafe"
)

type Writer struct {
	w       io.Writer
	offset  int64
	entries []FileHeader
}

func NewWriter(w io.Writer) (*Writer, error) {
	_, err := w.Write([]byte{0xEE})
	if err != nil {
		return nil, err
	}

	return &Writer{w: w, offset: 1}, nil
}

func (w *Writer) AddStream(name string, size int64, mode uint32, r io.Reader) error {
	header := FileHeader{
		Name: name,
		Info: FileHeaderInfo{
			NameLength: uint16(len(name)),
			Size:       size,
			Mode:       mode,
			Offset:     w.offset,
		},
	}

	written, err := io.Copy(w.w, r)
	if err != nil {
		return err
	}
	if written != size {
		return fmt.Errorf("written %d instead of %d", written, size)
	}

	w.offset += written
	w.entries = append(w.entries, header)

	return nil
}

func (w *Writer) writeTOC() error {
	tocStartOffset := w.offset

	if err := binary.Write(w.w, binary.LittleEndian, uint32(len(w.entries))); err != nil {
		return err
	}

	for _, e := range w.entries {
		if err := binary.Write(w.w, binary.LittleEndian, e.Info); err != nil {
			return err
		}
		if _, err := w.w.Write(unsafe.Slice(unsafe.StringData(e.Name), len(e.Name))); err != nil {
			return err
		}
	}

	if err := binary.Write(w.w, binary.LittleEndian, tocStartOffset); err != nil {
		return err
	}

	if _, err := w.w.Write([]byte("SACK")); err != nil {
		return err
	}

	return nil
}

func (w *Writer) Close() error {
	return w.writeTOC()
}
