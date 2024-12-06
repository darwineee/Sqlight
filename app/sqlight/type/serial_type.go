package _type

// SerialType values stored in record header as VarInt
const (
	SerialTypeNull VarInt = 0

	SerialTypeInt8  VarInt = 1 // 8-bit signed integer
	SerialTypeInt16 VarInt = 2 // 16-bit signed integer
	SerialTypeInt24 VarInt = 3 // 24-bit signed integer
	SerialTypeInt32 VarInt = 4 // 32-bit signed integer
	SerialTypeInt48 VarInt = 5 // 48-bit signed integer
	SerialTypeInt64 VarInt = 6 // 64-bit signed integer

	SerialTypeFloat64 VarInt = 7 // 64-bit IEEE floating point

	SerialTypeFalse VarInt = 8 // Represents integer value 0
	SerialTypeTrue  VarInt = 9 // Represents integer value 1

	SerialTypeBLOBMin VarInt = 12 // N>=12 and even: BLOB value of (N-12)/2 bytes
	SerialTypeTextMin VarInt = 13 // N>=13 and odd: TEXT value of (N-13)/2 bytes
)

func GetContentTypeSize(t VarInt) int {
	switch {
	case t == SerialTypeNull:
		return 0
	case t == SerialTypeInt8:
		return 1
	case t == SerialTypeInt16:
		return 2
	case t == SerialTypeInt24:
		return 3
	case t == SerialTypeInt32:
		return 4
	case t == SerialTypeInt48:
		return 6
	case t == SerialTypeInt64:
		return 8
	case t == SerialTypeFloat64:
		return 8
	case t == SerialTypeFalse:
		return 0
	case t == SerialTypeTrue:
		return 0
	case t >= SerialTypeBLOBMin && t%2 == 0:
		return GetBlobSize(t)
	case t >= SerialTypeTextMin && t%2 == 1:
		return GetTextSize(t)
	default:
		return 0 // Invalid or reserved types
	}
}

func IsValidSerialType(t VarInt) bool {
	return (t >= 0 && t <= 9) || t >= 12
}

func IsBlobType(t VarInt) bool {
	return t >= SerialTypeBLOBMin && t%2 == 0
}

func IsTextType(t VarInt) bool {
	return t >= SerialTypeTextMin && t%2 == 1
}

func GetBlobSize(t VarInt) int {
	if IsBlobType(t) {
		return int((t - SerialTypeBLOBMin) / 2)
	}
	return -1
}

func GetTextSize(t VarInt) int {
	if IsTextType(t) {
		return int((t - SerialTypeTextMin) / 2)
	}
	return -1
}

func BlobSerialType(length int) VarInt {
	return SerialTypeBLOBMin + VarInt(length*2)
}

func TextSerialType(length int) VarInt {
	return SerialTypeTextMin + VarInt(length*2)
}
