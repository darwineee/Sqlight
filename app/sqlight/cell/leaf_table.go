package cell

import (
	_type "com.sentry.dev/app/sqlight/type"
	"encoding/binary"
	"fmt"
)

// LeafTable represents a leaf cell in a table B-tree (type 0x0d)
type LeafTable struct {
	PayloadSize  _type.VarInt // VarInt: Total number of bytes of payload
	RowID        _type.VarInt // VarInt: Integer key (rowid)
	Payload      Record       // Payload in record format
	OverflowPage uint32       // Optional: Page number of first overflow page (if payload doesn't fit)
}

func ParseLeafTable(data []byte) (LeafTable, error) {
	var table LeafTable
	var pos int

	// 1. Read total payload size
	payloadSize, bytesRead, err := _type.ReadVarInt(data)
	if err != nil {
		return LeafTable{}, fmt.Errorf("reading payload size: %w", err)
	}
	table.PayloadSize = payloadSize
	pos += bytesRead

	// 2. Read rowid
	rowID, bytesRead, err := _type.ReadVarInt(data[pos:])
	if err != nil {
		return LeafTable{}, fmt.Errorf("reading rowid: %w", err)
	}
	table.RowID = rowID
	pos += bytesRead

	// 3. Parse the record itself
	record, err := ParseRecord(data[pos:])
	if err != nil {
		return LeafTable{}, fmt.Errorf("parsing record: %w", err)
	}
	table.Payload = record

	// 4. Optional: Check for overflow page
	// If payload doesn't fit in this page, there will be a 4-byte overflow page number
	actualPayloadSize := len(record.Values[0]) // Sum up all value lengths
	for i := 1; i < len(record.Values); i++ {
		actualPayloadSize += len(record.Values[i])
	}

	if actualPayloadSize < int(payloadSize) {
		// We have overflow
		pos += actualPayloadSize
		if pos+4 > len(data) {
			return LeafTable{}, fmt.Errorf("data too short for overflow page")
		}
		table.OverflowPage = binary.BigEndian.Uint32(data[pos : pos+4])
	}

	return table, nil
}
