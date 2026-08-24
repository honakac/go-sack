package lib

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"unsafe"
)

type Reader struct {
	r     io.ReadSeeker
	Files []FileHeader
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

	var tocStartOffset int64
	if err := binary.Read(r.r, binary.LittleEndian, &tocStartOffset); err != nil {
		return err
	}

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

	var count uint32
	if err := binary.Read(r.r, binary.LittleEndian, &count); err != nil {
		return err
	}

	r.Files = make([]FileHeader, count)
	for i := range count {
		if err := binary.Read(r.r, binary.LittleEndian, &r.Files[i].Info); err != nil {
			return err
		}

		name := make([]byte, r.Files[i].Info.NameLength)
		n, err := r.r.Read(name)
		if err != nil {
			return err
		}
		if uint16(n) != r.Files[i].Info.NameLength {
			return fmt.Errorf("invalid read name, file size(%d) != buffer(%d)",
				n, r.Files[i].Info.Size)
		}

		r.Files[i].Name = unsafe.String(unsafe.SliceData(name), len(name))
	}

	return nil
}
