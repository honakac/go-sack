package sack

type FileHeaderInfo struct {
	NameLength uint16
	Size       int64
	Offset     int64
	Mode       uint32
}
