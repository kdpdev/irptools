package irp

const (
	FrequencyRc5  = 36000
	FrequencyRc5x = 36000
	FrequencyRc6  = 36000
)

func NewIrpRc5(protocol string) Irp {
	return &irpImpl{
		protocol:  protocol,
		frequency: FrequencyRc5,
		decode:    DecodeRc5,
	}
}

func NewIrpRc5x(protocol string) Irp {
	return NewIrpUnsupported(protocol, FrequencyRc5x)
}

func NewIrpRc6(protocol string) Irp {
	return &irpImpl{
		protocol:  protocol,
		frequency: FrequencyRc6,
		decode:    DecodeRc6,
	}
}

func DecodeRc5(code SignalCode) (SignalData, error) {
	addressBits := GetBits8(code.Address[0], false)
	commandBits := GetBits8(code.Command[0], false)

	bits := make([]bool, 0, 16)
	bits = append(bits, addressBits[3:]...)
	bits = append(bits, commandBits[2:]...)

	data := SignalData{RC5_BIT, RC5_BIT, RC5_BIT + RC5_BIT, RC5_BIT} // 2 start 1-bits + toggle bit

	lastBit := false
	for _, b := range bits {
		if lastBit {
			if b {
				data = append(data, RC5_BIT, RC5_BIT)
			} else {
				data[len(data)-1] += RC5_BIT
				data = append(data, RC5_BIT)
			}
		} else {
			if b {
				data[len(data)-1] += RC5_BIT
				data = append(data, RC5_BIT)
			} else {
				data = append(data, RC5_BIT, RC5_BIT)
			}
		}
		lastBit = b
	}

	if len(data)%2 == 0 {
		data = data[0 : len(data)-1]
	}

	return data, nil
}

func DecodeRc6(code SignalCode) (SignalData, error) {
	addressBits := GetBits8(code.Address[0], false)
	commandBits := GetBits8(code.Command[0], false)

	bits := make([]bool, 0, 16)
	bits = append(bits, addressBits[:]...)
	bits = append(bits, commandBits[:]...)

	data := SignalData{
		RC6_PREAMBLE_MARK, RC6_PREAMBLE_SPACE, // preamble
		RC6_BIT, RC6_BIT + RC6_BIT, RC6_BIT, // start 1-bit + 1st mode 0-bit
		RC6_BIT, RC6_BIT, RC6_BIT, RC6_BIT, // 2nd and 3rd mode 0-bits
		RC6_T_BIT, RC6_T_BIT, // toggle 1-bit
	}

	lastBit := false
	for _, b := range bits {
		if (lastBit && b) || (!lastBit && !b) {
			data = append(data, RC6_BIT, RC6_BIT)
		} else if (lastBit && !b) || (!lastBit && b) {
			data[len(data)-1] += RC6_BIT
			data = append(data, RC6_BIT)
		} else {
			return nil, Error("logic error")
		}
		lastBit = b
	}

	if len(data)%2 == 0 {
		data = data[0 : len(data)-1]
	}

	return data, nil
}

/***************************************************************************************************
*   RC5 protocol description
*   https://www.mikrocontroller.net/articles/IRMP_-_english#RC5_.2B_RC5X
*   https://radiohlam.ru/rc-5/
*   https://www.kernel.org/doc/html/latest/userspace-api/media/rc/rc-protos.html
****************************************************************************************************
*                                       Manchester/biphase
*                                           Modulation
*
*                              888/1776 - bit (x2 for toggle bit)
*
*                           __  ____    __  __  __  __  __  __  __  __
*                         __  __    ____  __  __  __  __  __  __  __  _
*                         | 1 | 1 | 0 |      ...      |      ...      |
*                           s  si   T   address (MSB)   command (MSB)
*
*    Note: manchester starts from space timing, so it have to be handled properly
*    s - start bit (always 1)
*    si - RC5: start bit (always 1), RC5X - 7-th bit of address (in our case always 0)
*    T - toggle bit, change it's value every button press
*    address - 5 bit
*    command - 6/7 bit
***************************************************************************************************/

const (
	RC5_PREAMBLE_MARK      = 0
	RC5_PREAMBLE_SPACE     = 0
	RC5_BIT                = 888       // half of time-quant for 1 bit
	RC5_PREAMBLE_TOLERANCE = 200       // us
	RC5_BIT_TOLERANCE      = 120       // us
	RC5_SILENCE            = 2700 * 10 // protocol allows 2700 silence, but it is hard to send 1 message without repeat */
	RC5_MIN_SPLIT_TIME     = 2700
)

/***************************************************************************************************
*   RC6 protocol description
*   https://www.mikrocontroller.net/articles/IRMP_-_english#RC6_.2B_RC6A
*   https://snrlab.in/8051/rc-6-protocol-and-interfacing-with-microcontroller/
*   http://www.pcbheaven.com/userpages/The_Philips_RC6_Protocol/
****************************************************************************************************
*      Preamble                       Manchester/biphase                       Silence
*     mark/space                          Modulation
*
*    2666     889        444/888 - bit (x2 for toggle bit)                       2666
*
*  ________         __    __  __  __    ____  __  __  __  __  __  __  __  __
* _        _________  ____  __  __  ____    __  __  __  __  __  __  __  __  _______________
*                   | 1 | 0 | 0 | 0 |   0   |      ...      |      ...      |             |
*                     s  m2  m1  m0     T     address (MSB)   command (MSB)
*
*    s - start bit (always 1)
*    m0-2 - mode (000 for RC6)
*    T - toggle bit, twice longer
*    address - 8 bit
*    command - 8 bit
***************************************************************************************************/

const (
	RC6_CARRIER_FREQUENCY  = 36000
	RC6_DUTY_CYCLE         = 0.33
	RC6_PREAMBLE_MARK      = 2666
	RC6_PREAMBLE_SPACE     = 889
	RC6_BIT                = 444 // half of time-quant for 1 bit
	RC6_T_BIT              = RC6_BIT * 2
	RC6_PREAMBLE_TOLERANCE = 200
	RC6_BIT_TOLERANCE      = 120
	RC6_SILENCE            = 2700 * 10 // protocol allows 2700 silence, but it is hard to send 1 message without repeat
	RC6_MIN_SPLIT_TIME     = 2700
)
