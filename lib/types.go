package lib

type FileHeaderInfo struct {
	NameLength uint16
	Size       int64
	Offset     int64
	Mode       uint32
}

// type FileHeader struct {
// 	Name string
// 	Info FileHeaderInfo
// }

// func (f *FileHeader) String() string {
// 	return
// }
