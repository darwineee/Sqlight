package cell

import _type "com.sentry.dev/app/sqlight/type"

// InteriorIndex represents an interior cell in an index B-tree (type 0x02)
type InteriorIndex struct {
	LeftChildPage uint32       // 4-byte page number which is the left child pointer
	PayloadSize   _type.VarInt // VarInt: Total number of bytes of key payload
	Payload       Record       // Initial portion of the payload (the key)
	OverflowPage  uint32       // Optional: Page number of first overflow page (if payload doesn't fit)
}
