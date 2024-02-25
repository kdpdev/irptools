package irp

const (
	FrequencySirc12 = 40000
	FrequencySirc15 = 40000
	FrequencySirc20 = 40000
)

func NewIrpSirc12(protocol string) Irp {
	return &irpImpl{
		protocol:  protocol,
		frequency: FrequencySirc12,
		decode:    DecodeSirc12,
	}
}

func NewIrpSirc15(protocol string) Irp {
	return &irpImpl{
		protocol:  protocol,
		frequency: FrequencySirc15,
		decode:    DecodeSirc15,
	}
}

func NewIrpSirc20(protocol string) Irp {
	return &irpImpl{
		protocol:  protocol,
		frequency: FrequencySirc20,
		decode:    DecodeSirc20,
	}
}

func DecodeSirc12(code SignalCode) (SignalData, error) {
	commandBits8 := GetBits8(code.Command[0], true)
	commandBits := commandBits8[0:7]

	addressBits8 := GetBits8(code.Address[0], true)
	addressBits := addressBits8[0:5]

	data := NewSignalData()
	data.Add(SIRC_PREAMBLE_MARK, SIRC_PREAMBLE_SPACE)
	data.AddBits(commandBits, SIRC_BIT1_MARK, SIRC_BIT1_SPACE, SIRC_BIT0_MARK, SIRC_BIT0_SPACE)
	data.AddBits(addressBits, SIRC_BIT1_MARK, SIRC_BIT1_SPACE, SIRC_BIT0_MARK, SIRC_BIT0_SPACE)
	data.Pop() // last space is not needed

	return data, nil
}

func DecodeSirc15(code SignalCode) (SignalData, error) {
	commandBits8 := GetBits8(code.Command[0], true)
	commandBits := commandBits8[0:7]

	addressBits8 := GetBits8(code.Address[0], true)
	addressBits := addressBits8[0:8]

	data := NewSignalData()
	data.Add(SIRC_PREAMBLE_MARK, SIRC_PREAMBLE_SPACE)
	data.AddBits(commandBits, SIRC_BIT1_MARK, SIRC_BIT1_SPACE, SIRC_BIT0_MARK, SIRC_BIT0_SPACE)
	data.AddBits(addressBits, SIRC_BIT1_MARK, SIRC_BIT1_SPACE, SIRC_BIT0_MARK, SIRC_BIT0_SPACE)
	data.Pop() // last space is not needed

	return data, nil
}

func DecodeSirc20(code SignalCode) (SignalData, error) {
	commandBits8 := GetBits8(code.Command[0], true)
	commandBits := commandBits8[0:7]

	addressBits8_0 := GetBits8(code.Address[0], true)
	addressBits8_1 := GetBits8(code.Address[1], true)
	addressBits := make([]bool, 0, 13)
	addressBits = append(addressBits, addressBits8_0[:]...)
	addressBits = append(addressBits, addressBits8_1[0:5]...)

	data := NewSignalData()
	data.Add(SIRC_PREAMBLE_MARK, SIRC_PREAMBLE_SPACE)
	data.AddBits(commandBits, SIRC_BIT1_MARK, SIRC_BIT1_SPACE, SIRC_BIT0_MARK, SIRC_BIT0_SPACE)
	data.AddBits(addressBits, SIRC_BIT1_MARK, SIRC_BIT1_SPACE, SIRC_BIT0_MARK, SIRC_BIT0_SPACE)
	data.Pop() // last space is not needed

	return data, nil
}

/***************************************************************************************************
*   Sony SIRC protocol description
*   https://www.sbprojects.net/knowledge/ir/sirc.php
*   http://picprojects.org.uk/
*   https://radiohlam.ru/sirc/
****************************************************************************************************
*      Preamble  Preamble     Pulse Width Modulation           Pause             Entirely repeat
*        mark     space                                     up to period             message..
*
*        2400      600      12/15/20 bits (600, 1200)         ...45000          2400      600
*     __________          _ _ _ _  _  _  _ _ _  _  _ _ _                    __________          _ _
* ____          __________ _ _ _ __ __ __ _ _ __ __ _ _ ____________________          __________ _
*                        |    command    |   address    |
*                 SIRC   |     7b LSB    |    5b LSB    |
*                 SIRC15 |     7b LSB    |    8b LSB    |
*                 SIRC20 |     7b LSB    |    13b LSB   |
*
* No way to determine either next message is repeat or not,
* so recognize only fact message received.Sony remotes always send at least 3 messages.
* Assume 8 last extended bits for SIRC20 are address bits.
***************************************************************************************************/

const (
	SIRC_CARRIER_FREQUENCY  = 40000
	SIRC_PREAMBLE_MARK      = 2400
	SIRC_PREAMBLE_SPACE     = 600
	SIRC_BIT1_MARK          = 1200
	SIRC_BIT1_SPACE         = 600
	SIRC_BIT0_MARK          = 600
	SIRC_BIT0_SPACE         = 600
	SIRC_PREAMBLE_TOLERANCE = 200
	SIRC_BIT_TOLERANCE      = 120
	SIRC_SILENCE            = 10000
	SIRC_MIN_SPLIT_TIME     = (SIRC_SILENCE - 1000)
	SIRC_REPEAT_PERIOD      = 45000
)
