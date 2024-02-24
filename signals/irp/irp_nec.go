package irp

const (
	FrequencyNec    = 38000
	FrequencyNecExt = 38000
	FrequencyNec42  = 38000
)

func NewIrpNec(protocol string) Irp {
	return &irpImpl{
		protocol:  protocol,
		frequency: FrequencyNec,
		decode:    GetNecDecoder(1),
	}
}

func NewIrpNecExt(protocol string) Irp {
	return &irpImpl{
		protocol:  protocol,
		frequency: FrequencyNecExt,
		decode:    GetNecExtDecoder(1),
	}
}

func NewIrpNec42(protocol string) Irp {
	return &irpImpl{
		protocol:  protocol,
		frequency: FrequencyNec42,
		decode:    GetNecExtDecoder(1),
	}
}

func GetNecDecoder(repeatCodes int) func(code SignalCode) (SignalData, error) {
	return func(code SignalCode) (SignalData, error) {
		return DecodeNec(code, repeatCodes)
	}
}

func GetNecExtDecoder(repeatCodes int) func(code SignalCode) (SignalData, error) {
	return func(code SignalCode) (SignalData, error) {
		return DecodeNecExt(code, repeatCodes)
	}
}

func DecodeNec(code SignalCode, repeatCodes int) (SignalData, error) {
	data := NewSignalData()
	data.Add(NEC_PREAMBLE_MARK, NEC_PREAMBLE_SPACE)
	data.AddBits8(GetBits8(code.Address[0], true), NEC_BIT1_MARK, NEC_BIT1_SPACE, NEC_BIT0_MARK, NEC_BIT0_SPACE)
	data.AddBits8(GetBits8(code.Address[0], true), NEC_BIT0_MARK, NEC_BIT0_SPACE, NEC_BIT1_MARK, NEC_BIT1_SPACE)
	data.AddBits8(GetBits8(code.Command[0], true), NEC_BIT1_MARK, NEC_BIT1_SPACE, NEC_BIT0_MARK, NEC_BIT0_SPACE)
	data.AddBits8(GetBits8(code.Command[0], true), NEC_BIT0_MARK, NEC_BIT0_SPACE, NEC_BIT1_MARK, NEC_BIT1_SPACE)
	data.Add(NEC_BIT1_MARK)

	addNecRepeatCodes(&data, repeatCodes)

	return data, nil
}

func DecodeNecExt(code SignalCode, repeatCodes int) (SignalData, error) {
	data := NewSignalData()
	data.Add(NEC_PREAMBLE_MARK, NEC_PREAMBLE_SPACE)
	data.AddBits8(GetBits8(code.Address[0], true), NEC_BIT1_MARK, NEC_BIT1_SPACE, NEC_BIT0_MARK, NEC_BIT0_SPACE)
	data.AddBits8(GetBits8(code.Address[1], true), NEC_BIT1_MARK, NEC_BIT1_SPACE, NEC_BIT0_MARK, NEC_BIT0_SPACE)
	data.AddBits8(GetBits8(code.Command[0], true), NEC_BIT1_MARK, NEC_BIT1_SPACE, NEC_BIT0_MARK, NEC_BIT0_SPACE)
	data.AddBits8(GetBits8(code.Command[1], true), NEC_BIT1_MARK, NEC_BIT1_SPACE, NEC_BIT0_MARK, NEC_BIT0_SPACE)
	data.Add(NEC_BIT1_MARK)

	addNecRepeatCodes(&data, repeatCodes)

	return data, nil
}

func addNecRepeatCodes(data *SignalData, repeatCodes int) {
	if repeatCodes < 0 {
		return
	}
	data.Add(Micros(NEC_REPEAT_PERIOD) - data.Duration())
	const leadingBurst = Micros(NEC_REPEAT_MARK)
	const space = Micros(NEC_REPEAT_SPACE)
	const mark = Micros(NEC_BIT0_MARK)
	const pause = Micros(NEC_REPEAT_PERIOD) - leadingBurst - space - mark
	for i := 0; i < repeatCodes; i++ {
		data.Add(leadingBurst, space, mark, pause)
	}
}

/***************************************************************************************************
*   NEC protocol description
*   https://radioparty.ru/manuals/encyclopedia/213-ircontrol?start=1
*   https://radiohlam.ru/nec/
*   https://techdocs.altium.com/display/FPGA/NEC+Infrared+Transmission+Protocol
****************************************************************************************************
*     Preamble   Preamble      Pulse Distance/Width          Pause       Preamble   Preamble  Stop
*       mark      space            Modulation             up to period    repeat     repeat    bit
*                                                                          mark       space
*
*        9000      4500         32 bit + stop bit         ...110000         9000       2250
*     __________          _ _ _ _  _  _  _ _ _  _  _ _ _                ___________            _
* ____          __________ _ _ _ __ __ __ _ _ __ __ _ _ ________________           ____________ ___
*
***************************************************************************************************/

const (
	NEC_PREAMBLE_MARK      = 9000
	NEC_PREAMBLE_SPACE     = 4500
	NEC_BIT1_MARK          = 560
	NEC_BIT1_SPACE         = 1690
	NEC_BIT0_MARK          = 560
	NEC_BIT0_SPACE         = 560
	NEC_REPEAT_PERIOD      = 108000
	NEC_SILENCE            = NEC_REPEAT_PERIOD
	NEC_MIN_SPLIT_TIME     = NEC_REPEAT_PAUSE_MIN
	NEC_REPEAT_PAUSE_MIN   = 4000
	NEC_REPEAT_PAUSE_MAX   = 150000
	NEC_REPEAT_MARK        = 9000
	NEC_REPEAT_SPACE       = 2250
	NEC_PREAMBLE_TOLERANCE = 200
	NEC_BIT_TOLERANCE      = 120
)
