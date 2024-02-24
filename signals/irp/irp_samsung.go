package irp

const (
	FrequencySamsung32 = 37900
)

func NewIrpSamsung32(protocol string) Irp {
	return &irpImpl{
		protocol:  protocol,
		frequency: FrequencySamsung32,
		decode:    DecodeSamsung32,
	}
}

func DecodeSamsung32(code SignalCode) (SignalData, error) {
	addressBits := GetBits8(code.Address[0], true)
	commandBits := GetBits8(code.Command[0], true)
	commandBitsReversed := GetBits8(0xff^(code.Command[0]), true)

	data := NewSignalData()
	data.Add(SAMSUNG_PREAMBLE_MARK, SAMSUNG_PREAMBLE_SPACE)
	data.AddBits8(addressBits, SAMSUNG_BIT1_MARK, SAMSUNG_BIT1_SPACE, SAMSUNG_BIT0_MARK, SAMSUNG_BIT0_SPACE)
	data.AddBits8(addressBits, SAMSUNG_BIT1_MARK, SAMSUNG_BIT1_SPACE, SAMSUNG_BIT0_MARK, SAMSUNG_BIT0_SPACE)
	data.AddBits8(commandBits, SAMSUNG_BIT1_MARK, SAMSUNG_BIT1_SPACE, SAMSUNG_BIT0_MARK, SAMSUNG_BIT0_SPACE)
	data.AddBits8(commandBitsReversed, SAMSUNG_BIT1_MARK, SAMSUNG_BIT1_SPACE, SAMSUNG_BIT0_MARK, SAMSUNG_BIT0_SPACE)
	data.Add(SAMSUNG_BIT1_MARK)

	return data, nil
}

/***************************************************************************************************
*   SAMSUNG32 protocol description
*   https://www.mikrocontroller.net/articles/IRMP_-_english#SAMSUNG
****************************************************************************************************
*  Preamble   Preamble     Pulse Distance/Width        Pause       Preamble   Preamble  Bit1  Stop
*    mark      space           Modulation                           repeat     repeat          bit
*                                                                    mark       space
*
*     4500      4500        32 bit + stop bit       40000/100000     4500       4500
*  __________          _  _ _  _  _  _ _ _  _  _ _                ___________            _    _
* _          __________ __ _ __ __ __ _ _ __ __ _ ________________           ____________ ____ ___
*
***************************************************************************************************/

/* Samsung silence have to be greater than REPEAT MAX
 * otherwise there can be problems during unit tests parsing
 * of some data. Real tolerances we don't know, but in real life
 * silence time should be greater than max repeat time. This is
 * because of similar preambule timings for repeat and first messages. */

const (
	SAMSUNG_PREAMBLE_MARK      = 4500
	SAMSUNG_PREAMBLE_SPACE     = 4500
	SAMSUNG_BIT1_MARK          = 550
	SAMSUNG_BIT1_SPACE         = 1650
	SAMSUNG_BIT0_MARK          = 550
	SAMSUNG_BIT0_SPACE         = 550
	SAMSUNG_REPEAT_PAUSE_MIN   = 30000
	SAMSUNG_REPEAT_PAUSE1      = 46000
	SAMSUNG_REPEAT_PAUSE2      = 97000
	SAMSUNG_MIN_SPLIT_TIME     = 5000
	SAMSUNG_SILENCE            = 145000
	SAMSUNG_REPEAT_PAUSE_MAX   = 140000
	SAMSUNG_REPEAT_MARK        = 4500
	SAMSUNG_REPEAT_SPACE       = 4500
	SAMSUNG_PREAMBLE_TOLERANCE = 200
	SAMSUNG_BIT_TOLERANCE      = 120
)
