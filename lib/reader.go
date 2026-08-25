package lib

import (
	"encoding/binary"
	"errors"
	"io"
	"unsafe"
)

type Reader struct {
	r     io.ReadSeeker
	Files []FileHeader
	buf   [22]byte
}

func NewReader(r io.ReadSeeker) (*Reader, error) {
	reader := &Reader{r: r}

	if err := reader.readTOC(); err != nil {
		return nil, err
	}

	return reader, nil
}

// func (r *Reader) ReadFile(name string) io.Reader {
// 	for _, f := range r.Files {
// 		if f.Name == name {
// 			_, err := r.r.Seek(header.Offset, io.SeekStart)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 3. Возвращаем LimitReader: он позволит прочитать РОВНО header.Size байт,
// 	// после чего честно вернёт io.EOF, не заглядывая в чужие файлы!
// 	return io.LimitReader(r.r, header.Size), nil
// 		}
// 	}
// 	return nil
// }

func (r *Reader) OpenFile(header FileHeader) (io.Reader, error) {
	_, err := r.r.Seek(header.Info.Offset, io.SeekStart)
	if err != nil {
		return nil, err
	}

	return io.LimitReader(r.r, header.Info.Size), nil
}

// Read Table of Contents
func (r *Reader) readTOC() error {
	// 12 bytes
	// TOC offset (8 bytes) + SACK (4 bytes)
	if _, err := r.r.Seek(-int64(len(Magic))-8, io.SeekEnd); err != nil {
		return err
	}

	if _, err := io.ReadFull(r.r, r.buf[0:8]); err != nil {
		return err
	}
	tocStartOffset := int64(binary.LittleEndian.Uint64(r.buf[0:8]))

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

	if _, err := io.ReadFull(r.r, r.buf[0:4]); err != nil {
		return err
	}
	count := binary.LittleEndian.Uint32(r.buf[0:4])

	r.Files = make([]FileHeader, count)
	for i := range count {
		if _, err := io.ReadFull(r.r, r.buf[0:22]); err != nil {
			return err
		}
		r.Files[i].Info.NameLength = binary.LittleEndian.Uint16(r.buf[0:2])
		r.Files[i].Info.Size = int64(binary.LittleEndian.Uint64(r.buf[2:10]))
		r.Files[i].Info.Offset = int64(binary.LittleEndian.Uint64(r.buf[10:18]))
		r.Files[i].Info.Mode = binary.LittleEndian.Uint32(r.buf[18:22])

		name := make([]byte, r.Files[i].Info.NameLength)
		if _, err := io.ReadFull(r.r, name); err != nil {
			return err
		}

		r.Files[i].Name = unsafe.String(unsafe.SliceData(name), len(name))
	}

	return nil
}
