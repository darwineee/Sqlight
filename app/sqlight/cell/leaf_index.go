package cell

import _type "com.sentry.dev/app/sqlight/type"

// LeftIndex represents a leaf cell in an index B-tree (type 0x0a)
type LeftIndex struct {
	PayloadSize  _type.VarInt // VarInt: Total number of bytes of key payload
	Payload      Record       // Initial portion of the payload (the key)
	OverflowPage uint32       // Optional: Page number of first overflow page (if payload doesn't fit)
}
