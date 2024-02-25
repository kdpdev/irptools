package irp

const (
	FrequencyRca = 56000
)

func NewIrpRca(protocol string) Irp {
	return &irpImpl{
		protocol:  protocol,
		frequency: FrequencyRca,
		decode:    DecodeRca,
	}
}

func DecodeRca(code SignalCode) (SignalData, error) {
	commandBits8 := GetBits8(code.Command[0], false)
	commandBits := commandBits8[:]

	addressBits8 := GetBits8(code.Address[0], false)
	addressBits := addressBits8[4:]

	data := NewSignalData()
	data.Add(RCA_PREAMBLE_MARK, RCA_PREAMBLE_SPACE)
	data.AddBits(addressBits, RCA_BIT1_MARK, RCA_BIT1_SPACE, RCA_BIT0_MARK, RCA_BIT0_SPACE)
	data.AddBits(commandBits, RCA_BIT1_MARK, RCA_BIT1_SPACE, RCA_BIT0_MARK, RCA_BIT0_SPACE)
	data.AddBits(addressBits, RCA_BIT0_MARK, RCA_BIT0_SPACE, RCA_BIT1_MARK, RCA_BIT1_SPACE)
	data.AddBits(commandBits, RCA_BIT0_MARK, RCA_BIT0_SPACE, RCA_BIT1_MARK, RCA_BIT1_SPACE)
	data.Pop()
	dur := data.Duration()
	tie := RCA_SIGNAL_DUR - dur
	data.Add(tie)

	return data, nil
}

/***************************************************************************************************
*   https://www.sbprojects.net/knowledge/ir/rca.php
****************************************************************************************************
 */

const (
	RCA_PREAMBLE_MARK  = 4000
	RCA_PREAMBLE_SPACE = 4000
	RCA_BIT1_MARK      = 500
	RCA_BIT1_SPACE     = 2000
	RCA_BIT0_MARK      = 500
	RCA_BIT0_SPACE     = 1000
	RCA_SIGNAL_DUR     = 64000
)
