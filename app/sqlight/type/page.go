package _type

type Page byte

const (
	InteriorIndex Page = 0x02
	InteriorTable Page = 0x05
	LeafIndex     Page = 0x0a
	LeafTable     Page = 0x0d
)

func (t Page) IsValid() bool {
	switch t {
	case InteriorIndex, InteriorTable, LeafIndex, LeafTable:
		return true
	default:
		return false
	}
}

func (t Page) IsInteriorPage() bool {
	return t == InteriorIndex || t == InteriorTable
}
