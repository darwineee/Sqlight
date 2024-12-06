package cell

import _type "com.sentry.dev/app/sqlight/type"

// InteriorTable represents an interior cell in a table B-tree (type 0x05)
type InteriorTable struct {
	LeftChildPage uint32       // 4-byte page number which is the left child pointer
	RowID         _type.VarInt // VarInt: Integer key
}
