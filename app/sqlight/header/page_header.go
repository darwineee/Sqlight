package header

import (
	"com.sentry.dev/app/sqlight/type"
	"encoding/binary"
	"fmt"
)

const (
	PageHeaderSizeInterior = 12
	PageHeaderSizeLeaf     = 8
)

func (h *PageHeader) Size() int {
	if h.PageType.IsInteriorPage() {
		return PageHeaderSizeInterior
	}
	return PageHeaderSizeLeaf
}

// PageHeader represents the header structure of a B-tree page
type PageHeader struct {
	PageType            _type.Page
	FirstFreeBlockPtr   uint16
	CellCount           uint16
	CellContentPtr      uint16
	FragmentedFreeBytes byte
	RightmostPtr        uint32 // Only present in interior pages
}

// OfPage parses a B-tree page header from a byte slice
func OfPage(data []byte) (*PageHeader, error) {
	if len(data) < PageHeaderSizeLeaf {
		return nil, fmt.Errorf("data too short: need at least 8 bytes, got %d", len(data))
	}

	pageType := _type.Page(data[0])
	if !pageType.IsValid() {
		return nil, fmt.Errorf("invalid page type: 0x%02x", data[0])
	}

	header := &PageHeader{
		PageType:            pageType,
		FirstFreeBlockPtr:   binary.BigEndian.Uint16(data[1:3]),
		CellCount:           binary.BigEndian.Uint16(data[3:5]),
		CellContentPtr:      binary.BigEndian.Uint16(data[5:7]),
		FragmentedFreeBytes: data[7],
	}

	if pageType.IsInteriorPage() {
		if len(data) < PageHeaderSizeInterior {
			return nil, fmt.Errorf("data too short for interior page: need 12 bytes, got %d", len(data))
		}
		header.RightmostPtr = binary.BigEndian.Uint32(data[8:12])
	}

	return header, nil
}

func (h *PageHeader) GetRealCellContentPointer() uint {
	if h.CellContentPtr == 0 {
		return 65536
	} else {
		return uint(h.CellContentPtr)
	}
}

func (h *PageHeader) IsRemainFreeBlock() bool {
	return h.FirstFreeBlockPtr != 0
}
