package irp

const (
	FrequencyKaseikyo = 36700
)

func NewIrpKaseikyo(protocol string) Irp {
	return &irpImpl{
		protocol:  protocol,
		frequency: FrequencyKaseikyo,
		decode:    DecodeKaseikyo,
	}
}

func DecodeKaseikyo(code SignalCode) (SignalData, error) {
	vendor := (uint16(code.Address[2]) << 8) | uint16(code.Address[1])
	address := code.Address[0]
	command := (uint16(code.Command[1]) << 8) | uint16(code.Command[0])

	vendorParity := uint8(vendor ^ (vendor >> 8))
	vendorParity = (vendorParity ^ (vendorParity >> 4)) & 0xF

	lowWord := uint16(address) << KASEIKYO_VENDOR_ID_PARITY_BITS
	lowByte := uint8(lowWord >> 8)
	lowMidByte := uint8(lowWord & 0xff)
	highMidByte := uint8(0)
	highByte := uint8(0)

	lowByte |= vendorParity
	highMidByte = uint8(command)
	highByte = uint8(command ^ uint16(lowByte) ^ uint16(lowMidByte))

	ulong := (uint32(lowByte) << 24) |
		(uint32(lowMidByte) << 16) |
		(uint32(highMidByte) << 8) |
		(uint32(highByte))

	data := NewSignalData()
	data.Add(KASEIKYO_PREAMBLE_MARK, KASEIKYO_PREAMBLE_SPACE)
	data.AddBits16(GetBits16(vendor, true), KASEIKYO_BIT1_MARK, KASEIKYO_BIT1_SPACE, KASEIKYO_BIT0_MARK, KASEIKYO_BIT0_SPACE)
	data.AddBits32(GetBits32(ulong, false), KASEIKYO_BIT1_MARK, KASEIKYO_BIT1_SPACE, KASEIKYO_BIT0_MARK, KASEIKYO_BIT0_SPACE)
	data.Add(KASEIKYO_BIT1_MARK)

	return data, nil
}

const (
	KASEIKYO_VENDOR_ID_PARITY_BITS = 4
	KASEIKYO_UNIT                  = 432
	KASEIKYO_PREAMBLE_MARK         = (8 * KASEIKYO_UNIT)
	KASEIKYO_PREAMBLE_SPACE        = (4 * KASEIKYO_UNIT)
	KASEIKYO_BIT1_MARK             = KASEIKYO_UNIT
	KASEIKYO_BIT1_SPACE            = (3 * KASEIKYO_UNIT)
	KASEIKYO_BIT0_MARK             = KASEIKYO_UNIT
	KASEIKYO_BIT0_SPACE            = KASEIKYO_UNIT
	KASEIKYO_REPEAT_PERIOD         = 130000
	KASEIKYO_SILENCE               = KASEIKYO_REPEAT_PERIOD
	KASEIKYO_MIN_SPLIT_TIME        = KASEIKYO_REPEAT_PAUSE_MIN
	KASEIKYO_REPEAT_PAUSE_MIN      = 4000
	KASEIKYO_REPEAT_PAUSE_MAX      = 150000
	KASEIKYO_REPEAT_MARK           = KASEIKYO_PREAMBLE_MARK
	KASEIKYO_REPEAT_SPACE          = (KASEIKYO_REPEAT_PERIOD - 56000)
	KASEIKYO_PREAMBLE_TOLERANCE    = 200
	KASEIKYO_BIT_TOLERANCE         = 120
)

/*
mark(KASEIKYO_HEADER_MARK);   (8 * KASEIKYO_UNIT) // 3456
space(KASEIKYO_HEADER_SPACE); (4 * KASEIKYO_UNIT) // 1728


sendPulseDistanceWidthData(
aOneMarkMicros   KASEIKYO_BIT_MARK,           KASEIKYO_UNIT       // 432
aOneSpaceMicros  KASEIKYO_ONE_SPACE,          (3 * KASEIKYO_UNIT) // 1296
aZeroMarkMicros  KASEIKYO_BIT_MARK,           KASEIKYO_UNIT       // 432
aZeroSpaceMicros KASEIKYO_ZERO_SPACE,         KASEIKYO_UNIT       // 432
aVendorCode,                 code
KASEIKYO_VENDOR_ID_BITS,     16
PROTOCOL_IS_LSB_FIRST,       true
SEND_NO_STOP_BIT             false
);


                |   0  |   1   |   2    |   3   |   4    |   5   |   6    |   7   |   8   |   9    |   10  |   11  |   12   |   13   |   14  |  15   |       |       |       |
                |   0  |   0   |   1    |   0   |   1    |   0   |   1    |   0   |   0   |   1    |   0   |   0   |   1    |    1   |   0   |   0   |       |       |       |
data: 3363 1685 407 436 411 432 415 1240 434 410 437 1245 439 404 433 1249 435 408 439 431 406 1249 435 435 412 405 442 1241 433 1249 435 408 439 405 442 428 409 434 413 430 407 411 436 433 414 429 408 1248 436 407 440 1243 441 428 409 434 413 431 406 1249 435 1248 436 406 441 1242 442 1240 434 409 438 431 416 428 409 408 439 430 407 411 436 407 440 429 408 436 411 432 415 402 435 1247 437 1245 439 1243 441 1238 436
protocol: Kaseikyo
address: 41 54 32 00
command: 1B 00 00 00

0010 1010 0100 1100

0011 0010 0101 0100
3      2    5    4
*/
