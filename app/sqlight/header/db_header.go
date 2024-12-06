package header

import (
	"encoding/binary"
	"fmt"
)

const (
	DbHeaderSize = 100
)

// DbHeader represents the database header format
type DbHeader struct {
	HeaderString        [16]byte // Offset 0: "SQLite format 3\000"
	PageSize            uint16   // Offset 16: Database page size
	WriteVersion        uint8    // Offset 18: File format write version
	ReadVersion         uint8    // Offset 19: File format read version
	ReservedSpace       uint8    // Offset 20: Bytes of unused space
	MaxPayloadFraction  uint8    // Offset 21: Maximum payload fraction
	MinPayloadFraction  uint8    // Offset 22: Minimum payload fraction
	LeafPayloadFraction uint8    // Offset 23: Leaf payload fraction
	FileChangeCounter   uint32   // Offset 24: File change counter
	DatabaseSize        uint32   // Offset 28: Size in pages
	FirstFreelistPage   uint32   // Offset 32: First freelist trunk page
	FreelistPages       uint32   // Offset 36: Total freelist pages
	SchemaCookie        uint32   // Offset 40: Schema cookie
	SchemaFormat        uint32   // Offset 44: Schema format number
	DefaultPageCache    uint32   // Offset 48: Default page cache size
	LargestRootBTree    uint32   // Offset 52: Largest root b-tree page
	TextEncoding        uint32   // Offset 56: Database text encoding
	UserVersion         uint32   // Offset 60: User version
	IncrementalVacuum   uint32   // Offset 64: Incremental-vacuum mode flag
	ApplicationID       uint32   // Offset 68: Application ID
	Reserved            [20]byte // Offset 72: Reserved for expansion
	VersionValidFor     uint32   // Offset 92: Version valid for number
	SQLiteVersionNumber uint32   // Offset 96: SQLITE_VERSION_NUMBER
}

// OfDatabase parses a byte slice into a DbHeader struct
func OfDatabase(data []byte) (*DbHeader, error) {
	if len(data) < 100 {
		return nil, fmt.Errorf("data too short: expected 100 bytes, got %d", len(data))
	}

	header := &DbHeader{}
	copy(header.HeaderString[:], data[0:16])

	header.PageSize = binary.BigEndian.Uint16(data[16:18])
	header.WriteVersion = data[18]
	header.ReadVersion = data[19]
	header.ReservedSpace = data[20]
	header.MaxPayloadFraction = data[21]
	header.MinPayloadFraction = data[22]
	header.LeafPayloadFraction = data[23]
	header.FileChangeCounter = binary.BigEndian.Uint32(data[24:28])
	header.DatabaseSize = binary.BigEndian.Uint32(data[28:32])
	header.FirstFreelistPage = binary.BigEndian.Uint32(data[32:36])
	header.FreelistPages = binary.BigEndian.Uint32(data[36:40])
	header.SchemaCookie = binary.BigEndian.Uint32(data[40:44])
	header.SchemaFormat = binary.BigEndian.Uint32(data[44:48])
	header.DefaultPageCache = binary.BigEndian.Uint32(data[48:52])
	header.LargestRootBTree = binary.BigEndian.Uint32(data[52:56])
	header.TextEncoding = binary.BigEndian.Uint32(data[56:60])
	header.UserVersion = binary.BigEndian.Uint32(data[60:64])
	header.IncrementalVacuum = binary.BigEndian.Uint32(data[64:68])
	header.ApplicationID = binary.BigEndian.Uint32(data[68:72])
	copy(header.Reserved[:], data[72:92])
	header.VersionValidFor = binary.BigEndian.Uint32(data[92:96])
	header.SQLiteVersionNumber = binary.BigEndian.Uint32(data[96:100])

	return header, nil
}

func (dbHeader *DbHeader) GetRealPageSize() uint {
	if dbHeader.PageSize == 1 {
		return 65536
	} else {
		return uint(dbHeader.PageSize)
	}
}
