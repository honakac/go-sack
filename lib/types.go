package lib

import (
	"fmt"
	"os"
)

type FileHeaderInfo struct {
	NameLength uint16
	Size       int64
	Offset     int64
	Mode       uint32
}

type FileHeader struct {
	Name string
	Info FileHeaderInfo
}

func (f *FileHeader) String() string {
	return fmt.Sprintf("%s '%s' %d bytes", os.FileMode(f.Info.Mode), f.Name, f.Info.Size)
}
