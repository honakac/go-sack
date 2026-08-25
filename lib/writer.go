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
	buf     [22]byte
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

	binary.LittleEndian.PutUint32(w.buf[0:4], uint32(len(w.entries)))
	if _, err := w.w.Write(w.buf[0:4]); err != nil {
		return err
	}

	for _, e := range w.entries {
		binary.LittleEndian.PutUint16(w.buf[0:2], e.Info.NameLength)
		binary.LittleEndian.PutUint64(w.buf[2:10], uint64(e.Info.Size))
		binary.LittleEndian.PutUint64(w.buf[10:18], uint64(e.Info.Offset))
		binary.LittleEndian.PutUint32(w.buf[18:22], e.Info.Mode)

		if _, err := w.w.Write(w.buf[0:22]); err != nil {
			return err
		}
		if _, err := w.w.Write(unsafe.Slice(unsafe.StringData(e.Name), len(e.Name))); err != nil {
			return err
		}
	}

	binary.LittleEndian.PutUint64(w.buf[0:8], uint64(tocStartOffset))
	if _, err := w.w.Write(w.buf[0:8]); err != nil {
		return err
	}

	if _, err := w.w.Write(Magic[0:4]); err != nil {
		return err
	}

	return nil
}

func (w *Writer) Close() error {
	return w.writeTOC()
}
