package irp

import (
	"strconv"
	"strings"
)

func GetBits8(value uint8, lsb bool) [8]bool {
	const size = 8
	result := [size]bool{}
	for i := 0; i < size; i++ {
		if value&uint8(0x01) == 1 {
			if lsb {
				result[i] = true
			} else {
				result[size-1-i] = true
			}

		}
		value >>= 1
	}
	return result
}

func GetBits16(value uint16, lsb bool) [16]bool {
	result := [16]bool{}
	if lsb {
		source := GetBits8(uint8(value), true)
		copy(result[:8], source[:])
		source = GetBits8(uint8(value>>8), true)
		copy(result[8:], source[:])
	} else {
		source := GetBits8(uint8(value>>8), false)
		copy(result[:8], source[:])
		source = GetBits8(uint8(value), false)
		copy(result[8:], source[:])
	}
	return result
}

func GetBits32(value uint32, lsb bool) [32]bool {
	result := [32]bool{}
	if lsb {
		source := GetBits16(uint16(value), true)
		copy(result[:16], source[:])
		source = GetBits16(uint16(value>>16), true)
		copy(result[16:], source[:])
	} else {
		source := GetBits16(uint16(value>>16), false)
		copy(result[:16], source[:])
		source = GetBits16(uint16(value), false)
		copy(result[16:], source[:])
	}
	return result
}

func GetBytes32(val uint32) [4]uint8 {
	const size = 4
	result := [4]uint8{}
	for i := 0; i < size; i++ {
		result[size-i-1] = uint8(val >> (i * 8))
	}
	return result
}

func ParseHex32(str string) ([4]uint8, error) {
	const size = 4
	result := [size]uint8{}
	data := strings.Split(str, " ")
	if len(data) == 0 {
		return result, Errorf("bad hex32: %v", data)
	}
	if data[0] == "0x" {
		data = data[1:]
	}

	if !(0 < len(data) && len(data) <= size) {
		return result, Errorf("bad hex32: %v", data)
	}

	offset := size - len(data)
	for i, hex8 := range data {
		val, err := strconv.ParseUint(hex8, 16, 8)
		if err != nil {
			return result, Errorf("bad hex32: %w", err)
		}
		result[offset+i] = uint8(val)
	}

	return result, nil
}

func ParseMicros(str string) (Micros, error) {
	res, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, Errorf("bad micros: %w", err)
	}
	return Micros(res), nil
}

func ParseFrequency(str string) (Frequency, error) {
	res, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, Errorf("bad frequency: %w", err)
	}
	return Frequency(res), nil
}

func SplitToMicrosArr(str string, separator string) ([]Micros, error) {
	data := make([]Micros, 0, 0)
	for i, elemStr := range strings.Split(str, separator) {
		elem, err := ParseMicros(strings.TrimSpace(elemStr))
		if err != nil {
			return nil, Errorf("bad micros elem[%v]: %w", i, err)
		}
		data = append(data, elem)
	}
	return data, nil
}
