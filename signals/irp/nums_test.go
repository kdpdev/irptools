package irp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Nums_GetBits8(t *testing.T) {
	assert.Equal(t, [8]bool{true, false, true, false, true, false, true, false}, GetBits8(0b10101010, false))
	assert.Equal(t, [8]bool{false, true, false, true, false, true, false, true}, GetBits8(0b10101010, true))
}

func Test_Nums_GetBits16(t *testing.T) {
	assert.Equal(t, [16]bool{
		true, false, true, false, true, false, true, false,
		true, false, true, false, true, false, true, false,
	}, GetBits16(0b1010101010101010, false))

	assert.Equal(t, [16]bool{
		false, true, false, true, false, true, false, true,
		false, true, false, true, false, true, false, true,
	}, GetBits16(0b1010101010101010, true))
}

func Test_Nums_GetBits32(t *testing.T) {
	assert.Equal(t, [32]bool{
		true, false, true, false, true, false, true, false,
		true, false, true, false, true, false, true, false,
		true, false, true, false, true, false, true, false,
		true, false, true, false, true, false, true, false,
	}, GetBits32(0b10101010101010101010101010101010, false))

	assert.Equal(t, [32]bool{
		false, true, false, true, false, true, false, true,
		false, true, false, true, false, true, false, true,
		false, true, false, true, false, true, false, true,
		false, true, false, true, false, true, false, true,
	}, GetBits32(0b10101010101010101010101010101010, true))
}

func Test_Nums_GetBytes32(t *testing.T) {
	assert.Equal(t, [4]uint8{0x00, 0x00, 0x00, 0x00}, GetBytes32(0x00))
	assert.Equal(t, [4]uint8{0x00, 0x00, 0x00, 0x01}, GetBytes32(0x01))
	assert.Equal(t, [4]uint8{0x00, 0x00, 0x00, 0xab}, GetBytes32(0xab))
	assert.Equal(t, [4]uint8{0x00, 0x00, 0xab, 0xcd}, GetBytes32(0xabcd))
	assert.Equal(t, [4]uint8{0x00, 0xab, 0xcd, 0xef}, GetBytes32(0xabcdef))
	assert.Equal(t, [4]uint8{0x01, 0x23, 0x45, 0x67}, GetBytes32(0x01234567))
}

func Test_Nums_Parse32(t *testing.T) {

	testOk := func(str string, expected uint32) {
		result, err := ParseHex32(str)
		assert.NoError(t, err, "str = '%v'", str)
		assert.Equal(t, GetBytes32(expected), result, "str = '%v'", str)
	}

	testErr := func(str string) {
		result, err := ParseHex32(str)
		assert.Error(t, err, "str = '%v'; res = 0x%x", str, result)
	}

	testErr("")
	testErr("x")
	testErr("xyz")
	testErr("AB xyz")
	testErr("a0 a0 a0 a0 a0")
	testErr("abc")
	testOk("a", 0x0a)
	testOk("a b", 0x0a0b)
	testOk("a b c", 0x0a0b0c)
	testOk("a b c d", 0x0a0b0c0d)
	testOk("aa", 0xaa)
	testOk("aa bb", 0x0000aabb)
	testOk("aa bb cc", 0x00aabbcc)
	testOk("aa bb cc dd", 0xaabbccdd)
	testOk("0x aa", 0x000000aa)
	testOk("0x aa bb", 0x0000aabb)
	testOk("0x aa bb cc", 0x00aabbcc)
	testOk("0x aa bb cc dd", 0xaabbccdd)
}
