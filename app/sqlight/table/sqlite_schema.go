package table

import (
	"com.sentry.dev/app/sqlight/cell"
	"com.sentry.dev/app/sqlight/header"
	_type "com.sentry.dev/app/sqlight/type"
	"encoding/binary"
	"fmt"
)

type SqliteSchema struct {
	DbHeader     *header.DbHeader
	PageHeader   *header.PageHeader
	CellPointers []_type.CellPtr
	CellContent  []cell.SchemaRecord
}

func ParseSqliteSchema(data []byte) (SqliteSchema, error) {
	if len(data) < 100 {
		return SqliteSchema{}, fmt.Errorf("data too short for database header")
	}

	// Parse database header (first 100 bytes)
	dbHeader, err := header.OfDatabase(data[:header.DbHeaderSize])
	if err != nil {
		return SqliteSchema{}, fmt.Errorf("parsing db header: %w", err)
	}

	// Parse page header (next 8 bytes as it's always a leaf page)
	cellPtrStart := header.DbHeaderSize + header.PageHeaderSizeLeaf
	pageHeader, err := header.OfPage(data[header.DbHeaderSize:cellPtrStart])
	if err != nil {
		return SqliteSchema{}, fmt.Errorf("parsing page header: %w", err)
	}

	// Validate it's a leaf table page (should be 0x0d)
	if pageHeader.PageType != 0x0d {
		return SqliteSchema{}, fmt.Errorf("unexpected page type for first page: %x", pageHeader.PageType)
	}

	// Parse cell pointers array
	cellPtrCount := int(pageHeader.CellCount)
	cellPointers := make([]_type.CellPtr, cellPtrCount)

	for i := 0; i < cellPtrCount; i++ {
		offset := cellPtrStart + (i * _type.CellPtrSize)
		ptr := binary.BigEndian.Uint16(data[offset : offset+2])
		cellPointers[i] = _type.CellPtr(ptr)
	}

	// Parse each cell into a SchemaRecord
	schemaRecords := make([]cell.SchemaRecord, cellPtrCount)
	for i, ptr := range cellPointers {
		leafTable, err := cell.ParseLeafTable(data[ptr:])
		if err != nil {
			return SqliteSchema{}, fmt.Errorf("parsing leaf table at offset %d: %w", ptr, err)
		}

		schemaRecord, err := cell.RecordToSchema(leafTable.Payload)
		if err != nil {
			return SqliteSchema{}, fmt.Errorf("converting to schema record at offset %d: %w", ptr, err)
		}

		schemaRecords[i] = schemaRecord
	}

	return SqliteSchema{
		DbHeader:     dbHeader,
		PageHeader:   pageHeader,
		CellPointers: cellPointers,
		CellContent:  schemaRecords,
	}, nil
}
