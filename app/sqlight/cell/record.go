package cell

import (
	_type "com.sentry.dev/app/sqlight/type"
	"encoding/binary"
	"fmt"
)

const (
	SchemaTypeIndex     = 0
	SchemaNameIndex     = 1
	SchemaTblNameIndex  = 2
	SchemaRootPageIndex = 3
	SchemaSqlIndex      = 4
)

type Record struct {
	HeaderSize _type.VarInt   // Total bytes in header
	Types      []_type.VarInt // Serial type for each column
	Values     [][]byte       // Actual column values
}

type SchemaRecord struct {
	Type     _type.SchemaType
	Name     string
	TblName  string
	RootPage int64
	SQL      string
}

func ParseRecord(data []byte) (Record, error) {
	// First read header size
	headerSize, bytesRead, err := _type.ReadVarInt(data)
	if err != nil {
		return Record{}, fmt.Errorf("reading header size: %w", err)
	}

	// Read serial types until we reach header size
	var types []_type.VarInt
	pos := bytesRead
	headerEnd := int(headerSize)

	for pos < headerEnd {
		serialType, n, err := _type.ReadVarInt(data[pos:])
		if err != nil {
			return Record{}, fmt.Errorf("reading serial type: %w", err)
		}
		types = append(types, serialType)
		pos += n
	}

	// Read values based on their serial types
	values := make([][]byte, len(types))
	dataStart := headerEnd

	for i, typ := range types {
		size := _type.GetContentTypeSize(typ)
		if size > 0 {
			if dataStart+size > len(data) {
				return Record{}, fmt.Errorf("data too short for value %d", i)
			}
			values[i] = data[dataStart : dataStart+size]
			dataStart += size
		}
	}

	return Record{
		HeaderSize: headerSize,
		Types:      types,
		Values:     values,
	}, nil
}

func RecordToSchema(r Record) (SchemaRecord, error) {
	if len(r.Values) != 5 {
		return SchemaRecord{}, fmt.Errorf("schema record must have 5 values, got %d", len(r.Values))
	}

	// Parse root page value (should be at index 3)
	rootPage, _ := binary.Varint(r.Values[SchemaRootPageIndex]) //TODO check document to fix it later

	schemaType := string(r.Values[SchemaTypeIndex])
	if !_type.IsValidSchemaType(schemaType) {
		return SchemaRecord{}, fmt.Errorf("invalid schema type: %s", schemaType)
	}

	return SchemaRecord{
		Type:     _type.SchemaType(schemaType),
		Name:     string(r.Values[SchemaNameIndex]),
		TblName:  string(r.Values[SchemaTblNameIndex]),
		RootPage: rootPage,
		SQL:      string(r.Values[SchemaSqlIndex]),
	}, nil
}
